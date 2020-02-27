package main

import (
	"fmt"
	"sync"
)

// Communicate by Sharing memory

func addByShareMemory(n int) []int {
	var ints []int
	var wg sync.WaitGroup
	var mux sync.Mutex

	wg.Add(n) //add n counter
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done() // counter--
			mux.Lock()
			ints = append(ints, i)
			mux.Unlock()
		}(i)
	}

	wg.Wait()
	return ints
}
// share memory by communicate
func addByShareCommunicate(n int) []int {
	var ints []int
	channel := make(chan int, n)// build a channel

	for i := 0; i < n; i++ {
		go func(channel chan<- int, order int) { // write into channel value
			channel <- order
		}(channel, i)
	}

	for i := range channel {
		ints = append(ints, i)

		if len(ints) == n {
			break
		}
	}
	close(channel)

	return ints
}
func main() {
	// foo := addByShareMemory(10)
	foo := addByShareCommunicate(10)
	fmt.Println(len(foo))
	fmt.Println(foo)
}
