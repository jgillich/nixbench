package modules

import (
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
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
	Aes    float64
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

	stat.aes()

	return nil
}

func (stat *CPU) Print() {
	fmt.Printf("Sha256  : %.2f seconds\n", stat.Sha256)
	fmt.Printf("Gzip    : %.2f seconds\n", stat.Gzip)
	fmt.Printf("AES     : %.2f seconds\n", stat.Aes)
}

// aes encrypts 100MB of random data using AES-256-GCM
// crypto sourced from https://github.com/gtank/cryptopasta
func (stat *CPU) aes() error {

	data := [bytefmt.MEGABYTE]byte{}
	if _, err := io.ReadFull(rand.Reader, data[:]); err != nil {
		return err
	}

	key := [32]byte{}
	if _, err := io.ReadFull(rand.Reader, key[:]); err != nil {
		return err
	}

	aesStart := time.Now()
	for i := 0; i < bytefmt.GIGABYTE*5; i += bytefmt.MEGABYTE {

		block, err := aes.NewCipher(key[:])
		if err != nil {
			return err
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return err
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return err
		}

		gcm.Seal(nonce, nonce, data[:], nil)
	}
	stat.Aes = time.Since(aesStart).Seconds()

	return nil
}
