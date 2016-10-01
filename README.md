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
