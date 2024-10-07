package run

import (
	"fmt"
	"os"
	"time"

	"github.com/jiangklijna/res-maker/cmd"
	"github.com/jiangklijna/res-maker/def"
	"github.com/jiangklijna/res-maker/res"
)

// Run main run
func Run() {
	p := cmd.GetParameter()
	if p == nil {
		return
	}
	// exit after duration if specified
	if p.GetDuration() > 0 {
		go func() {
			time.Sleep(p.GetDuration())
			os.Exit(0)
		}()
	}
	// run by parameter or config file
	if p.IsConfigFile() {
		runConfigFile(p.GetConfigFile())
	} else {
		runParameter(p.GetParameter())
	}
	select {}
}

// runConfigFile run by config file
func runConfigFile(m map[def.ResType][24]uint) {
	now := time.Now()
	hour := now.Hour()

	// run by this hour config
	resMap := runParameter((func() map[def.ResType]uint {
		_parameter := map[def.ResType]uint{}
		for resType, ns := range m {
			_parameter[resType] = ns[hour]
		}
		return _parameter
	})())
	// sleep to next hour
	time.Sleep(time.Now().Truncate(time.Hour).Add(time.Hour).Sub(time.Now()))

	for {
		// check parameter, whether to replace Res
		go func() {
			now = time.Now()
			hour = now.Hour()
			for resType, r := range resMap {
				ns := m[resType]
				if r.Num() != ns[hour] {
					r.Free()
					time.Sleep(time.Second)
					r = res.NewRes(ns[hour], resType)
					go r.Eat()
					resMap[resType] = r
				}
				fmt.Println(now.Format(time.DateTime), resType.Name(), "Usage:", ns[hour], resType.Unit())
			}
		}()
		// sleep to next hour
		time.Sleep(time.Hour)
	}
}

// runParameter run by parameter
func runParameter(m map[def.ResType]uint) map[def.ResType]res.Res {
	arr := make(map[def.ResType]res.Res)
	for resType, num := range m {
		newRes := res.NewRes(num, resType)
		go newRes.Eat()
		arr[resType] = newRes
		fmt.Println(time.Now().Format(time.DateTime), resType.Name(), "Usage:", num, resType.Unit())
	}
	return arr
}
