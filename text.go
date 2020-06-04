package blanc

import (
	"github.com/dgryski/go-linebreak"
	"github.com/gdamore/tcell"
	"github.com/lithdew/blanc/layout"
	"github.com/mattn/go-runewidth"
)

type StyleFunc func(int) tcell.Style

type Text struct {
	dirty bool // does this view need to be recomputed?
	lastW int  // last rendered width

	style StyleFunc
	wrap  bool

	text  string
	buf   []rune
	width int
}

func NewText(text string) Text { return Text{text: text, buf: []rune(text), lastW: -1} }

func (t *Text) SetText(text string) {
	t.dirty = true
	t.text = text
	t.width = runewidth.StringWidth(text)
}

func (t *Text) SetWrap(wrap bool) {
	t.dirty = !t.wrap && wrap
	t.wrap = wrap
}

func (t *Text) SetStyle(style tcell.Style)   { t.style = func(int) tcell.Style { return style } }
func (t *Text) SetStyleFunc(style StyleFunc) { t.style = style }

func (t Text) Width() int   { return t.width }
func (t Text) Text() string { return t.text }

func (t *Text) Draw(s tcell.Screen, r layout.Rect) {
	if r.W != t.lastW {
		t.dirty = true
		t.lastW = r.W
	}

	if t.dirty {
		if t.wrap {
			t.buf = []rune(linebreak.Wrap(t.text, r.W, r.W))
		} else {
			t.buf = []rune(t.text)
		}
		t.dirty = false
	}

	x := r.X
	y := r.Y

	for i := 0; i < len(t.buf) && x <= r.X+r.W && y <= r.Y+r.H; i++ {
		if t.buf[i] == '\n' {
			y += 1
			x = r.X
		} else {
			style := tcell.StyleDefault
			if t.style != nil {
				style = t.style(i)
			}
			s.SetContent(x, y, t.buf[i], nil, style)
			x += runewidth.RuneWidth(t.buf[i])
		}
	}
}
