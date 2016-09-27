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
	Status string // AC
	Memory int    // KB
	Time   int    // MS
	Nth    int
}

func (r Result) Json() (string, error) {
	v, err := json.Marshal(r)
	return string(v), err
}

func main() {
	var (
		source   string
		binary   string
		language string
		input    string
		output   string
		memory   int64
		time     int64
	)
	app := cli.NewApp()
	app.Author = "ggaaooppeenngg"
	app.Version = "0.0.3"
	app.Email = "peng.gao.dut@gmail.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "lang,l", Usage: "source code languge", Destination: &language},
		cli.Int64Flag{Name: "time,t", Value: 1000, Usage: "time limit in MS", Destination: &time},
		cli.Int64Flag{Name: "memory,m", Value: 1000 * 1024, Usage: "memory limit in KB", Destination: &memory},
		cli.BoolFlag{Name: "compile,c", Usage: "wether complie before running"},
		cli.StringFlag{Name: "input,i", Value: "/dev/null", Usage: "input file path", Destination: &input},
		cli.StringFlag{Name: "output,o", Value: "/dev/null", Usage: "output file path", Destination: &output},
		cli.StringFlag{Name: "source,s", Usage: "source file path", Destination: &source},
		cli.StringFlag{Name: "binary,b", Usage: "binary file path", Destination: &binary},
	}
	app.Name = "sandbox"
	app.Usage = ""
	app.Version = "0.2"
	// libcontainer的边界应该是 binary-> sandbox -> output
	app.Action = func(c *cli.Context) error {
		id := uuid.NewV1()
		defaultbinpath := filepath.Join(os.TempDir(), id.String()+".binary")
		if binary == "" {
			binary = defaultbinpath
		}
		for _, pair := range []struct {
			K string
			V string
		}{
			{K: "language", V: language},
			{K: "source", V: source},
		} {
			if pair.V == "" {
				return cli.NewExitError(fmt.Sprintf("%s is empty", pair.K), 1)
			}
		}
		err := compile(source, binary, language)
		if err != nil {
			return err
		}
		input, err := os.Open(input)
		if err != nil {
			return err
		}
		defer input.Close()
		outputF, err := os.Open(output)
		if err != nil {
			return err
		}
		defer outputF.Close()
		sandbox, err := libsandbox.NewStdSandbox(libsandbox.Config{
			Args:   []string{binary},
			Input:  input,
			Time:   time,
			Memory: memory,
		})
		if err != nil {
			return err
		}
		outputStd, err := ioutil.ReadAll(outputF)
		if err != nil {
			return err
		}
		output, err := sandbox.Run() // 这个output 应该是已经申明过了才对?...
		if err != nil {
			return err
		}
		if bytes.Equal(output, outputStd) {
			res, err := Result{
				Status: "AC",
			}.Json()
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
		}
		output = bytes.Map(func(r rune) rune {
			if r == '\t' || r == '\n' || r == ' ' {
				return -1
			}
			return r
		}, output)
		outputStd = bytes.Map(func(r rune) rune {
			if r == '\t' || r == '\n' || r == ' ' {
				return -1
			}
			return r
		}, outputStd)

		if bytes.Equal(output, outputStd) {
			res, err := Result{
				Status: "FE",
			}.Json()
			if err != nil {
				return err
			}
			fmt.Println(res)
			return nil
		}
		return nil
	}
	app.Run(os.Args)
}
