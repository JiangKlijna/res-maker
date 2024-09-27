package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

func runCpu(nCores int) {
	run := func() {
		const z = 0
		const j = 2342
		var m = 9823
		var k = 31455
		var l = 16452
		var i = 1000000
		for {
			m = m ^ l
			k = (k / m * j) % i
			l = z * m * k
			i = (z * k) ^ m
			k = (k / m * j) % i
			m = m ^ l
			m = m ^ l
			i = (z * k) ^ m
			k = (k / m * j) % i
			m = i * i * i * i * i * i * i // m=k*l*j*l;
			m = m ^ l
			k = (k / m * j) % i
			l = z * m * k
			i = (z * k) ^ m
			l = (k / m * j) % i
			m = m ^ l
			m = m ^ l
			i = (z * k) ^ m
			k = (k / m * j) % i
			m = k*k*k*k - m/i
		}
	}
	// Ensure we don't request more cores than available.
	if nCores > runtime.NumCPU() {
		fmt.Printf("Warning: Requested cores (%d) exceed the available cores (%d). Using all available cores.\n", nCores, runtime.NumCPU())
		nCores = runtime.NumCPU()
	}
	if nCores > 0 {
		// Set the number of cores to use.
		runtime.GOMAXPROCS(nCores)

		// Start the specified number of goroutines.
		for i := 0; i < nCores; i++ {
			go func() {
				run()
			}()
		}
	}
}

func runMem(nMemGs int) {
	if nMemGs <= 0 {
		return
	}
	// Convert GB to bytes
	requestedSizeInBytes := uint64(nMemGs) * (1024 * 1024 * 1024)

	// Get current memory statistics
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	// Estimate available memory based on MemStats
	// Assuming system has at least 2GB free for other processes
	safeLimit := memStats.Sys - (2 * 1024 * 1024 * 1024)

	if requestedSizeInBytes > safeLimit {
		fmt.Printf("Warning: Requested memory (%d GB) exceeds safe limit by %d GB.\n", nMemGs, (requestedSizeInBytes-safeLimit)/(1024*1024*1024))
		fmt.Printf("Adjusting to maximum safe size (%d GB).\n", safeLimit/(1024*1024*1024))
		requestedSizeInBytes = safeLimit
	}

	// Allocate memory using make([]byte, sizeInBytes) instead of unsafe package
	memoryBlock := make([]byte, requestedSizeInBytes)

	// Print the address of the first byte to verify allocation
	fmt.Printf("Allocated %d GB of memory at address: %p\n", requestedSizeInBytes/(1024*1024*1024), unsafe.Pointer(&memoryBlock[0]))

	// Simulate using the memory block
	for {
		for i := range memoryBlock {
			memoryBlock[i] = byte(i)
		}

		time.Sleep(time.Minute)
	}
}

func main() {
	// Define the flag to specify the number of cores.
	nCores := flag.Int("c", 0, "Number of cores to use")
	nMemGs := flag.Int("m", 0, "Size of memory to allocate in GB")
	flag.Parse()

	if *nCores == 0 && *nMemGs == 0 {
		flag.Usage()
		return
	}

	go runMem(*nMemGs)
	go runCpu(*nCores)

	// Keep the main goroutine alive and listening for a signal to stop.
	select {}
	// Optionally you can add logic here to listen for signals and close the stopChan.
	// <-stopChan
	// close(stopChan)
	// wg.Wait()
}
