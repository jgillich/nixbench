// +build !linux

package main

func Geekbench() (*GeekbenchStat, error) {
	return &GeekbenchStat{
		SingleCore: 0,
		MultiCore:  0,
		ID:         0,
		URL:        "",
	}, nil
}
