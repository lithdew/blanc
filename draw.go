package main

import (
	"github.com/gdamore/tcell"
	"github.com/lithdew/blanc/layout"
)

var TcellFrame = [...]rune{
	tcell.RuneULCorner,
	tcell.RuneURCorner,
	tcell.RuneLLCorner,
	tcell.RuneLRCorner,
	tcell.RuneHLine,
	tcell.RuneHLine,
	tcell.RuneVLine,
	tcell.RuneVLine,
}

var (
	AsciiFrame       = [...]rune{'-', '-', '-', '-', '-', '-', '|', '|'}
	UnicodeFrame     = [...]rune{'┏', '┓', '┗', '┛', '━', '━', '┃', '┃'}
	UnicodeAltFrame  = [...]rune{'▛', '▜', '▙', '▟', '▀', '▄', '▌', '▐'}
	UnicodeAlt2Frame = [...]rune{'╔', '╗', '╚', '╝', '═', '═', '║', '║'}
	SpaceFrame       = [...]rune{' ', ' ', ' ', ' ', ' ', ' ', ' ', ' '}
)

func Clear(s tcell.Screen, style tcell.Style, r layout.Rect) {
	clear(s, style, r.Left(), r.Top(), r.Right(), r.Bottom())
}

func clear(s tcell.Screen, style tcell.Style, x1, y1, x2, y2 int) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}
}

func box(s tcell.Screen, style tcell.Style, x1, y1, x2, y2 int, frames [8]rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}
	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, frames[4], nil, style)
		s.SetContent(col, y2, frames[5], nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, frames[6], nil, style)
		s.SetContent(x2, row, frames[7], nil, style)
	}
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, frames[0], nil, style)
		s.SetContent(x2, y1, frames[1], nil, style)
		s.SetContent(x1, y2, frames[2], nil, style)
		s.SetContent(x2, y2, frames[3], nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}
}
