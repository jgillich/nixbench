package main

import (
	"compress/gzip"
	"crypto/sha256"
	"io"
	"os"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

// DiskStat stores disk speeds in MB/s
type CPUStat struct {
	Sha256 float64
	Gzip   float64
}

func CPU() (*CPUStat, error) {
	stat := CPUStat{}

	zero, err := os.Open("/dev/zero")
	defer zero.Close()
	if err != nil {
		return nil, err
	}

	null, err := os.Create("/dev/null")
	defer null.Close()
	if err != nil {
		return nil, err
	}

	hashStart := time.Now()
	hash := sha256.New()
	if _, err := io.CopyN(hash, zero, bytefmt.GIGABYTE); err != nil {
		return nil, err
	}
	hash.Sum(nil)
	stat.Sha256 = time.Since(hashStart).Seconds()

	gzipStart := time.Now()
	gz := gzip.NewWriter(null)
	if _, err := io.CopyN(gz, zero, bytefmt.GIGABYTE); err != nil {
		return nil, err
	}
	gz.Close()
	stat.Gzip = time.Since(gzipStart).Seconds()

	return &stat, nil

}
