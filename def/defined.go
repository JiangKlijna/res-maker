package def

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
