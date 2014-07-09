package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/ggaaooppeenngg/sandbox"
)

func checkStatus(obj *sandbox.RunningObject) {
	switch obj.Status {
	case sandbox.AC:
		fmt.Printf("AC")
	case sandbox.MLE:
		fmt.Printf("MLE")
	case sandbox.TLE:
		fmt.Printf("TLE")
	case sandbox.WA:
		fmt.Printf("WA")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "sandbox"
	app.Usage = "test untrused source code"
	app.Author = "ggaaooppeenngg"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{"lang", "c,cpp,go", "source code languge"},
		cli.IntFlag{"time", 1000, "time limit in MS"},
		cli.IntFlag{"memory", 10000, "memory limit in KB"},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 2 {
			time := int64(c.Int("time"))
			memory := int64(c.Int("memory"))
			if c.String("lang") == "c" {
				if err := sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.C); err != nil {
					fmt.Printf("CE")
					os.Exit(2)
				} else {
					obj := sandbox.Run(c.Args()[1], []string{"tmp"}, time, memory)
					checkStatus(obj)
				}
			}
			if c.String("lang") == "cpp" {
				if err := sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.CPP); err != nil {
					fmt.Printf("CE")
					os.Exit(2)
				} else {
					obj := sandbox.Run(c.Args()[1], []string{"tmp"}, time, memory)
					checkStatus(obj)
				}
			}
			if c.String("lang") == "go" {
				if err := sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.GO); err != nil {
					fmt.Printf("CE")
					os.Exit(2)
				} else {
					obj := sandbox.Run(c.Args()[1], []string{"tmp"}, time, memory)
					checkStatus(obj)
				}
			}
		} else {
			println("miss input source file and output destination")
			os.Exit(1)
		}
	}
	app.Run(os.Args)
}
