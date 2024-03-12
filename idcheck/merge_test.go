package idcheck

import (
	"fmt"
	"testing"
)

func TestMerge(t *testing.T) {
	fmt.Println(merge("one", "two"))
	fmt.Println(merge(merge("one", "two"), "three"))
	fmt.Println(merge("zero", merge("one", "two")))
	fmt.Println(merge(merge("one", "two"), merge("three", "four")))
}

func TestMergeZero(t *testing.T) {
	fmt.Println(merge("", ""))
	fmt.Println(merge("one", ""))
	fmt.Println(merge("", "one"))
}

func TestMergeComplicatedZero(t *testing.T) {
	fmt.Println(merge(merge("one", ""), ""))
	fmt.Println(merge(merge("one", "two"), ""))
	fmt.Println(merge(merge("one", ""), "two"))
	fmt.Println(merge(merge("", "one"), "two"))
	fmt.Println(merge("", merge("one", "two")))
	fmt.Println(merge("zero", merge("", "one")))
	fmt.Println(merge("zero", merge("one", "")))
	fmt.Println(merge("", merge("one", "")))
	fmt.Println(merge("", merge("", "one")))
	fmt.Println(merge("one", merge("", "")))
	fmt.Println(merge(merge("one", ""), merge("three", "four")))
	fmt.Println(merge(merge("one", "two"), merge("", "four")))
	fmt.Println(merge(merge("one", "two"), merge("three", "")))
	fmt.Println(merge(merge("", "two"), merge("three", "four")))

	fmt.Println(merge(merge("", ""), merge("three", "four")))
	fmt.Println(merge(merge("", "two"), merge("three", "")))
	fmt.Println(merge(merge("one", ""), merge("three", "")))
	fmt.Println(merge(merge("", "two"), merge("", "four")))
	fmt.Println(merge(merge("one", ""), merge("", "four")))
	fmt.Println(merge(merge("one", "two"), merge("", "")))

	fmt.Println(merge(merge("", ""), merge("", "four")))
	fmt.Println(merge(merge("", "two"), merge("", "")))
	fmt.Println(merge(merge("", ""), merge("three", "")))
	fmt.Println(merge(merge("one", ""), merge("", "")))
}
