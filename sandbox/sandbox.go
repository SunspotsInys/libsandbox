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

const (
	BINARY  = "binary"
	COMPILE = "compile"
	SOURCE  = "source"
	TIME    = "time"
	MEMORY  = "memory"
	INPUT   = "input"
	OUTPUT  = "output"
	LANG    = "lang"

	C   = "c"
	CPP = "cpp"
	GO  = "go"

	DELIM = "!-_-\n"
)

func panicErr(e error) {
	if e != nil {
		fmt.Println(e)
	}
}

//read byte from the file
func readFile(f *os.File) (testOut []byte) {
	tmp := make([]byte, 256)
	for n, err := f.Read(tmp); err != io.EOF; n, err = f.Read(tmp) {
		testOut = append(testOut, tmp[:n]...)
	}
	return testOut
}

//obj records process information and n is the nth test,if n is 0 ,all test are passed
func checkStatus(obj *sandbox.RunningObject, n int) {
	switch obj.Status {
	case sandbox.AC:
		fmt.Printf("AC:%d:%d:%d", obj.Memory, obj.Time, n)
	case sandbox.MLE:
		fmt.Printf("MLE:%d:%d:%d", obj.Memory, obj.Time, n)
	case sandbox.TLE:
		fmt.Printf("TLE:%d:%d:%d", obj.Memory, obj.Time, n)
	default:
		fmt.Printf("FE:%d:%d:%d", obj.Memory, obj.Time, n)
	}
}

func main() {

	app := cli.NewApp()
	app.Name = "sandbox"
	app.Usage = `test untrusted source code'
example:
	sandbox --lang=c -c -s src/main.c -b bin/main --memory=10000 --time=1000 --input=judge/input --output==judge/output
	sandbox --lang=c -b bin/main -i judge/input -o judge/output
	echo "input" | sandbox --lang=c  bin/main 
result:
	status:time:memory:times`
	app.Author = "ggaaooppeenngg"
	app.Version = "0.0.2"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "lang,l", Value: "c,cpp,go", Usage: "source code languge"},
		cli.IntFlag{Name: "time,t", Value: 1000, Usage: "time limit in MS"},
		cli.IntFlag{Name: "memory,m", Value: 10000, Usage: "memory limit in KB"},
		cli.BoolFlag{Name: "compile,c", Usage: "wether complie before running", EnvVar: ""},
		cli.StringFlag{Name: "input,i", Value: "", Usage: "input file path"},
		cli.StringFlag{Name: "output,o", Value: "", Usage: "output file path"},
		cli.StringFlag{Name: "source,s", Value: "", Usage: "source file path"},
		cli.StringFlag{Name: "binary,b", Value: "", Usage: "binary file path"},
	}
	app.Action = func(c *cli.Context) {
		var in *os.File  //input file instance
		var out *os.File //output file instance
		var src string   // source file path
		var bin string   //binary file path
		var err error
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		if c.String(LANG) == "" {
			println("needs to specify a language")
			return
		}
		//target binary file path is neccessary
		if c.String(BINARY) != "" {
			p := c.String(BINARY)
			if path.IsAbs(p) {
				bin = p
			} else {
				bin = path.Join(pwd, p)
			}
		} else {
			println("needs target binary file path as argument")
			return
		}

		//if input is not set , use /dev/null as input
		if c.String(INPUT) == "" {
			in, err = os.Open(os.DevNull)
		} else {
			p := c.String(INPUT)
			if path.IsAbs(p) {
				in, err = os.Open(p)
			} else {
				in, err = os.Open(path.Join(pwd, p))
			}
		}
		if err != nil {
			panic(err)
		}
		defer in.Close()

		if c.Bool(COMPILE) {
			if c.String(SOURCE) == "" {
				println("compiler needs source file!")
				return
			} else {
				//get source file path
				p := c.String(SOURCE)
				if path.IsAbs(p) {
					src = p
				} else {
					src = path.Join(pwd, p)
				}
				//compile code ,if compile set , not compile
				if c.Bool(COMPILE) {
					var language uint64
					switch c.String(LANG) {
					case C:
						language = sandbox.C
					case CPP:
						language = sandbox.CPP
					case GO:
						language = sandbox.GO
					}
					if err = sandbox.Complie(src, bin, language); err != nil {
						fmt.Printf("CE:0:0")
						return
					}
				}
			}
		}

		var obj = &sandbox.RunningObject{}
		time := int64(c.Int(TIME))
		memory := int64(c.Int(MEMORY))
		if c.String(OUTPUT) != "" {
			//get out test and check if every output matches the single input
			outPath := c.String(OUTPUT)
			if !path.IsAbs(outPath) {
				outPath = path.Join(pwd, outPath)
			}
			out, err = os.Open(outPath)
		} else {
			out, err = os.Open(os.DevNull)

		}
		if err != nil {
			panic(err)
		}
		defer out.Close()
		//form a  scope
		if c.String(OUTPUT) != "" {

			//get input tests and run every test one by one
			i := readFile(in)
			inputs := bytes.Split(i, []byte(DELIM))
			o := readFile(out)
			outputs := bytes.Split(o, []byte(DELIM))
			for i, v := range inputs {
				inBytes := bytes.NewBuffer(v)
				out := bytes.Buffer{}
				obj = sandbox.Run(bin, inBytes, &out, []string{""}, time, memory)
				if !bytes.Equal(out.Bytes(), outputs[i]) {
					fmt.Printf("WA:%d:%d:%d", obj.Memory, obj.Time, i+1)
					return
				}
			}
		} else {
			out := bytes.Buffer{}
			in := bytes.NewBuffer([]byte{})
			obj = sandbox.Run(bin, in, &out, []string{""}, time, memory)
		}
		checkStatus(obj, 0)
		return
	}
	app.Run(os.Args)
}
