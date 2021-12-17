package main

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"runtime"
	"strconv"

	"github.com/remeh/sizedwaitgroup"
)

const startNPrime = 2

func main() {
	//Logging setup
	logName := "nPrimes.log"
	lf, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer lf.Close()

	mw := io.MultiWriter(os.Stdout, lf)
	log.SetOutput(mw)

	//Init
	var x int64 = startNPrime
	buf := ""

	//Wait group with cpu threds
	swg := sizedwaitgroup.New(runtime.NumCPU())
	log.Println("Starting", runtime.NumCPU(), "threads.")

	log.Println("Creating first big int...")
	// Count up
	var z int64 = 0
	for z = 1; z <= x; z++ {
		buf = buf + strconv.FormatInt(z, 10)
	}
	//log.Println(buf)

	log.Println("Checking for n=x primes: ")
	for x = z; x < 9223372036854775807; x++ {
		swg.Add()
		//log.Print("n=", x, "? ")
		nbuf := buf
		go func(nbuf string, val int64, z int64) {
			var y int64 = 0
			//log.Print("Creating big int for n=", val)

			// Count up
			for y = z; y <= val; y++ {
				nbuf = nbuf + strconv.FormatInt(y, 10)
			}
			//Count down
			//for y = y - 1; y > 0; y-- {
			//	buf = buf + strconv.FormatInt(y, 10)
			//}

			temp := big.NewInt(0)
			temp.SetString(nbuf, 10)

			//log.Print("Checking if n=", val, " is a probable prime.")
			if temp.ProbablyPrime(0) {
				log.Println("POSSIBLE PRIME, VERIFYING: n=", val)
				isPrime(val, temp)
			} else {
				//log.Println("not a probable prime: n=", val)
				//log.Println(nbuf)
			}
			swg.Done()
		}(nbuf, x, z)
	}
	swg.Wait()
}

func isPrime(x int64, num *big.Int) bool {
	i := big.NewInt(0)
	iSq := big.NewInt(0)
	iSq = iSq.Sqrt(num)

	for i.SetString("2", 10); i.Cmp(iSq) == -1; i.Add(i, big.NewInt(1)) {
		if num.Mod(num, i) == big.NewInt(0) {
			log.Println("*** NOT PRIME *** n=", x, " divisible by", i, ", ")
			return false
		}
	}
	log.Println("\n***** n=", x, " is prime! *****")
	return true
}
