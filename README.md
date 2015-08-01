#sandbox
---

Sandbox for online judge.

##install

*command*

```
$ go get github.com/ggaaooppeenngg/sandbox

//to install sandbox cmd without Docker
$ cd $GOPATH/github.com/ggaaooppeenngg/sandbox/sandbox
$ go install

```

###Usage:

```
    //example:
    //compile before running
    sandbox --lang=c -c -s src/main.c -b bin/main --memory=10000 --time=1000 --input=judge/input --output==judge/output
    //running without compile
    sandbox --lang=c -b bin/main -i judge/input -o judge/output
    //if input or output not set, use /dev/null instead
    sandbox --lang=c -b bin/main 
    //result:
    //output fllows the order below, if result is wrong answer, 5th argument will be attached.
    //status:time:memory:times:wrong_answer

```


##How to implement?

Set process traced with ptrace syscall and send signal SIGALRM to the process every time clock.
Check the /proc/[id]/status virtual memory size, start time and the /proc/status uptime to calculate memory and time consumed.

Also check the signal received,if not SIGAARM (like SIGEGV), it should be some runtime error.
If time or memoery exceed or other signal received, top tracing and return error,else accept.
