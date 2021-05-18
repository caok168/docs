package main

import (
	"math"
	"fmt"
)

type MinStack struct {
	stack []int
	minStack []int
}

func Constructor() MinStack {
	return MinStack{
		stack:[]int{},
		minStack:[]int{math.MaxInt64},
	}
}

func (this *MinStack) Push(val int){
	this.stack = append(this.stack, val)
	this.minStack = append(this.minStack, min(this.minStack[len(this.minStack) - 1], val))
}

func (this *MinStack) Pop() {
	this.stack = this.stack[:len(this.stack)-1]
	this.minStack = this.minStack[:len(this.minStack)-1]
}

func (this *MinStack) Top() int {
	return this.stack[len(this.stack)-1]
}

func (this *MinStack) GetMin() int{
	return this.minStack[len(this.minStack)-1]
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func main() {
	obj := Constructor()
	obj.Push(3)
	obj.Push(4)
	obj.Push(5)

	min := obj.GetMin()
	fmt.Println("min:", min)

	obj.Push(1)
	min = obj.GetMin()
	fmt.Println("min:", min)
}