package main

import (
	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
)

func box(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, r rune) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	for col := x1; col <= x2; col++ {
		s.SetContent(col, y1, tcell.RuneHLine, nil, style)
		s.SetContent(col, y2, tcell.RuneHLine, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.SetContent(x1, row, tcell.RuneVLine, nil, style)
		s.SetContent(x2, row, tcell.RuneVLine, nil, style)
	}
	if y1 != y2 && x1 != x2 {
		s.SetContent(x1, y1, tcell.RuneULCorner, nil, style)
		s.SetContent(x2, y1, tcell.RuneURCorner, nil, style)
		s.SetContent(x1, y2, tcell.RuneLLCorner, nil, style)
		s.SetContent(x2, y2, tcell.RuneLRCorner, nil, style)
	}
	for row := y1 + 1; row < y2; row++ {
		for col := x1 + 1; col < x2; col++ {
			s.SetContent(col, row, r, nil, style)
		}
	}
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	var deferred []rune

	i := 0
	zwj := false
	width := 0

	for _, r := range str {
		if r == '\u200d' {
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				width = 1
			}
			deferred = append(deferred, r)
			zwj = true
			continue
		}
		if zwj {
			deferred = append(deferred, r)
			zwj = false
			continue
		}
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				width = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += width
			}
			deferred = nil
			width = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += width
			}
			deferred = nil
			width = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += width
	}
}
