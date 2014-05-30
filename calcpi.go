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

func calcMember(k int64) *big.Rat {
	member := big.NewRat(1103+26390*k, 1)
	for i := int64(1); i <= k; i++ {
		member = member.Mul(member, big.NewRat(i*(i+k)*(i+2*k)*(i+3*k), 1))

		member = member.Mul(member, big.NewRat(1, (396*396*396*396)))

		member = member.Mul(member, big.NewRat(1, i*i*i*i))
	}
	return member
}

func calcPi(from, to int64) *big.Rat {
	pi := big.NewRat(0, 1)
	for i := from; i < to; i++ {
		member := calcMember(i)
		pi = pi.Add(pi, member)
	}
	return pi
}

func main() {
	startTime := time.Now()
	sum := big.NewRat(0, 1)
	ch := make(chan *big.Rat)

	parseArgs()

	log("Starting execution with arguments:\n")
	log(fmt.Sprintf("	-p           =  %v - precision as the number of members in the sequence\n", precision))
	log(fmt.Sprintf("	-t, --tasks  =  %v - number of async tasks / go-routines\n", tasks))
	log(fmt.Sprintf("	-q,          =  %v - whether to log more verbose output\n", quiteMode))
	if outputFile != "<not given>" {
		log(fmt.Sprintf("	-o           =  %v - log file\n", outputFile))
	}

	var membersPerTask int = int(math.Ceil(float64(precision) / float64(tasks)))
	for i := 0; i < tasks; i++ {
		var from, to int64 = int64(i * membersPerTask), int64((i + 1) * membersPerTask)

		log(fmt.Sprintf("Calculating the sum of sequence members %v to %v...\n", from, to))

		go func(from, to int64) {
			pi := calcPi(from, to)
			ch <- pi
		}(from, to)
	}

	for i := 0; i < tasks; i++ {
		value := <-ch

		sum.Add(sum, value)
		log(fmt.Sprintf("... temporary sum of the sequence %v\n", sum.FloatString(20)))
	}

	log(fmt.Sprintf("Total sum of the sequence %v\n", sum.FloatString(20)))

	coef := big.NewRat(1, 1)
	coef.SetFloat64(math.Sqrt(8) / (99 * 99))
	inversed_ans := coef.Mul(coef, sum)

	pi := inversed_ans.Inv(inversed_ans)

	log(fmt.Sprintf("Pi: %v\n", pi.FloatString(20)))

	log(fmt.Sprintf("Total time: %v\n", time.Since(startTime)))
	fmt.Printf("Total time: %v\n", time.Since(startTime))

	fmt.Println("***FIX INFORMATIONAL LOGGING")
}
