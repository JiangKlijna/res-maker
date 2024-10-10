package def

import (
	"fmt"
	"time"
)

// Name of this program
const Name = "res-marker"

// Version of this program
const Version = "0.3"

// ResType Resource type
type ResType uint

const (
	Cpu ResType = iota + 1
	Mem
)

// Name of resource type
func (rt ResType) Name() string {
	switch rt {
	case Cpu:
		return "cpu"
	case Mem:
		return "mem"
	default:
		panic("unreachable code")
	}
}

// Unit of resource type
func (rt ResType) Unit() string {
	switch rt {
	case Cpu:
		return "core"
	case Mem:
		return "G"
	default:
		panic("unreachable code")
	}
}

// LogInfo log info
func LogInfo(a ...any) {
	fmt.Print(time.Now().Format(time.DateTime), " [Info] ")
	fmt.Println(a...)
}

// LogWarning log warning
func LogWarning(a ...any) {
	fmt.Print(time.Now().Format(time.DateTime), " [Warning] ")
	fmt.Println(a...)
}

// LogError log error
func LogError(a ...any) {
	fmt.Print(time.Now().Format(time.DateTime), " [Error] ")
	fmt.Println(a...)
}
