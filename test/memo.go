package main

import (
	"time"
)

func main() {
	var a [10000][]byte
	for i := 0; i < 10000; i++ {
		a[i] = make([]byte, 1024)
		for j := 0; j < 1024; j++ {
			a[i][j] = 1
		}
	}
	time.Sleep(time.Second * 5)
}
