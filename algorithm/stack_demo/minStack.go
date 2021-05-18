package main

import (
	"math"
	"fmt"
)

type MinStack struct {
	stack    []int
	minStack []int
}

func Construct() MinStack{
	return MinStack{
		stack: []int{},
		minStack: []int{math.MaxInt64},
	}
}


func (this *MinStack) Push(val int){
	this.stack = append(this.stack, val)
	this.minStack = append(this.minStack, min(this.minStack[len(this.minStack) -1 ], val))
}

func (this *MinStack) Pop(){
	this.stack = this.stack[:len(this.stack) - 1]
	this.minStack = this.minStack[:len(this.minStack) - 1]
}

func (this *MinStack) Top() int {
	return this.stack[len(this.stack) - 1]
}

func (this *MinStack) GetMin() int {
	return this.minStack[len(this.minStack) - 1]
}


func min(x, y int) int {
	if x < y {
		return x
	}

	return y
}

func main() {
	stack := Construct()
	stack.Push(3)
	stack.Push(2)
	stack.Push(5)
	min := stack.GetMin()
	fmt.Println(min)
}