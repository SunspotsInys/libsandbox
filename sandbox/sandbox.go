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
	app.Usage = "test untrused source code,result format 'status:time:memory'"
	app.Author = "ggaaooppeenngg"
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "lang", Value: "c,cpp,go", Usage: "source code languge"},
		cli.IntFlag{Name: "time", Value: 1000, Usage: "time limit in MS"},
		cli.IntFlag{Name: "memory", Value: 10000, Usage: "memory limit in KB"},
	}
	app.Action = func(c *cli.Context) {
		if len(c.Args()) >= 2 {
			time := int64(c.Int("time"))
			memory := int64(c.Int("memory"))
			var err error
			pwd, err := os.Getwd()
			if err != nil {
				panic(err)
			}
			src := path.Join(pwd, c.Args()[1])
			var inPath string
			var in *os.File
			if len(c.Args()) >= 3 {
				inPath = path.Join(pwd, c.Args()[2])
				in, err = os.Open(inPath)
			} else {
				in, err = os.Open(os.DevNull)
			}
			if err != nil {
				panic(err)
			}
			defer in.Close()
			var out bytes.Buffer
			var obj = &sandbox.RunningObject{}
			if c.String("lang") == "c" {
				if err = sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.C); err != nil {
					fmt.Printf("CE:0:0")
					return
				} else {
					obj = sandbox.Run(src, in, &out, []string{"tmp"}, time, memory)
					goto testOutput
				}
			}
			if c.String("lang") == "cpp" {
				if err = sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.CPP); err != nil {
					fmt.Printf("CE:0:0")
					return
				} else {
					obj = sandbox.Run(src, in, &out, []string{"tmp"}, time, memory)
					goto testOutput
				}
			}
			if c.String("lang") == "go" {
				if err = sandbox.Complie(c.Args()[0], c.Args()[1], sandbox.GO); err != nil {
					fmt.Printf("CE:0:0")
					return
				} else {
					obj = sandbox.Run(src, in, &out, []string{"tmp"}, time, memory)
					goto testOutput
				}
			}
			//it's convinient to use goto  in the Action context
		testOutput:
			if len(c.Args()) >= 4 {
				outPath := path.Join(pwd, c.Args()[3])
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
		} else {
			println("miss input source file and output destination,please use -h to help.")
		}
	}
	app.Run(os.Args)
}
