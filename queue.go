package main

import (
	"container/heap"
	"github.com/gdamore/tcell"
)

type EventHandlerFunc func(ev tcell.Event) bool

func (fn EventHandlerFunc) HandleEvent(ev tcell.Event) bool {
	return fn(ev)
}

var _ heap.Interface = (*InputQueue)(nil)

type InputQueue []tcell.EventHandler

func (q InputQueue) Len() int            { return len(q) }
func (q InputQueue) Less(i, j int) bool  { panic("implement me") }
func (q InputQueue) Swap(i, j int)       { panic("implement me") }
func (q *InputQueue) Push(x interface{}) { panic("implement me") }
func (q *InputQueue) Pop() interface{}   { panic("implement me") }
