package main

import (
	"container/heap"
	"github.com/gdamore/tcell"
)

type InputListener interface {
	tcell.EventHandler
}

var _ heap.Interface = (*InputQueue)(nil)

type InputQueue []InputListener

func (q InputQueue) Len() int            { return len(q) }
func (q InputQueue) Less(i, j int) bool  { panic("implement me") }
func (q InputQueue) Swap(i, j int)       { panic("implement me") }
func (q *InputQueue) Push(x interface{}) { panic("implement me") }
func (q *InputQueue) Pop() interface{}   { panic("implement me") }
