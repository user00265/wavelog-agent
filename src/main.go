package main

import (
	"fmt"
	"sync"
)

func main() {
	// Recreated Main function in main.go and run both other main functions.
	// If the go routines hold up this should work till the cmd is exited
	var wg sync.WaitGroup
	wg.Add(2)
	go rigmain(&wg)
	go wsjtxmain(&wg)

	fmt.Println("Waiting for Routines to finish...")
	wg.Wait()
	fmt.Println("Done")
}
