package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/codegangsta/cli"
	"github.com/ggaaooppeenngg/sandbox"
)

func checkStatus(obj *sandbox.RunningObject) {
	switch obj.Status {
	case sandbox.AC:
		fmt.Printf("AC:%d:%d", obj.Memory, obj.Time)
	case sandbox.MLE:
		fmt.Printf("MLE:%d:%d", obj.Memory, obj.Time)
	case sandbox.TLE:
		fmt.Printf("TLE:%d:%d", obj.Memory, obj.Time)
	case sandbox.WA:
		fmt.Printf("WA:%d:%d", obj.Memory, obj.Time)
	default:
		fmt.Printf("FE:%d:%d", obj.Memory, obj.Time)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "sandbox"
	app.Usage = `test untrusted source code'
example:
	sandbox --lang=c src/main.c bin/main judge/input judge/output
result:
	status:time:memory`
	app.Author = "ggaaooppeenngg"
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "lang,l", Value: "c,cpp,go", Usage: "source code languge"},
		cli.IntFlag{Name: "time,t", Value: 1000, Usage: "time limit in MS"},
		cli.IntFlag{Name: "memory,m", Value: 10000, Usage: "memory limit in KB"},
		cli.BoolFlag{Name: "compile,c", Usage: "wether complie before running", EnvVar: ""},
	}
	app.Action = func(c *cli.Context) {
		var bin string
		var inPath string
		var outPath string
		var err error
		var in *os.File
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if c.Bool("compile") {
			if len(c.Args()) < 1 {
				println("not match the number of arguments,please use -h to help.")
				return
			} else {
				bin = path.Join(pwd, c.Args()[1])
				if len(c.Args()) >= 4 {
					outPath = path.Join(pwd, c.Args()[3])
				}
				if len(c.Args()) >= 3 {
					inPath = path.Join(pwd, c.Args()[2])
					in, err = os.Open(inPath)
				} else {
					in, err = os.Open(os.DevNull)
				}
			}
		} else {
			if len(c.Args()) < 2 {
				println("not match the number of arguments,please use -h to help.")
				return
			} else {
				bin = path.Join(pwd, c.Args()[0])
				if len(c.Args()) >= 3 {
					outPath = path.Join(pwd, c.Args()[2])
				}
				if len(c.Args()) >= 2 {
					inPath = path.Join(pwd, c.Args()[1])
					in, err = os.Open(inPath)
				} else {
					in, err = os.Open(os.DevNull)
				}
			}
		}
		time := int64(c.Int("time"))
		memory := int64(c.Int("memory"))

		if err != nil {
			panic(err)
		}
		defer in.Close()
		var out bytes.Buffer
		var obj = &sandbox.RunningObject{}

		//compile code ,if compile set , not compile
		if c.Bool("compile") {
			if c.String("lang") == "c" {
				if err = sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.C); err != nil {
					fmt.Printf("CE:0:0")
					return
				}
			}
			if c.String("lang") == "cpp" {
				if err = sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.CPP); err != nil {
					fmt.Printf("CE:0:0")
					return
				}
			}
			if c.String("lang") == "go" {
				if err = sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.GO); err != nil {
					fmt.Printf("CE:0:0")
					return
				}
			}
		}
		//fmt.Println(bin, inPath, outPath)
		obj = sandbox.Run(bin, in, &out, []string{"tmp"}, time, memory)
		if outPath != "" {
			//get output
			outFile, err := os.Open(outPath)
			defer outFile.Close()
			if err != nil {
				panic(err)
			}
			var testOut []byte
			tmp := make([]byte, 256)
			for n, err := outFile.Read(tmp); err != io.EOF; n, err = outFile.Read(tmp) {
				testOut = append(testOut, tmp[:n]...)
			}
			if bytes.Equal(out.Bytes(), testOut) {
				fmt.Printf("AC:%d:%d", obj.Memory, obj.Time)
				return
			}
		}
		checkStatus(obj)
	}
	app.Run(os.Args)
}
