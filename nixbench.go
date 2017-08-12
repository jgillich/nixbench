package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	cli "gopkg.in/urfave/cli.v2"
	yaml "gopkg.in/yaml.v2"

	"github.com/jgillich/nixbench/modules"
)

// VERSION is set at build time
var VERSION = "master"

func main() {

	app := &cli.App{
		Name:        "nixbench",
		Usage:       "A better benchmarking tool for servers",
		Description: fmt.Sprintf("Loaded modules: %s", strings.Trim(fmt.Sprintf("%v", reflect.ValueOf(modules.Modules).MapKeys()), "[]")),
		Version:     VERSION,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "yaml",
				Usage: "Output as yaml",
			},
			&cli.StringSliceFlag{
				Name:    "modules",
				Aliases: []string{"m"},
				Usage:   "Modules to enable",
				Value:   cli.NewStringSlice("host", "cpu", "disk", "net"),
			},
		},
		Action: func(c *cli.Context) error {
			if !c.Bool("yaml") {
				fmt.Printf("nixbench %s - https://github.com/jgillich/nixbench", VERSION)
			}

			for _, name := range c.StringSlice("modules") {
				module, ok := modules.Modules[name]

				if !ok {
					return fmt.Errorf("unknown module '%s'", name)
				}

				if err := module.Run(); err != nil {
					return err
				}

				if c.Bool("yaml") {
					var r map[string]interface{} = map[string]interface{}{}
					r[name] = module
					yml, err := yaml.Marshal(r)
					if err != nil {
						return err
					}
					fmt.Printf(string(yml))
				} else {
					fmt.Printf("\n\n%s\n", name)
					for i := 1; i <= len(name); i++ {
						fmt.Print("-")
					}
					fmt.Print("\n")
					module.Print()
				}
			}

			return nil
		},
	}

	app.Run(os.Args)
}
