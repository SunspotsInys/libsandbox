#sandbox
---

sandbox in Golang, for online judge.

##usage: 

see tracer\_test.go and other test files for Public funtion usage

*command*

```
$ cd $GOPATH/github.com/ggaaooppeenngg/sandbox/sandbox
$ go install

sandbox --lang=go main.go main
$?
```
if exit is not 0, it exceeds limit.
//TODO:define relate exit code with specific limitation

##test:

run `./test.sh`

//TODO:rlimit RLIMIT_AS test doesn't signal but exit ,going to deal with it.
