package utils

import (
	"log"
	"sync/atomic"
	"time"

	// "github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/cpu"
	// "github.com/shirou/gopsutil/v3/cpu"
)

var cpuValue atomic.Value

const decay = 0.95

func init() {
	c := cpuinterval{500 * time.Microsecond}
	go c.Start()
}

func GetCpuUsage() float64 {
	v, ok := cpuValue.Load().(float64)
	if ok {
		return v
	}
	return 0
}

type cpuinterval struct {
	internal time.Duration
}

func (c *cpuinterval) Start() {
	tk := time.NewTimer(c.internal)
	for {
		<-tk.C
		old := GetCpuUsage()
		v := c.GetCPU()
		new := (1-decay)*old + decay*v
		log.Println(old, v, new)
		cpuValue.Store(new)
	}
}

func (c *cpuinterval) GetCPU() float64 {
	// info, _ := cpu.Info()
	// log.Println(info)
	// cpuValue.Store(info)
	r, _ := cpu.Percent(100*c.internal, false)
	return r[0]
}
