package main

func main() {
	var a [10000][]int
	for i := 0; i < 10000; i++ {
		a[i] = make([]int, 1024)
	}
}
