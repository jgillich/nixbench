package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

// DiskStat stores disk speeds in MB/s
type DiskStat struct {
	Average float64
	Speeds  []float64
}

func Disk() (*DiskStat, error) {
	stat := DiskStat{}

	zero, err := os.Open("/dev/zero")
	defer zero.Close()
	if err != nil {
		return nil, err
	}

	speeds := []float64{}

	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("diskbench_%v", i)

		f, err := os.Create(name)
		if err != nil {
			return nil, err
		}

		start := time.Now()

		if _, err := io.CopyN(f, zero, bytefmt.GIGABYTE); err != nil {
			return nil, err
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

	return &stat, nil
}
