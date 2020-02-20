package main

import (
	"fmt"
	"runtime"
	"sync"
)

var wg2 sync.WaitGroup

func print_characters() {
	defer wg2.Done()
	for ch := 'a'; ch < 'a' + 26; ch++ {
		fmt.Printf("%c ", ch)
	}
}

func print_digits() {
	defer wg2.Done()
	for number := 1; number < 27; number++ {
		fmt.Printf("%d ", number)
	}
}

func concurrent() {
	wg2.Add(2)
	fmt.Println("Starting go routines")
	go print_characters()
	go print_digits()
	fmt.Println("Waiting to Finish")
	wg2.Wait()
	fmt.Println("\nTerminating processing")
}

func main() {
	runtime.GOMAXPROCS(1)
	concurrent()
}
