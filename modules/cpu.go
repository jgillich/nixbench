package modules

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"runtime"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

func init() {
	Modules["cpu"] = &CPU{}
}

type CPU struct {
	SingleThread CPUStat
	MultiThread  CPUStat
}

type CPUStat struct {
	Threads int
	Sha256  float64
	Gzip    float64
	Aes     float64
}

type DevZero int

func (z DevZero) Read(b []byte) (n int, err error) {
	for i := range b {
		b[i] = 0
	}

	return len(b), nil
}

func (z DevZero) Write(b []byte) (n int, err error) {
	return len(b), nil
}

func (stat *CPU) Run() error {
	if err := stat.SingleThread.runThread(); err != nil {
		return err
	}

	stat.SingleThread.Threads = 1

	stats := make(chan CPUStat)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			stat := CPUStat{}
			if err := stat.runThread(); err != nil {
				panic(err) // TODO
			}
			stats <- stat
		}()
	}

	for i := 0; i < runtime.NumCPU(); i++ {
		s := <-stats
		stat.MultiThread.Sha256 += s.Sha256
		stat.MultiThread.Gzip += s.Gzip
		stat.MultiThread.Aes += s.Aes
	}

	stat.MultiThread.Threads = runtime.NumCPU()

	return nil
}

func (stat *CPU) Print() {
	fmt.Printf("Sha256 (%vx) : %7.2f MB/s\n", stat.SingleThread.Threads, stat.SingleThread.Sha256)
	fmt.Printf("Gzip (%vx)   : %7.2f MB/s\n", stat.SingleThread.Threads, stat.SingleThread.Gzip)
	fmt.Printf("AES (%vx)    : %7.2f MB/s\n", stat.SingleThread.Threads, stat.SingleThread.Aes)
	fmt.Printf("\n")
	fmt.Printf("Sha256 (%vx) : %7.2f MB/s\n", stat.MultiThread.Threads, stat.MultiThread.Sha256)
	fmt.Printf("Gzip (%vx)   : %7.2f MB/s\n", stat.MultiThread.Threads, stat.MultiThread.Gzip)
	fmt.Printf("AES (%vx)    : %7.2f MB/s\n", stat.MultiThread.Threads, stat.MultiThread.Aes)
}

func (stat *CPUStat) runThread() error {
	if res, err := stat.sha256(); err == nil {
		stat.Sha256 = res
	} else {
		return err
	}

	if res, err := stat.gzip(); err == nil {
		stat.Gzip = res
	} else {
		return err
	}

	if res, err := stat.aes(); err == nil {
		stat.Aes = res
	} else {
		return err
	}

	return nil
}

func (stat *CPUStat) sha256() (float64, error) {
	start := time.Now()
	hash := sha256.New()
	if _, err := io.CopyN(hash, DevZero(0), bytefmt.GIGABYTE); err != nil {
		return 0, err
	}
	hash.Sum(nil)
	return (bytefmt.GIGABYTE / time.Since(start).Seconds()) / bytefmt.MEGABYTE, nil
}

func (stat *CPUStat) gzip() (float64, error) {
	start := time.Now()
	gz := gzip.NewWriter(DevZero(0))
	if _, err := io.CopyN(gz, DevZero(0), bytefmt.GIGABYTE); err != nil {
		return 0, err
	}
	return (bytefmt.GIGABYTE / time.Since(start).Seconds()) / bytefmt.MEGABYTE, nil
}

// crypto sourced from https://github.com/gtank/cryptopasta
func (stat *CPUStat) aes() (float64, error) {

	data := [bytefmt.MEGABYTE]byte{}
	if _, err := io.ReadFull(rand.Reader, data[:]); err != nil {
		return 0, err
	}

	key := [32]byte{}
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		return 0, err
	}

	aesStart := time.Now()
	for i := 0; i < bytefmt.GIGABYTE; i += bytefmt.MEGABYTE {

		block, err := aes.NewCipher(key[:])
		if err != nil {
			return 0, err
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return 0, err
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return 0, err
		}

		gcm.Seal(nonce, nonce, data[:], nil)
	}

	return (bytefmt.GIGABYTE / time.Since(aesStart).Seconds()) / bytefmt.MEGABYTE, nil
}
