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
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/remeh/sizedwaitgroup"
)

// Constants/vars
var startNPrime int64 = 1000000

const logName = "nPrimes.log"

/* Progress reports */
const progressFile = "progress.dat"
const progressInterval = 30 * time.Second

var lastProgress time.Time
var progressLock sync.Mutex

func main() {

	//Vars
	var z int64 = 0
	var buf bytes.Buffer
	var number int64 = 0
	lastProgress = time.Now()

	prog, err := os.ReadFile(progressFile)
	if err != nil {
		log.Println("No progress file found, starting from scratch")
	} else {
		number, err = strconv.ParseInt(string(prog), 10, 64)
		if err != nil {
			log.Println("Error reading progress file, starting from scratch")
		}
		startNPrime = number
	}

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
	var x int64 = startNPrime - 1
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

	//Wait group with cpu threads
	threads := runtime.NumCPU()

	swg := sizedwaitgroup.New(threads)
	pcg := sizedwaitgroup.New(threads * 2)
	log.Printf("Detected %v vCPUs.\n", threads)

	//We basically buffer up a ton of big.ints we can process when a open thread appears
	for x = startNPrime; x < 9223372036854775807; x++ {
		pcg.Add() //Precalculate next n, within limits

		shiftDigits(&bigPrime, x) //Modifying big.int is slow

		go func(lx int64, nbp big.Int) {
			swg.Add() //We are ready, but wait our turn
			progressReport(lx, "ch-prob:")

			if nbp.ProbablyPrime(0) {
				log.Printf("* POSSIBLE PRIME: n=%v *\n", lx)
				if nbp.ProbablyPrime(20) {
					log.Printf("** PROBABLE PRIME: n=%v **\n", lx)
					isPrime(lx, &nbp)
				}
			}
			//Done, let new threads run
			pcg.Done()
			swg.Done()
		}(x, bigPrime)

	}

	//Wait for everything to finish before exiting.
	pcg.Wait()
	swg.Wait()
}

func isPrime(x int64, num *big.Int) bool {
	i := big.NewInt(0)
	iSq := big.NewInt(0)
	iSq = iSq.Sqrt(num)

	for i.SetString("2", 10); i.Cmp(iSq) == -1; i.Add(i, big.NewInt(1)) {
		progressReport(x, "ch-isPrime:")
		if num.Mod(num, i) == big.NewInt(0) {
			log.Printf("*** NOT PRIME: n=%v is divisible by %v ***\n", x, i)
			return false
		}
	}
	log.Printf("*** VERIFIED PRIME: n=%v ***\n", x)
	return true
}

func shiftDigits(bigPrime *big.Int, x int64) {
	//Shift over digits, this is faster than re-serializing the big.int
	//Calculate how many digits we need to move over
	toAdd := int64(math.Pow(10, float64(len(strconv.FormatInt(x, 10)))))

	progressReport(x, "mult:")
	//Mutiply to move required number of digits, for our new number
	bigPrime.Mul(bigPrime, big.NewInt(toAdd))
	progressReport(x, "add:")
	//Add our new digits
	bigPrime.Add(bigPrime, big.NewInt(x))
}

func progressReport(x int64, message string) {
	progressLock.Lock()
	if time.Since(lastProgress) > progressInterval {
		log.Println(message+" n=", x)

		//log.Println("Saving progress")
		err := os.WriteFile(progressFile, []byte(strconv.FormatInt(x, 10)), 0644)
		if err != nil {
			log.Println("Error saving progress:", err)
		}
		lastProgress = time.Now()
	}
	progressLock.Unlock()
}
