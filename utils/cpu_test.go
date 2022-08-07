package utils

import (
	"testing"
	"time"
)

func TestCpu(t *testing.T) {
	c := cpuinterval{50 * time.Microsecond}
	// c.GetCPU()
	t.Log(c.GetCPU())
}
