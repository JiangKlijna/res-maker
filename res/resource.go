package res

/*
#include <stdlib.h> // use malloc, free
*/
import "C"
import (
	"runtime"
	"sync"
	"time"
	"unsafe"

	"github.com/jiangklijna/res-maker/def"
)

// Res interface
type Res interface {
	Num() uint
	Type() def.ResType
	Eat()
	Free()
}

// NewRes return a new Res
func NewRes(num uint, rt def.ResType) Res {
	return &res{num: num, resType: rt}
}

// res struct
type res struct {
	num     uint
	resType def.ResType
	stopCh  chan uint
}

func (r *res) Num() uint {
	return r.num
}

func (r *res) Type() def.ResType {
	return r.resType
}

func (r *res) Eat() {
	r.stopCh = make(chan uint)
	if r.num == 0 {
		nilEat(r.stopCh)
	} else {
		switch r.resType {
		case def.Cpu:
			cpuEat(r.num, r.stopCh)
		case def.Mem:
			memEat(r.num, r.stopCh)
		}
	}
}

func (r *res) Free() {
	close(r.stopCh)
	r.stopCh = nil
}

// nilEat nothing
func nilEat(stopCh chan uint) {
	for {
		select {
		case <-stopCh:
			return
		default:
			time.Sleep(1 * time.Second)
		}
	}
}

// cpuEat cpu eating
func cpuEat(num uint, stopCh chan uint) {
	var wg sync.WaitGroup
	run := func() {
		defer wg.Done()

		const z = 0
		const j = 2342
		var m = 9823
		var k = 31455
		var l = 16452
		var i = 1000000

		for {
			select {
			case <-stopCh:
				return
			default:
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
	}
	runtime.GOMAXPROCS(int(num))
	// Start the specified number of goroutines.
	for i := uint(0); i < num; i++ {
		wg.Add(1)
		go run()
	}
	wg.Wait()
}

// memEat mem eating
func memEat(num uint, stopCh chan uint) {
	// Convert GB to bytes
	requestedSizeInBytes := uint64(num) * (1024 * 1024 * 1024)
	// Allocate memory using make([]byte, sizeInBytes) instead of unsafe package
	ptr := C.malloc(C.size_t(requestedSizeInBytes))
	if ptr == nil {
		def.LogWarning("Warning: Failed to allocate", num, "G memory.")
		nilEat(stopCh)
		return
	}
	memoryBlock := unsafe.Slice((*byte)(ptr), int(requestedSizeInBytes))
	for i := range memoryBlock {
		memoryBlock[i] = byte(i)
	}

	// Suspend access every 600 times to prevent being compressed into virtual memory by OS
	count := requestedSizeInBytes / 600
	cursor := uint64(0)
	// Simulate using the memory block
	for {
		select {
		case <-stopCh:
			C.free(ptr)
			ptr = nil
			return
		default:
			for j := uint64(0); j < count; j++ {
				cursor++
				if cursor >= requestedSizeInBytes {
					cursor = 0
				}
				memoryBlock[cursor] = byte(j)
			}
			time.Sleep(1 * time.Second)
		}
	}
}
