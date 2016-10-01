//sandbox is command line interface for the Sandbox without docker wrapped.
//  Example:
//      compile before running
//          sandbox --lang=c -c -s src/main.c -b bin/main --memory=10000 --time=1000 --input=judge/input --output==judge/output
//      running without compiling
//          sandbox --lang=c -b bin/main -i judge/input -o judge/output
//      if input or output not set, use /dev/null instead
//          sandbox --lang=c -b bin/main
//      result:
//          output fllows the order below,if result is wrong answer,5th argument will be attached.
//          status:time:memory:times:wrong_answer
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	uuid "github.com/satori/go.uuid"
	"github.com/urfave/cli"

	"github.com/ggaaooppeenngg/libsandbox"
)

// render more information
const (
	C   = "c"
	CPP = "cpp"
	GO  = "go"

	DELIM = "!-_-\n"

	OUTPUT_LIMIT = 255
)

const (
	FormatError = "FE"
)

// default comiple options
func compile(src string, des string, lan string) error {
	var cmd = new(exec.Cmd)
	switch lan {
	case C:
		cmd = exec.Command("gcc", "-o", des, src, "-lm") //-lm for gcc math link option
	case CPP:
		cmd = exec.Command("g++", "-o", des, src)
	case GO:
		cmd = exec.Command("go", "build", "-o", des, src)
	default:
		return fmt.Errorf("unspported or unknown language %s", lan)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return cli.NewExitError(fmt.Sprintf("%s err:%s", out, err), 1)
	} else {
		return nil
	}
}

type Result struct {
	Status    string // sandbox result status
	Error     string //
	Memory    int64  // KB
	Time      int64  // MS
	Nth       int
	Output    string // output
	StdOutput string // wanted output
}

func (r Result) Json() ([]byte, error) {
	return json.Marshal(r)
}

func runAction(c *cli.Context) (Result, error) {
	id := uuid.NewV1()
	defaultbinpath := filepath.Join(os.TempDir(), id.String()+".binary")
	var (
		binary   = c.String("binary")
		source   = c.String("source")
		memory   = c.Int64("memory")
		time     = c.Int64("time")
		language = c.String("lang")
		input    = c.String("input")
		output   = c.String("output")
	)
	if binary == "" {
		binary = defaultbinpath
	}
	for _, pair := range []struct {
		K string
		V string
	}{
		{K: "lang", V: language},
		{K: "source", V: source},
	} {
		if pair.V == "" {
			return Result{}, cli.NewExitError(fmt.Sprintf("%s is empty", pair.K), 1)
		}
	}
	err := compile(source, binary, language)
	if err != nil {
		return Result{}, err
	}

	// read test inputs
	inputF, err := os.Open(input)
	if err != nil {
		return Result{}, err
	}
	defer inputF.Close()
	inputs, err := ioutil.ReadAll(inputF)
	if err != nil {
		return Result{}, err
	}

	// read test outputs
	outputF, err := os.Open(output)
	if err != nil {
		return Result{}, err
	}
	defer outputF.Close()
	outputs, err := ioutil.ReadAll(outputF)
	if err != nil {
		return Result{}, err
	}

	noutputs := bytes.Split(outputs, []byte(DELIM))
	ninputs := bytes.Split(inputs, []byte(DELIM))
	if len(noutputs) != len(ninputs) {
		return Result{}, errors.New("length of test input and output not equal")
	}
	for n, input := range ninputs {
		sandbox, err := libsandbox.NewStdSandbox(libsandbox.Config{
			Args:   []string{binary},
			Input:  bytes.NewReader(input),
			Time:   time,
			Memory: memory,
		})
		if err != nil {
			return Result{}, err
		}
		outputBytes, err := sandbox.Run() // 这个output 应该是已经申明过了才对?...
		if err == libsandbox.OutOfTimeError {
			return Result{
				Status: "TL",
				Error:  err.Error(),
				Nth:    n + 1,
				Time:   sandbox.Time(),
			}, nil
		}
		if err == libsandbox.OutOfMemoryError {
			return Result{
				Status: "ML",
				Error:  err.Error(),
				Nth:    n + 1,
				Memory: sandbox.Memory(),
			}, nil
		}
		if err != nil {
			return Result{
				Status: "RE",
				Error:  err.Error(),
				Nth:    n + 1,
			}, nil
		}
		if bytes.Equal(outputBytes, noutputs[n]) {
			if n != len(ninputs) {
				continue
			}
		}
		outputBytesStd := bytes.Map(func(r rune) rune {
			if r == '\t' || r == '\n' || r == ' ' {
				return -1
			}
			return r
		}, outputBytes)
		outputStd := bytes.Map(func(r rune) rune {
			if r == '\t' || r == '\n' || r == ' ' {
				return -1
			}
			return r
		}, noutputs[n])

		if bytes.Equal(outputBytesStd, outputStd) {
			return Result{
				Status: "FE",
			}, nil
		}
		return Result{
			Status:    "WA",
			Nth:       n + 1,
			Output:    string(outputBytes),
			StdOutput: string(noutputs[n]),
		}, nil
	}
	return Result{
		Status: "AC",
	}, nil
}

func NewApp(action func(c *cli.Context) error) *cli.App {
	app := cli.NewApp()
	app.Author = "ggaaooppeenngg"
	app.Version = "0.0.3"
	app.Email = "peng.gao.dut@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "lang,l", Usage: "source code languge"},
		cli.Int64Flag{Name: "time,t", Value: 1000, Usage: "time limit in MS"},
		cli.Int64Flag{Name: "memory,m", Value: 1000 * 1024, Usage: "memory limit in KB"},
		cli.BoolFlag{Name: "compile,c", Usage: "wether complie before running"},
		cli.StringFlag{Name: "input,i", Value: "/dev/null", Usage: "input file path"},
		cli.StringFlag{Name: "output,o", Value: "/dev/null", Usage: "output file path"},
		cli.StringFlag{Name: "source,s", Usage: "source file path"},
		cli.StringFlag{Name: "binary,b", Usage: "binary file path"},
	}
	app.Name = "sandbox"
	app.Usage = ""
	app.Version = "0.2"
	app.Action = action
	return app
}

func main() {
	// libcontainer的边界应该是 binary-> sandbox -> output
	NewApp(func(c *cli.Context) error {
		result, err := runAction(c)
		if err != nil {
			return err
		}
		v, err := result.Json()
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", v)
		return nil
	}).Run(os.Args)
}
