package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func (dbs *dbScanner) failIfStrictf(format string, args ...interface{}) {
	if dbs.nonStrict {
		format = "WARNING: " + format
	} else {
		format = "ERROR: " + format
	}
	_, _ = fmt.Fprintf(os.Stderr, "\n"+format+"\n\n", args...)
	if !dbs.nonStrict {
		os.Exit(1)
	}
}

func fatalf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func printLine(s string) {
	fmt.Println(strings.Repeat(s, 78))
}

type progressBarHolder struct {
	disable bool
	pb      *pb.ProgressBar
}

func (pbh progressBarHolder) Finish() {
	if pbh.disable {
		return
	}
	pbh.pb.Finish()
}

func (pbh progressBarHolder) Increment() {
	// fmt.Printf("%#v\n", pbh)
	if pbh.disable {
		return
	}
	pbh.pb.Increment()
}

func createProgressBar(count int, disable bool) progressBarHolder {
	pbh := progressBarHolder{
		disable: disable,
	}
	if !disable {
		pbh.pb = pb.StartNew(count)
	}
	return pbh
}

func formatInt(n int, sep string) string {
	var numParts []int
	// left := n
	for {
		order := n % 1000
		n /= 1000
		numParts = append(numParts, order)
		if n == 0 {
			break
		}
	}
	reverseInts(numParts)
	s := make([]string, len(numParts))
	numFmt := "%3s"
	for j, order := range numParts {
		s[j] = fmt.Sprintf(numFmt, strconv.Itoa(order))
		if j == 0 {
			numFmt = "%03s"
		}
	}
	return strings.Join(s, sep)
}

func reverseInts(a []int) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

type pairFormatter struct {
	width int
}

func (p pairFormatter) Pairs(pairs ...string) string {
	format := "%" + fmt.Sprintf("%ds", p.width) + ": %s"
	lines := make([]string, 0, len(pairs)/2)
	for i := range pairs {
		if i%2 == 1 {
			continue
		}
		lines = append(lines, fmt.Sprintf(format, pairs[i], pairs[i+1]))
	}
	return strings.Join(lines, "\n")
}
