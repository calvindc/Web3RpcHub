package main

import (
	"fmt"
	"time"
)

func main() {
	startT := time.Now().UnixNano() / 1e6
	go waitDisp(time.Millisecond * 50)
	const n = 50
	fResult := fibonc(n)
	fmt.Println(fmt.Sprintf("\nFibonacci(%d) = %d ,use time = %.2f ms", n, fResult, float64(time.Now().UnixNano()/1e6-startT)))

}
func fibonc(x int) int {
	if x < 2 {
		return x
	}
	return fibonc(x-1) + fibonc(x-2)
}

func waitDisp(td time.Duration) {
	disWords := [...]string{"-", "\\", "|", "/"}
	for i := 1; true; i++ {
		time.Sleep(td)
		fmt.Printf("\r%s %d ms", disWords[i%4], int64(i)*td.Milliseconds())
	}
}
