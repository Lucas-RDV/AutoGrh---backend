package main

import (
	"fmt"
	"time"
)

func main() {
	s := "gopher"
	fmt.Printf("Hello and welcome, %s!", s)
	time.Sleep(10 * time.Second)
}
