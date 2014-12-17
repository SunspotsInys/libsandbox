#sandbox
---

sandbox in Golang, for online judge.

##build dockerfile:

```
docker build -t ubuntu/sandbox .

```

##usage: 

see tracer\_test.go and other test files for Public funtion usage

*command*

```
$ cd $GOPATH/github.com/ggaaooppeenngg/sandbox/sandbox
$ go install

```
###Usage:
```
sandbox -h

```

###with docker

```
//TODO:
1. produce more detail status(runtime error, presentation error, outputlimit error)

2. wrap with docker util the sandbox is stable

3. check out illegal system call to feedback errors

```

##test:

run `./test.sh`

##todo

change syscall package to apply to Go1.4
