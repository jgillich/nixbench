package main

import (
	"fmt"
	"log"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/urfave/cli"

	"github.com/briandowns/spinner"
)

type Result struct {
	CPU       *CPUStat
	Disk      *DiskStat
	Geekbench *GeekbenchStat
	Host      *HostStat
	Net       *NetStat
}

var geekbench bool

func main() {

	app := &cli.App{
		Name:  "nixbench",
		Usage: "A better benchmarking tool for servers",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "yaml",
				Usage: "Output as yaml",
			},
			&cli.BoolFlag{
				Name:  "no-geekbench",
				Usage: "Doesn't run Geekbench",
			},
		},
		Action: func(c *cli.Context) error {
			geekbench = !c.GlobalBool("no-geekbench")
			if c.GlobalBool("yaml") {
				runYaml()
			} else {
				runInteractive()
			}
			return nil
		},
	}

	app.Run(os.Args)
}

func runInteractive() {
	fmt.Print("nixbench 0.4.2 - https://github.com/jgillich/nixbench")
	spin := spinner.New(spinner.CharSets[12], time.Second/10)

	fmt.Print("\n\n")
	fmt.Print("Host\n")
	fmt.Print("----\n\n")
	host, err := Host()
	handleErr(err)
	fmt.Printf("%-10s: %s\n", "OS", host.OS)
	fmt.Printf("%-10s: %s\n", "Platform", host.Platform)
	fmt.Printf("%-10s: %s\n", "CPU", host.CPU)
	fmt.Printf("%-10s: %d\n", "Cores", host.Cores)
	fmt.Printf("%-10s: %d Mhz\n", "Clock", int(host.Clock))
	fmt.Printf("%-10s: %d MB\n", "RAM", host.RAM)

	fmt.Print("\n\n")
	fmt.Print("CPU\n")
	fmt.Print("---\n\n")
	spin.Start()
	cpu, err := CPU()
	spin.Stop()
	handleErr(err)
	fmt.Printf("Sha256  : %.2f seconds\n", cpu.Sha256)
	fmt.Printf("Gzip    : %.2f seconds\n", cpu.Gzip)

	fmt.Print("\n\n")
	fmt.Print("Disk\n")
	fmt.Print("----\n\n")
	spin.Start()
	disk, err := Disk()
	spin.Stop()
	handleErr(err)
	for i, speed := range disk.Speeds {
		fmt.Printf("%d. run   : %d MB/s\n", i+1, int(speed))
	}
	fmt.Printf("Average  : %d MB/s\n", int(disk.Average))

	if geekbench {
		fmt.Print("\n\n")
		fmt.Print("Geekbench\n")
		fmt.Print("---------\n\n")
		spin.Start()
		gb, err := Geekbench()
		spin.Stop()
		handleErr(err)
		fmt.Printf("Single-Core Score  : %d\n", gb.SingleCore)
		fmt.Printf("Multi-Core Score   : %d\n", gb.MultiCore)
		fmt.Printf("Result URL         : %s\n", gb.URL)
	}

	fmt.Print("\n\n")
	fmt.Print("Net\n")
	fmt.Print("---\n\n")
	spin.Start()
	net, err := Net()
	spin.Stop()
	handleErr(err)
	for _, f := range files {
		fmt.Printf("%-30s: %d MB/s\n", f.Key, int((*net)[f.Key]))
	}

}

func runYaml() {
	host, err := Host()
	handleErr(err)

	cpu, err := CPU()
	handleErr(err)

	disk, err := Disk()
	handleErr(err)

	gb, err := Geekbench()
	handleErr(err)

	net, err := Net()
	handleErr(err)

	res := Result{
		Host:      host,
		CPU:       cpu,
		Disk:      disk,
		Geekbench: gb,
		Net:       net,
	}

	yml, err := yaml.Marshal(res)
	handleErr(err)

	fmt.Printf(string(yml))
}

func handleErr(err error) {
	if err == nil {
		return
	}
	log.Fatalf("error: %v", err)
}
