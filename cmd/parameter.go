package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
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

	cpuNs [24]uint
	memNs [24]uint
}

// GetParameter get and check parameter
func GetParameter() Parameter {
	p := new(parameter)
	flag.IntVar(&p.nCores, "c", 0, "Number of cores to use")
	flag.IntVar(&p.nMemGs, "m", 0, "Size of memory to allocate in GB")
	flag.StringVar(&p.duration, "d", "", "Specify the duration (e.g., 1s, 1d, 3m) for the program to run")
	flag.StringVar(&p.configFile, "f", "", "Specify the configuration file for the program.")

	help := flag.Bool("h", false, "this help")
	version := flag.Bool("v", false, "show version and exit")

	flag.Parse()
	if *help {
		flag.Usage()
	} else if *version {
		println("version", def.Version)
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
		c := new(configFile)
		err = json.Unmarshal(content, c)
		if err != nil {
			panic(err)
		}
		p.cpuNs = convert(c.Cpu)
		p.memNs = convert(c.Mem)
	} else {
		if p.nCores < 0 {
			fmt.Println("Warning: cpu should be >= 0")
			p.nCores = 0
		}
		if p.nMemGs < 0 {
			fmt.Println("Warning: mem should be >= 0")
			p.nMemGs = 0
		}
	}
}

func (p *parameter) IsConfigFile() bool {
	return len(p.configFile) > 0
}

func (p *parameter) GetConfigFile() map[def.ResType][24]uint {
	m := make(map[def.ResType][24]uint)
	m[def.Cpu] = p.cpuNs
	m[def.Mem] = p.memNs
	return m
}

func (p *parameter) GetParameter() map[def.ResType]uint {
	m := make(map[def.ResType]uint)
	m[def.Cpu] = uint(p.nCores)
	m[def.Mem] = uint(p.nMemGs)
	return m
}

func (p *parameter) GetDuration() time.Duration {
	return 0
}

type configFile struct {
	Cpu map[string]*uint `json:"cpu"`
	Mem map[string]*uint `json:"mem"`
}

func convert(m map[string]*uint) [24]uint {
	var arr [24]uint
	if len(m) == 0 {
		return arr
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
	return arr
}

func formatIndex(i int) string {
	s := strconv.FormatInt(int64(i), 10)
	if i < 10 {
		return "0" + s
	}
	return s
}
