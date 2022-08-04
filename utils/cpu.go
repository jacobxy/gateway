package utils

import (
	"log"
	"sync/atomic"

	"github.com/shirou/gopsutil/cpu"
)

var cpuValue atomic.Value

func GetCPU() atomic.Value {
	info, _ := cpu.Info()
	log.Println(info)
	cpuValue.Store(info)
	return cpuValue
}
