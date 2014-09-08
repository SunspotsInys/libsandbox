package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/codegangsta/cli"
)

const (
	INPUT    = "input"
	OUTPUT   = "output"
	STANDARD = "standard"
	DELIM    = "!-_-\n"
)

//read byte from the file
func readFile(f *os.File) (testOut []byte) {
	tmp := make([]byte, 256)
	for n, err := f.Read(tmp); err != io.EOF; n, err = f.Read(tmp) {
		testOut = append(testOut, tmp[:n]...)
	}
	return testOut
}

// SaveFile saves content type '[]byte' to file by given path.
// It returns error when fail to finish operation.
func writeFile(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}
func getPath(c *cli.Context, t string, wd string) (p string, e error) {
	if c.String(t) != "" {
		p = c.String(t)
		if !path.IsAbs(p) {
			p = path.Join(wd, p)
		}
		return p, nil
	} else {
		return p, errors.New(t + " is needed")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "generator"
	app.Usage = `generate ouput with input and standard progra
example:
	generate -i input -s standard -o outputpath`
	app.Author = "ggaaooppeenngg"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "input,i", Value: ""},
		cli.StringFlag{Name: "standard,s", Value: ""},
		cli.StringFlag{Name: "output,o", Value: ""},
	}
	app.Action = func(c *cli.Context) {
		wd, e := os.Getwd()
		if e != nil {
			panic(e)
		}
		var i string //input file path
		var s string //standard program path
		var o string //output path
		i, e = getPath(c, INPUT, wd)
		if e != nil {
			fmt.Println(e)
			return
		}
		s, e = getPath(c, STANDARD, wd)
		if e != nil {
			fmt.Println(e)
			return
		}
		o, e = getPath(c, OUTPUT, wd)
		if e != nil {
			fmt.Println(e)
			return
		}

		input, e := os.Open(i)
		if e != nil {
			panic(e)
		}
		defer input.Close()
		bs := readFile(input)
		inputs := bytes.Split(bs, []byte(DELIM))
		var outputs [][]byte
		for _, v := range inputs {
			input := bytes.NewBuffer(v)
			cmd := exec.Command(s)
			cmd.Stdin = input
			o, e := cmd.CombinedOutput()
			if e != nil {
				//fmt.Printf("%s", o)
				//panic(e)
			}
			outputs = append(outputs, o)
		}
		_, e = writeFile(o, bytes.Join(outputs, []byte(DELIM)))
		if e != nil {
			panic(e)
		}
	}
	app.Run(os.Args)
}
