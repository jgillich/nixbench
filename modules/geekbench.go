// +build linux

package modules

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

func init() {
	Modules["geekbench"] = &Geekbench{}
}

type Geekbench struct {
	SingleCore int
	MultiCore  int
	URL        string
	ID         int
}

func (stat *Geekbench) Run() error {
	res, err := http.Get("http://cdn.primatelabs.com/Geekbench-4.1.0-Linux.tar.gz")
	if err != nil {
		return err
	}

	if err := extract(res.Body); err != nil {
		return err
	}

	gb := exec.Command("build.pulse/dist/Geekbench-4.1.0-Linux/geekbench4")

	out, err := gb.Output()
	if err != nil {
		fmt.Println(string(out))
		return err
	}

	if err := os.RemoveAll("build.pulse"); err != nil {
		return err
	}

	r, err := regexp.Compile("https://browser.geekbench.com/v4/cpu/([0-9]*)")
	if err != nil {
		return err
	}

	match := r.FindStringSubmatch(string(out))
	if len(match) < 2 {
		return errors.New("geekbench did not return result url")
	}

	stat.URL = match[0]
	id, _ := strconv.Atoi(match[1])
	stat.ID = id

	stat.SingleCore, stat.MultiCore, err = scrape(stat.URL)
	if err != nil {
		return err
	}

	return nil
}

func (stat *Geekbench) Print() {
	fmt.Printf("Single-Core Score  : %d\n", stat.SingleCore)
	fmt.Printf("Multi-Core Score   : %d\n", stat.MultiCore)
	fmt.Printf("Result URL         : %s\n", stat.URL)
}

func scrape(url string) (int, int, error) {
	res, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, 0, err
	}

	singleReg, err := regexp.Compile("<th class='name'>Single-Core Score</th>\n<th class='score'>([0-9]*)</th>")
	if err != nil {
		return 0, 0, err
	}

	single := singleReg.FindStringSubmatch(string(body))
	if len(single) < 2 {
		return 0, 0, errors.New("failed to scrape score")
	}

	multiReg, err := regexp.Compile("<th class='name'>Multi-Core Score</th>\n<th class='score'>([0-9]*)</th>")
	if err != nil {
		return 0, 0, err
	}

	multi := multiReg.FindStringSubmatch(string(body))
	if len(multi) < 2 {
		return 0, 0, errors.New("failed to scrape score")
	}

	singleScore, _ := strconv.Atoi(single[1])
	multiScore, _ := strconv.Atoi(multi[1])

	return singleScore, multiScore, nil
}

func extract(r io.Reader) error {

	gzf, err := gzip.NewReader(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tarReader := tar.NewReader(gzf)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(header.Name, os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg:
			f, err := os.Create(header.Name)
			if err != nil {
				return err
			}
			defer f.Close()

			if _, err := io.Copy(f, tarReader); err != nil {
				return err
			}

			os.Chmod(header.Name, os.ModePerm)
		default:
			return fmt.Errorf("%s %c %s %s", "unexpected file type", header.Typeflag, "in file", header.Name)
		}
	}

}
