package main

import (
	"github.com/dgryski/go-linebreak"
	"github.com/gdamore/tcell"
	"github.com/lithdew/blanc/layout"
	"github.com/mattn/go-runewidth"
)

type StyleFunc func(int) tcell.Style

type Text struct {
	style StyleFunc
	wrap  bool

	text  string
	width int
}

func NewText(text string) Text { return Text{text: text} }

func (t *Text) SetStyle(style tcell.Style)   { t.style = func(int) tcell.Style { return style } }
func (t *Text) SetStyleFunc(style StyleFunc) { t.style = style }

func (t *Text) SetText(text string) { t.text = text; t.width = runewidth.StringWidth(text) }
func (t *Text) SetWrap(wrap bool)   { t.wrap = wrap }

func (t Text) Width() int   { return t.width }
func (t Text) Text() string { return t.text }

func (t Text) Draw(s tcell.Screen, r layout.Rect) {
	x := r.X
	y := r.Y

	var text []rune
	if t.wrap {
		text = []rune(linebreak.Wrap(t.text, r.W, r.W))
	} else {
		text = []rune(t.text)
	}

	style := tcell.StyleDefault
	if t.style != nil {
		style = t.style(0)
	}

	clear(s, style, r.Left(), r.Top(), r.Right(), r.Bottom())

	for i := 0; i < len(text) && x <= r.X+r.W && y <= r.Y+r.H; i++ {
		if text[i] == '\n' {
			y += 1
			x = r.X
		} else {
			style := tcell.StyleDefault
			if t.style != nil {
				style = t.style(i)
			}
			s.SetContent(x, y, text[i], nil, style)
			x += runewidth.RuneWidth(text[i])
		}
	}
}
