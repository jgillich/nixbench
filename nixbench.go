package main

import (
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"

	"github.com/urfave/cli"
)

// VERSION is set at build time
var VERSION = "master"

type Result struct {
	CPU       *CPUStat
	Disk      *DiskStat
	Geekbench *GeekbenchStat
	Host      *HostStat
	Net       *NetStat
}

func main() {

	app := &cli.App{
		Name:    "nixbench",
		Usage:   "A better benchmarking tool for servers",
		Version: VERSION,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "yaml",
				Usage: "Output as yaml",
			},
		},
		Action: func(c *cli.Context) error {
			nixbench := Nixbench{
				Yaml: c.GlobalBool("yaml"),
			}

			return nixbench.Run()
		},
	}

	app.Run(os.Args)
}

type Nixbench struct {
	Yaml bool
}

func (n *Nixbench) Printf(format string, a ...interface{}) {
	if !n.Yaml {
		fmt.Printf(format, a...)
	}
}

func (n *Nixbench) Run() error {
	n.Printf("nixbench %s - https://github.com/jgillich/nixbench", VERSION)

	n.Printf("\n\n")
	n.Printf("Host\n")
	n.Printf("----\n\n")
	host, err := Host()
	if err != nil {
		return err
	}
	n.Printf("%-10s: %s\n", "OS", host.OS)
	n.Printf("%-10s: %s\n", "Platform", host.Platform)
	n.Printf("%-10s: %s\n", "CPU", host.CPU)
	n.Printf("%-10s: %d\n", "Cores", host.Cores)
	n.Printf("%-10s: %d Mhz\n", "Clock", int(host.Clock))
	n.Printf("%-10s: %d MB\n", "RAM", host.RAM)

	n.Printf("\n\n")
	n.Printf("CPU\n")
	n.Printf("---\n\n")
	cpu, err := CPU()
	if err != nil {
		return err
	}
	n.Printf("Sha256  : %.2f seconds\n", cpu.Sha256)
	n.Printf("Gzip    : %.2f seconds\n", cpu.Gzip)

	n.Printf("\n\n")
	n.Printf("Disk\n")
	n.Printf("----\n\n")
	disk, err := Disk()
	if err != nil {
		return err
	}
	for i, speed := range disk.Speeds {
		n.Printf("%d. run   : %d MB/s\n", i+1, int(speed))
	}
	n.Printf("Average  : %d MB/s\n", int(disk.Average))

	n.Printf("\n\n")
	n.Printf("Geekbench\n")
	n.Printf("---------\n\n")
	gb, err := Geekbench()
	if err != nil {
		return err
	}
	n.Printf("Single-Core Score  : %d\n", gb.SingleCore)
	n.Printf("Multi-Core Score   : %d\n", gb.MultiCore)
	n.Printf("Result URL         : %s\n", gb.URL)

	n.Printf("\n\n")
	n.Printf("Net\n")
	n.Printf("---\n\n")
	net, err := Net()
	if err != nil {
		return err
	}

	for _, f := range files {
		n.Printf("%-30s: %-6.2f MB/s\n", f.Key, (*net)[f.Key])
	}

	if n.Yaml {
		res := Result{
			Host:      host,
			CPU:       cpu,
			Disk:      disk,
			Geekbench: gb,
			Net:       net,
		}

		yml, err := yaml.Marshal(res)
		if err != nil {
			return err
		}
		fmt.Printf(string(yml))
	}

	return nil
}
