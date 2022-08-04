package utils

import "testing"

func TestCpu(t *testing.T) {
	c := GetCPU()
	t.Log(c)
}
