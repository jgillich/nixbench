package modules

import (
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

func init() {
	Modules["cpu"] = &CPU{}
}

type CPU struct {
	Sha256 float64
	Gzip   float64
}

func (stat *CPU) Run() error {
	zero, err := os.Open("/dev/zero")
	defer zero.Close()
	if err != nil {
		return err
	}

	null, err := os.Create("/dev/null")
	defer null.Close()
	if err != nil {
		return err
	}

	hashStart := time.Now()
	hash := sha256.New()
	if _, err := io.CopyN(hash, zero, bytefmt.GIGABYTE); err != nil {
		return err
	}
	hash.Sum(nil)
	stat.Sha256 = time.Since(hashStart).Seconds()

	gzipStart := time.Now()
	gz := gzip.NewWriter(null)
	if _, err := io.CopyN(gz, zero, bytefmt.GIGABYTE); err != nil {
		return err
	}
	gz.Close()
	stat.Gzip = time.Since(gzipStart).Seconds()

	return nil
}

func (stat *CPU) Print() {
	fmt.Printf("Sha256  : %.2f seconds\n", stat.Sha256)
	fmt.Printf("Gzip    : %.2f seconds\n", stat.Gzip)
}
