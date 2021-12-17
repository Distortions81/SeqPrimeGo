package main

import (
	"fmt"
	"math/big"
	"runtime"

	"github.com/remeh/sizedwaitgroup"
)

func main() {
	x := big.NewInt(0)

	swg := sizedwaitgroup.New(runtime.NumCPU() * 2)
	fmt.Println("Starting", runtime.NumCPU(), "threads...")

	for x.SetString("1", 0); true; x.Add(x, big.NewInt(1)) {
		swg.Add()
		go func() {
			buf := ""
			y := big.NewInt(0)
			for y.SetString("1", 0); y.Cmp(x) == -1; y.Add(y, big.NewInt(1)) {
				buf = buf + y.String()
			}
			temp := big.NewInt(0)
			temp.SetString(buf, 0)

			if temp.ProbablyPrime(0) {
				fmt.Print(temp, ", ")
			} else {
				//fmt.Print("NOT: ", temp, ", ")
			}
			swg.Done()
		}()
		swg.Wait()
	}
}
