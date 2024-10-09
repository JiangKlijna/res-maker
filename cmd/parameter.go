package cmd

/*
#include <stdio.h>

// print c compiler version
void print_c_compiler_version() {
    #ifdef __GNUC__
        printf("c compiler: GCC %d.%d.%d\n", __GNUC__, __GNUC_MINOR__, __GNUC_PATCHLEVEL__);
    #elif defined(__clang__)
        printf("c compiler: Clang %s\n", __clang_version__);
    #elif defined(_MSC_VER)
        printf("c compiler: MSVC %d\n", _MSC_VER);
    #else
        printf("c compiler: Unknown\n");
    #endif
}
*/
import "C"
import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jiangklijna/res-maker/def"
)

// Parameter interface
type Parameter interface {
	IsConfigFile() bool
	GetConfigFile() map[def.ResType][24]uint
	GetParameter() map[def.ResType]uint
	GetDuration() time.Duration
}

// parameter Command line parameters
type parameter struct {
	nCores     int
	nMemGs     int
	duration   string
	configFile string

	timeDuration time.Duration

	cpuNs *[24]uint
	memNs *[24]uint
}

// GetParameter get and check parameter
func GetParameter() Parameter {
	p := new(parameter)
	flag.IntVar(&p.nCores, "c", 0, "Number of cores to use")
	flag.IntVar(&p.nMemGs, "m", 0, "Size of memory to allocate in GB")
	flag.StringVar(&p.duration, "d", "", "Specify the duration (e.g., 1s, 2m, 3h, 4d) for the program to run")
	flag.StringVar(&p.configFile, "f", "", "Specify the configuration file for the program.")

	help := flag.Bool("h", false, "this help")
	version := flag.Bool("v", false, "show version and exit")

	flag.Parse()
	if *help {
		flag.Usage()
	} else if *version {
		fmt.Println(def.Name, "version:", def.Version)
		fmt.Println("golang", "version:", runtime.Version())
		C.print_c_compiler_version()
	} else {
		p.init()
		return p
	}
	return nil
}

// init
func (p *parameter) init() {
	if len(p.configFile) > 0 {
		content, err := os.ReadFile(p.configFile)
		if err != nil {
			panic(err)
		}
		config := make(map[string]map[string]*uint)
		err = json.Unmarshal(content, &config)
		if err != nil {
			panic(err)
		}
		if m, isOk := config[def.Cpu.Name()]; isOk {
			p.cpuNs = convert(m)
		}
		if m, isOk := config[def.Mem.Name()]; isOk {
			p.memNs = convert(m)
		}
		if p.cpuNs == nil && p.memNs == nil {
			fmt.Println(time.Now().Format(time.DateTime), "Error: config[", p.configFile, "], cpu or mem is not null")
			os.Exit(1)
		}
	} else {
		if p.nCores < 0 {
			fmt.Println(time.Now().Format(time.DateTime), "Warning: cpu should be >= 0")
			p.nCores = 0
		}
		if p.nMemGs < 0 {
			fmt.Println(time.Now().Format(time.DateTime), "Warning: mem should be >= 0")
			p.nMemGs = 0
		}
		if p.nCores == 0 && p.nMemGs == 0 {
			fmt.Println(time.Now().Format(time.DateTime), "Error: cpu or mem must be >= 0")
			os.Exit(1)
		}
	}
	if len(p.duration) > 0 {
		d, err := parseDuration(p.duration)
		if err != nil {
			fmt.Println(time.Now().Format(time.DateTime), "Error:")
			panic(err)
		}
		p.timeDuration = d
	}
}

func (p *parameter) IsConfigFile() bool {
	return len(p.configFile) > 0
}

func (p *parameter) GetConfigFile() map[def.ResType][24]uint {
	m := make(map[def.ResType][24]uint)
	if p.cpuNs != nil {
		m[def.Cpu] = *p.cpuNs
	}
	if p.memNs != nil {
		m[def.Mem] = *p.memNs
	}
	if len(m) == 0 {
		panic("config file error")
	}
	return m
}

func (p *parameter) GetParameter() map[def.ResType]uint {
	m := make(map[def.ResType]uint)
	if p.nCores > 0 {
		m[def.Cpu] = uint(p.nCores)
	}
	if p.nMemGs > 0 {
		m[def.Mem] = uint(p.nMemGs)
	}
	if len(m) == 0 {
		panic("parameter error")
	}
	return m
}

func (p *parameter) GetDuration() time.Duration {
	return p.timeDuration
}

func convert(m map[string]*uint) *[24]uint {
	var arr [24]uint
	if len(m) == 0 {
		return &arr
	}
	var pre uint
	for i := 23; i >= 0; i-- {
		u := m[formatIndex(i)]
		if u != nil {
			pre = *u
			break
		}
	}
	for i := 0; i < 24; i++ {
		u := m[formatIndex(i)]
		if u != nil {
			pre = *u
		}
		arr[i] = pre
	}
	return &arr
}

func formatIndex(i int) string {
	s := strconv.FormatInt(int64(i), 10)
	if i < 10 {
		return "0" + s
	}
	return s
}

// parseDuration parses strings like "1s", "2m", "3h", "4d" and returns the corresponding time.Duration.
func parseDuration(input string) (time.Duration, error) {
	// Remove the last character which represents the unit
	valueStr := strings.TrimRight(input, "dsmh")
	unit := strings.ToLower(input[len(valueStr):])

	// Convert the string value to an integer
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("-d invalid input: %s", input)
	}

	// Calculate the duration based on the unit
	switch unit {
	case "", "s":
		return time.Duration(value) * time.Second, nil
	case "m":
		return time.Duration(value) * time.Minute, nil
	case "h":
		return time.Duration(value) * time.Hour, nil
	case "d":
		return time.Duration(value) * 24 * time.Hour, nil
	default:
		return 0, fmt.Errorf("-d unsupported unit: %s", unit)
	}
}
