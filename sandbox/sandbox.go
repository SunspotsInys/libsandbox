package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/ggaaooppeenngg/sandbox"
)

func checkStatus(obj *sandbox.RunningObject) {
	if obj.Status == sandbox.AC {
		os.Exit(0)
	} else {
		os.Exit(3)
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
			//			println(time)
			memory := int64(c.Int("memory"))
			if c.String("lang") == "c" {
				if err := sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.C); err != nil {
					println(err.Error())
					os.Exit(2)
				} else {
					obj := sandbox.Run(c.Args()[1], []string{"tmp"}, time, memory)
					checkStatus(obj)
				}
			}
			if c.String("lang") == "cpp" {
				if err := sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.CPP); err != nil {
					println(err.Error())
					os.Exit(2)
				} else {
					obj := sandbox.Run(c.Args()[1], []string{"tmp"}, time, memory)
					checkStatus(obj)
				}
			}
			if c.String("lang") == "go" {
				if err := sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.GO); err != nil {
					println(err.Error())
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
