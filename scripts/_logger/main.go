//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func myUsage() {
	fmt.Printf("Usage: %s [flags] outfile\n", os.Args[0])
	flag.PrintDefaults()
}

var debug bool

func main() {
	var noBuffer bool

	flag.Usage = myUsage
	flag.BoolVar(&noBuffer, "no-buffer", false, "disables buffer")
	flag.BoolVar(&debug, "debug", false, "print debug info")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	flush := func() {}

	outFile := flag.Arg(0)
	fmt.Println("file:", flag.Arg(0))

	var w io.Writer
	f, err := os.Create(outFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	w = f

	var bufWriter *bufio.Writer
	if !noBuffer {
		bufWriter = bufio.NewWriter(f)
		w = bufWriter
		flush = func() {
			err := bufWriter.Flush()
			if err != nil {
				log.Println("flush failed:", err)
			}
			if debug {
				log.Println("flush")
			}
		}
	}

	// read stdin
	reopen := make(chan struct{})
	outgoing := make(chan []byte)
	go func() {
		inputReader := bufio.NewReader(os.Stdin)
		for {
			b, err := inputReader.ReadBytes('\n')
			if err != nil {
				log.Fatal("read error: %v", err)
				return
			}
			outgoing <- b
		}
	}()

	stop := make(chan bool)
	tick := time.Tick(time.Second)
	go func() {
		for {
			select {
			case <-tick:
				flush()
			case b := <-outgoing:
				n, err := w.Write(b)
				if debug {
					log.Println("write ", n, "bytes")
				}
				if err != nil {
					log.Fatal("write failed:", err)
					return
				}

			case <-reopen:
				flush()
				f.Close()

				fmt.Println("reopen", outFile)
				f, err = os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Fatal("reopen failed:", err)
					return
				}
				w = f
				if !noBuffer {
					bufWriter = bufio.NewWriter(f)
					w = bufWriter
				}
			case <-stop:
				flush()
				return
			}
		}
	}()

	finish := make(chan bool)
	var sigUSR2 = make(chan os.Signal, 1)
	var sigINT = make(chan os.Signal, 1)
	signal.Notify(sigUSR2, syscall.SIGUSR2)
	signal.Notify(sigINT, syscall.SIGINT) // inside the goroutine
	go func() {
		for {
			select {
			case <-sigUSR2:
				reopen <- struct{}{}
			case <-sigINT:
				// stop work
				stop <- true
				close(finish)
			}
		}
	}()

	<-finish
}
