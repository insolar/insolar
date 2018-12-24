package main

import (
	"fmt"
	"time"
)

func main() {
	t := make(chan interface{})

	go func() {
		<-t
		fmt.Println("1")
	}()

	go func() {
		<-t
		fmt.Println("2")
	}()

	go func() {
		<-t
		fmt.Println("3")
	}()

	go func() {
		<-t
		fmt.Println("4")
	}()

	go func() {
		<-t
		fmt.Println("5")
	}()

	time.Sleep(time.Second * 3)
	close(t)

	time.Sleep(time.Second * 5)
}
