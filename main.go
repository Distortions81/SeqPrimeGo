package main

import (
	"fmt"
	"math/big"
	"runtime"
	"strconv"

	"github.com/remeh/sizedwaitgroup"
)

func main() {
	var x int64 = 0

	swg := sizedwaitgroup.New(runtime.NumCPU())
	fmt.Println("Starting", runtime.NumCPU(), "threads.")

	for x = 1; x < 9223372036854775807; x++ {
		swg.Add()
		go func(val int64) {

			buf := ""
			var y int64 = 0
			// Count up
			for y = 1; y < x; y++ {
				buf = buf + strconv.FormatInt(y, 10)
			}
			//Count down
			for y = y - 1; y > 0; y-- {
				buf = buf + strconv.FormatInt(y, 10)
			}

			temp := big.NewInt(0)
			temp.SetString(buf, 10)

			if isPrime(temp) {
				fmt.Print(val, ", ")

			} else {
				//fmt.Print("!N=", val, ", ")
			}
			swg.Done()
		}(x)
	}
	swg.Wait()
}

func isPrime(num *big.Int) bool {
	i := big.NewInt(0)
	iSq := big.NewInt(0)
	iSq = iSq.Sqrt(num)

	for i.SetString("2", 10); i.Cmp(iSq) == -1; i.Add(i, big.NewInt(1)) {
		if num.Mod(num, i) == big.NewInt(0) {
			return false
		}
	}
	return true
}
