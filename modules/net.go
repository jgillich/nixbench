package modules

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"code.cloudfoundry.org/bytefmt"
)

func init() {
	Modules["net"] = &Net{}
}

type Net map[string]float64

type file struct {
	Key string
	URL string
}

var files = []file{
	{Key: "CDN", URL: "http://cachefly.cachefly.net/100mb.test"},
	{Key: "Amsterdam, The Netherlands", URL: "http://speedtest.ams01.softlayer.com/downloads/test100.zip"},
	{Key: "Dallas, USA", URL: "http://speedtest.dal05.softlayer.com/downloads/test100.zip"},
	{Key: "Frankfurt, Germany", URL: "http://speedtest.fra02.softlayer.com/downloads/test100.zip"},
	{Key: "Hong Kong, China", URL: "http://speedtest.hkg02.softlayer.com/downloads/test100.zip"},
	{Key: "London, United Kingdoms", URL: "http://speedtest.lon02.softlayer.com/downloads/test100.zip"},
	{Key: "Melbourne, Australia", URL: "http://speedtest.mel01.softlayer.com/downloads/test100.zip"},
	{Key: "Oslo, Norway", URL: "http://speedtest.osl01.softlayer.com/downloads/test100.zip"},
	{Key: "Paris, France", URL: "http://speedtest.par01.softlayer.com/downloads/test100.zip"},
	{Key: "Queretaro Mexico", URL: "http://speedtest.mex01.softlayer.com/downloads/test100.zip"},
	{Key: "San Jose, USA", URL: "http://speedtest.sjc01.softlayer.com/downloads/test100.zip"},
	{Key: "Sao Paulo, Brazil", URL: "http://speedtest.sao01.softlayer.com/downloads/test100.zip"},
	{Key: "Seoul, Korea", URL: "http://speedtest.seo01.softlayer.com/downloads/test100.zip"},
	{Key: "Singapore, Singapore", URL: "http://speedtest.sng01.softlayer.com/downloads/test100.zip"},
	{Key: "Tokyo, Japan", URL: "http://speedtest.tok02.softlayer.com/downloads/test100.zip"},
	{Key: "Toronto, Canada", URL: "http://speedtest.tor01.softlayer.com/downloads/test100.zip"},
	{Key: "Washington, USA", URL: "http://speedtest.wdc01.softlayer.com/downloads/test100.zip"},
}

func (stat *Net) Run() error {
	var client = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 5 * time.Second,
		},
		Timeout: time.Minute * 5,
	}

	for _, file := range files {
		res, err := client.Get(file.URL)
		if err != nil {
			(*stat)[file.Key] = 0
			continue
		}

		if res.ContentLength < bytefmt.MEGABYTE*100 {
			return fmt.Errorf("unexpected content length %v at %s", res.ContentLength, file.Key)
		}

		start := time.Now()
		_, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		duration := time.Since(start)

		(*stat)[file.Key] = float64(res.ContentLength/bytefmt.MEGABYTE) / duration.Seconds()
	}

	return nil
}

func (stat *Net) Print() {
	for _, f := range files {
		if (*stat)[f.Key] == 0 {
			fmt.Printf("%-30s: Failed\n", f.Key)
		} else {
			fmt.Printf("%-30s: %-6.2f MB/s\n", f.Key, (*stat)[f.Key])
		}
	}
}
