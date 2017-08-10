package main

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

type HostStat struct {
	OS       string
	Platform string
	Virt     string
	CPU      string
	Cores    int
	Clock    float64
}

// Host returns general system information
func Host() (*HostStat, error) {
	stat := HostStat{}

	host, err := host.Info()
	if err != nil {
		return nil, err
	}

	stat.OS = host.OS
	stat.Platform = fmt.Sprintf("%s %s", host.Platform, host.PlatformVersion)
	stat.Virt = host.VirtualizationSystem

	cpu, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	stat.CPU = cpu[0].ModelName
	stat.Cores = len(cpu)
	stat.Clock = cpu[0].Mhz

	return &stat, nil
}
