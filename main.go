package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/remeh/sizedwaitgroup"
)

const startNPrime = 1000000
const debug = true
const logName = "nPrimes.log"
const reportSeconds = 60

var lastReport time.Time

func main() {

	//Vars
	var x int64 = startNPrime - 1
	var z int64 = 0
	var buf bytes.Buffer
	lastReport = time.Now()

	//Logging setup
	lf, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer lf.Close()
	mw := io.MultiWriter(os.Stdout, lf) //log and stdout
	log.SetOutput(mw)

	//Preperation
	log.Println("Creating string for n=", x+1)
	for z = 1; z < x; z++ {
		buf.WriteString(strconv.FormatInt(z, 10))
	}
	buf.WriteString(strconv.FormatInt(x, 10))
	log.Println("String size:", humanize.Bytes(uint64(buf.Len())))
	log.Println("Making big.int for n=", x+1)
	var bigPrime big.Int
	bigPrime.SetString(buf.String(), 10)
	log.Println("Checking for n=x primes: ")

	//Wait group with cpu threds
	threads := runtime.NumCPU()
	swg := sizedwaitgroup.New(threads)
	log.Println("Starting", threads, "new threads.")

	for x = startNPrime; x < 9223372036854775807; x++ {

		shiftDigits(&bigPrime, x)

		//Add to wait group
		swg.Add()
		go func(lx int64, nbp big.Int) {
			defer swg.Done()

			isdebug(fmt.Sprintf("n=%v, ", lx))
			if nbp.ProbablyPrime(0) {
				log.Println("POSSIBLE PRIME: n=", lx)
				if nbp.ProbablyPrime(20) {
					log.Println("PROBABLE PRIME, VERIFYING: n=", lx)
					isPrime(lx, &nbp)
				}
			}
		}(x, bigPrime)
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

func isdebug(str string) {
	if debug && time.Since(lastReport) > reportSeconds*time.Second {
		fmt.Print(str)
		lastReport = time.Now()
	}
}

func shiftDigits(bigPrime *big.Int, x int64) {
	//Shift over digits, this is faster than re-serializing the big.int
	//Calculate how many digits we need to move over
	toAdd := int64(math.Pow(10, float64(len(strconv.FormatInt(x, 10)))))
	//Mutiply to move required number of digits, for our new number
	bigPrime.Mul(bigPrime, big.NewInt(toAdd))
	//Add our new digits
	bigPrime.Add(bigPrime, big.NewInt(x))
}
