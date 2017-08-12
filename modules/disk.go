package modules

import (
	"fmt"
	"io"
	"os"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

func init() {
	Modules["disk"] = &Disk{}
}

type Disk struct {
	Average float64
	Speeds  []float64
}

func (stat *Disk) Run() error {
	zero, err := os.Open("/dev/zero")
	defer zero.Close()
	if err != nil {
		return err
	}

	speeds := []float64{}

	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("diskbench_%v", i)

		f, err := os.Create(name)
		if err != nil {
			return err
		}

		start := time.Now()

		if _, err := io.CopyN(f, zero, bytefmt.GIGABYTE); err != nil {
			return err
		}
		f.Sync()
		f.Close()

		duration := time.Since(start)
		speeds = append(speeds, 1/duration.Seconds()*1024)

		os.Remove(name)
	}

	var total float64
	for _, value := range speeds {
		total += value
	}

	stat.Average = total / float64(len(speeds))
	stat.Speeds = speeds

	return nil
}

func (stat *Disk) Print() {
	for i, speed := range stat.Speeds {
		fmt.Printf("%d. run   : %d MB/s\n", i+1, int(speed))
	}
	fmt.Printf("Average  : %d MB/s\n", int(stat.Average))
}
