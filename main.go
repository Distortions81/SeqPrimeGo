package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/dustin/go-humanize"
)

const startNPrime = 1000000

func main() {

	//Init
	var x int64 = startNPrime
	var buf bytes.Buffer

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

	//Create inital big.int string
	log.Println("Creating first big int buffer for n=", x)
	var z int64 = 0
	for z = 1; z < x; z++ {
		buf.WriteString(strconv.FormatInt(z, 10))
	}
	log.Println("Buffer size:", humanize.Bytes(uint64(buf.Len())))

	//Start checking:
	log.Println("Checking for n=x primes: ")
	for x = z; x < 9223372036854775807; x++ {

		buf.WriteString(strconv.FormatInt(x, 10))

		fmt.Print("Making big.int for n=", x, ", ")
		temp := big.NewInt(0)
		temp.SetString(buf.String(), 10)

		fmt.Print("Checking n=", x, ", ")
		if temp.ProbablyPrime(0) {
			log.Println("POSSIBLE PRIME: n=", x)
			if temp.ProbablyPrime(20) {
				log.Println("PROBABLE PRIME, VERIFYING: n=", x)
				isPrime(x, temp)
			}
		} else {
			//Print failure, do not log
			fmt.Print("!n=", x, ", ")
		}
	}
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
