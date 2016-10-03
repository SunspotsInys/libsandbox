[![Go Report Card](https://goreportcard.com/badge/github.com/ggaaooppeenngg/libsandbox)](https://goreportcard.com/report/github.com/ggaaooppeenngg/libsandbox)
[![Coverage Status](https://coveralls.io/repos/github/ggaaooppeenngg/libsandbox/badge.svg?branch=master)](https://coveralls.io/github/ggaaooppeenngg/libsandbox?branch=master)
[![GoDoc](https://godoc.org/github.com/ggaaooppeenngg/libsandbox?status.svg)](https://godoc.org/github.com/ggaaooppeenngg/libsandbox)
[![Build Status](https://drone.io/github.com/ggaaooppeenngg/libsandbox/status.png)](https://drone.io/github.com/ggaaooppeenngg/libsandbox/latest)

#sandbox
---

Sandbox for online judge.

##install

*command*

```
$ go get github.com/ggaaooppeenngg/libsandbox

```

###Usage:

```
// compile before running and specify limit
sandbox --lang=c -c -s src/main.c -b bin/main --memory=10000 --time=1000 --input=judge/input --output==judge/output
// running with compiled binary file
sandbox --lang=c -b bin/main -i judge/input -o judge/output
```
