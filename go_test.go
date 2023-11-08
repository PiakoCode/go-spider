package main

import (
	"fmt"
	"testing"
)

func foo(c chan bool, cnt int) {
	fmt.Println("hello,", cnt)
	c <- true
}

func TestGo(t *testing.T) {
	c1 := make(chan bool)
	var cnt int
	for {
		go foo(c1, cnt)
		if cnt == 500 {
			break
		}
	}
	for {
		select {
		case <-c1:
			cnt += 1
		}
	}
	close(c1)

}
