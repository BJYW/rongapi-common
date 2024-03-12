package mgopool

import (
	"fmt"
	"testing"
)

func TestPoolSize(t *testing.T) {
	cha := make(chan int, 20)
	fmt.Println(len(cha))
	cha <- 1
	fmt.Println(len(cha))
	fmt.Println(cap(cha))
}
