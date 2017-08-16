// +build !openbsd

package modules

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"

	"code.cloudfoundry.org/bytefmt"
)

func init() {
	Modules["host"] = &Host{}
}

type Host struct {
	OS       string
	Platform string
	CPU      string
	Cores    int
	Clock    float64
	RAM      uint64
}

func (stat *Host) Run() error {
	host, err := host.Info()
	if err != nil {
		return err
	}

	stat.OS = host.OS
	stat.Platform = fmt.Sprintf("%s %s", host.Platform, host.PlatformVersion)

	cpu, err := cpu.Info()
	if err != nil {
		return err
	}

	stat.CPU = cpu[0].ModelName
	stat.Cores = len(cpu)
	stat.Clock = cpu[0].Mhz

	vm, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	stat.RAM = vm.Total / bytefmt.MEGABYTE

	return nil
}

func (stat *Host) Print() {
	fmt.Printf("%-10s: %s\n", "OS", stat.OS)
	fmt.Printf("%-10s: %s\n", "Platform", stat.Platform)
	fmt.Printf("%-10s: %s\n", "CPU", stat.CPU)
	fmt.Printf("%-10s: %d\n", "Cores", stat.Cores)
	fmt.Printf("%-10s: %d Mhz\n", "Clock", int(stat.Clock))
	fmt.Printf("%-10s: %d MB\n", "RAM", stat.RAM)
}
