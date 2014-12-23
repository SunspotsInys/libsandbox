//Dockerbox is sandbox wrapped with Docker
package main

import (
	"fmt"
	"os/exec"
	"strings"

	"code.google.com/p/go-uuid/uuid"
	"github.com/ggaaooppeenngg/util"
)

func uuPath(lang string) (id, path string) {
	ui := uuid.NewUUID()
	uid := strings.Replace(ui.String(), "-", "", -1)
	if lang == "go" {
		return uid, uid + ".go"
	} else {
		return "", ""
	}
}

//use command sed repleace SRCFILE to real source file
func genDocFile(path string) error {
	out, err := util.Run("sed", "s/SRCFILE/"+path+"/g", "Seedfile")
	if err != nil {
		fmt.Printf("%s", out)
		return err
	}
	_, err = util.WriteFile("Dockerfile", out)
	if err != nil {
		return err
	}
	return nil
}

//clean the container after running
func removeContainer(name string) {
	util.Run("docker", "rm", name)
}

func test(id, path string) []byte {
	defer removeContainer(id)
	err := genDocFile(path)
	if err != nil {
		panic(err)
	}
	//use Dockerfile to add source code
	cmd := exec.Command("docker", "build", "-t", id, ".")
	outs, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", outs)
		fmt.Println(err)
	}

	//build
	cmd = exec.Command("docker", "run", "-i", "--name="+id, id, "/home/GoPath/bin/sandbox", "--lang=go", "/home/"+path, "/home/"+id)
	outs, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("%s", outs)
		fmt.Println(err)
		return outs
	}

	return outs
}

func main() {
	test("test1234", "test1234.go")
	//fmt.Printf("%s", res)
}
