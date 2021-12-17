package main

import (
	"fmt"
	"math/big"
	"runtime"
	"strconv"

	"github.com/remeh/sizedwaitgroup"
)

const startNPrime = 2445

func main() {
	//Starting n=X
	var x int64 = 0

	//Wait group with cpu threds
	swg := sizedwaitgroup.New(runtime.NumCPU())
	fmt.Println("Starting", runtime.NumCPU(), "threads.")

	fmt.Print("Checking for n=x primes: ")
	for x = startNPrime; x < 9223372036854775807; x++ {
		swg.Add()
		fmt.Print("n=", x, "?, ")
		go func(val int64) {

			buf := ""
			var y int64 = 0
			// Count up
			for y = 1; y < val; y++ {
				buf = buf + strconv.FormatInt(y, 10)
			}
			//Count down
			for y = y - 1; y > 0; y-- {
				buf = buf + strconv.FormatInt(y, 10)
			}

			temp := big.NewInt(0)
			temp.SetString(buf, 10)

			if temp.ProbablyPrime(20) {
				fmt.Println("POSSIBLE PRIME: n=", val)
				isPrime(val, temp)
			}
			swg.Done()
		}(x)
	}
	swg.Wait()
}

func isPrime(x int64, num *big.Int) bool {
	i := big.NewInt(0)
	iSq := big.NewInt(0)
	iSq = iSq.Sqrt(num)

	for i.SetString("2", 10); i.Cmp(iSq) == -1; i.Add(i, big.NewInt(1)) {
		if num.Mod(num, i) == big.NewInt(0) {
			fmt.Println("n=", x, " divisible by", i, ", ")
			return false
		}
	}
	fmt.Println("\n***** n=", x, " is prime! *****")
	return true
}
