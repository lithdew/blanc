package main

import (
	"github.com/gdamore/tcell"
	"github.com/lithdew/blanc/layout"
)

type Cell struct {
	rect   layout.Rect
	cursor layout.Rect
	buf    []string
}

func NewCell(rect layout.Rect) *Cell {
	return &Cell{rect: rect, cursor: rect}
}

func (c *Cell) Render(screen tcell.Screen, style tcell.Style) {
	clear(screen, style, c.rect.Left(), c.rect.Top(), c.rect.Right(), c.rect.Bottom())
	style = style.Reverse(true)

}
