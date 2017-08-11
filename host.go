package main

import (
	"fmt"

	"code.cloudfoundry.org/bytefmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

type HostStat struct {
	OS       string
	Platform string
	CPU      string
	Cores    int
	Clock    float64
	RAM      uint64
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

	cpu, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	stat.CPU = cpu[0].ModelName
	stat.Cores = len(cpu)
	stat.Clock = cpu[0].Mhz

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	stat.RAM = vm.Total / bytefmt.MEGABYTE

	return &stat, nil
}
