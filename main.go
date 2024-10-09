//go:generate windres assets/main.rc -o main.syso
package main

import (
	"github.com/jiangklijna/res-maker/run"
)

func main() {
	run.Run()
}
