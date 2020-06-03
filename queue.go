package main

import (
	"github.com/gdamore/tcell"
)

type EventHandlerFunc func(ev tcell.Event) bool

func (fn EventHandlerFunc) HandleEvent(ev tcell.Event) bool {
	return fn(ev)
}
