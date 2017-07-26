package main

import (
	"fmt"
)

type Integer int

func (a Integer) Less(b Integer) bool {
	return a < b
}

func (a *Integer) Add(b Integer) {
	*a += b
}

type LessAdder interface {
	Less(b Integer) bool
	Add(b Integer)
}

type Lesser interface {
	Less(b Integer) bool
}


func main()  {
	var a Integer = 1
	var b LessAdder = &a
	//var c LessAdder = a

	var b1 Lesser = &a
	var b2 Lesser = a

	fmt.Println(b.Less(2))
	//fmt.Println(c.Less(2))
	fmt.Println(b1.Less(2))
	fmt.Println(b2.Less(2))

}

