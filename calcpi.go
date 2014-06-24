package main

import (
	"fmt"
	"math"
	"math/big"
	"os"
	"strconv"
	"syscall"
	"time"
)

var precision, tasks int = 100, 1
var quiteMode bool = false
var outputFile string = "<not given>"
var outputStream *os.File
var factorial []*big.Rat
var powerOf396 []*big.Rat

func log(msg string) {
	if quiteMode == false {
		outputStream.WriteString(msg)
	}
}

func parseArgs() {
	outputStream = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg == "-p" {
			i++
			precision, _ = strconv.Atoi(os.Args[i])
		} else if arg == "-t" || arg == "--tasks" {
			i++
			tasks, _ = strconv.Atoi(os.Args[i])
		} else if arg == "-q" {
			quiteMode = true
		} else if arg == "-o" {
			i++
			outputFile = os.Args[i]
			outputStream, _ = os.Create(outputFile)
		}
	}

}

func populateHelperSlices() {
	factorial[0] = big.NewRat(1, 1)
	factorial[1] = big.NewRat(1, 1)

	powerOf396[0] = big.NewRat(1, 1)
	powerOf396[1] = big.NewRat(1, 396)

	for i := int64(2); i < int64(5*precision); i++ {
		n := big.NewRat(i, 1)
		factorial[i] = n.Mul(n, factorial[i-1])

		denom := big.NewRat(1, 396)
		powerOf396[i] = denom.Mul(denom, powerOf396[i-1])

	}
}

func calcMember(k int64) *big.Rat {
	member := big.NewRat(1103+26390*k, 1)

	member.Mul(member, factorial[k*4])
	member.Mul(member, powerOf396[k*4])

	kFact := factorial[k]
	kFact.Mul(kFact, kFact)
	kFact.Mul(kFact, kFact)
	kFact.Inv(kFact)

	member.Mul(member, kFact)

	return member
}

func calcPi(from, to int64, ch chan *big.Rat) {
	startTime := time.Now()
	log(fmt.Sprintf("...... calcPi STARTED from %v to %v\n", from, to))
	pi := big.NewRat(0, 1)
	for i := from; i < to; i++ {
		member := calcMember(i)
		pi.Add(pi, member)
	}
	log(fmt.Sprintf("...... calcPi FINISHED from %v to %v for %v\n", from, to, time.Since(startTime)))
	ch <- pi
	// return pi
}

func main() {
	startTime := time.Now()

	parseArgs()

	log("Starting execution with arguments:\n")
	log(fmt.Sprintf("	-p           =  %v - precision as the number of members in the sequence\n", precision))
	log(fmt.Sprintf("	-t, --tasks  =  %v - number of async tasks / go-routines\n", tasks))
	log(fmt.Sprintf("	-q,          =  %v - whether to log more verbose output\n", quiteMode))
	if outputFile != "<not given>" {
		log(fmt.Sprintf("	-o           =  %v - log file\n", outputFile))
	}

	sum := big.NewRat(0, 1)
	ch := make(chan *big.Rat, tasks+1)

	factorial = make([]*big.Rat, 5*precision, 5*precision)
	powerOf396 = make([]*big.Rat, 5*precision, 5*precision)
	populateHelperSlices()

	var membersPerTask int = int(math.Ceil(float64(precision) / float64(tasks)))
	for i := 0; i < tasks; i++ {
		var from, to int64 = int64(i * membersPerTask), int64((i + 1) * membersPerTask)

		log(fmt.Sprintf("Calculating the sum of sequence members %v to %v...\n", from, to))

		// go func(from, to int64) {
		// 	pi := calcPi(from, to)
		// 	log("push to channel\n")
		// 	ch <- pi
		// 	log("pushed to channel\n")
		// }(from, to)
		go calcPi(from, to, ch)
	}
	for i := 0; i < tasks; i++ {

		log(fmt.Sprintf("... waiting for task %v of %v to finish\n", i, tasks))
		value := <-ch

		sum.Add(sum, value)
		log(fmt.Sprintf("... task %v of %v finished\n", i, tasks))
		// log(fmt.Sprintf("... temporary sum of the sequence %v\n", sum.FloatString(30)))
	}

	// log(fmt.Sprintf("Total sum of the sequence %v\n", sum.FloatString(30)))

	coef := big.NewRat(1, 1)
	coef.SetFloat64(math.Sqrt(8) / (99 * 99))
	inversed_ans := coef.Mul(coef, sum)

	pi := inversed_ans.Inv(inversed_ans)

	log(fmt.Sprintf("Pi: %v\n", pi.FloatString(30)))

	log(fmt.Sprintf("Total time: %v\n", time.Since(startTime)))
}
