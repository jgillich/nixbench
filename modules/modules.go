package modules

import "time"

var Modules map[string]Module = map[string]Module{}

type Module interface {
	Run() error
	Print()
	//Name() string
}

func duration(f func() error) (float64, error) {
	start := time.Now()
	err := f()
	return time.Since(start).Seconds(), err
}
