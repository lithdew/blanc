package main

import (
	"github.com/gdamore/tcell"
	"github.com/lithdew/blanc/layout"
)

type Textbox struct {
	label []rune
	text  []rune

	ptr int // cursor start index
	pos int // cursor end index
}

func newTextbox() *Textbox {
	return &Textbox{pos: -1}
}

func (t *Textbox) cursorPos() (start int, end int) {
	start = t.ptr // selection start index
	if t.pos != -1 && start > t.pos {
		start = t.pos
	}
	end = t.ptr // selection end index
	if t.pos != -1 && end < t.pos {
		end = t.pos
	}
	return start, end
}

func (t *Textbox) render(s tcell.Screen, style tcell.Style, r layout.Rect) {
	if r.H < 1 {
		return
	}

	cursor := tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite)

	j := 0
	for i := 0; i < len(t.label) && r.X+j < r.X+r.W-1; i, j = i+1, j+1 {
		s.SetContent(r.X+j, r.Y, t.label[i], nil, style)
	}

	start, end := t.cursorPos()

	for i := 0; i < len(t.text) && r.X+j < r.X+r.W-1; i, j = i+1, j+1 {
		if i >= start && i <= end {
			s.SetContent(r.X+j, r.Y, t.text[i], nil, cursor)
		} else {
			s.SetContent(r.X+j, r.Y, t.text[i], nil, style)
		}
	}

	if t.ptr == len(t.text) && t.pos == -1 { // render cursor
		s.SetContent(r.X+j, r.Y, tcell.RuneBlock, nil, cursor.Reverse(true))
	}
}

func (t *Textbox) selectLeft() {
	if t.ptr == 0 {
		return
	}
	if t.pos == -1 {
		t.pos = t.ptr - 1
	}
	t.ptr--
}

func (t *Textbox) selectRight() {
	if t.ptr == len(t.text) {
		return
	}
	if t.pos == -1 {
		t.pos = t.ptr
	} else {
		t.ptr++
	}
}

func (t *Textbox) moveLeft() {
	if t.pos != -1 {
		if t.ptr > t.pos {
			t.ptr = t.pos
		}
		t.pos = -1
		return
	}
	if t.ptr == 0 {
		return
	}
	t.ptr--
}

func (t *Textbox) moveRight() {
	if t.pos != -1 {
		if t.ptr < t.pos {
			t.ptr = t.pos
		}
		t.pos = -1
	}
	if t.ptr == len(t.text) {
		return
	}
	t.ptr++
}

func (t *Textbox) selectPrevWord() {
	if t.pos == -1 {
		t.pos = t.ptr - 1
	}

	for i := t.ptr - 2; i >= 0; i-- {
		if t.text[i] == ' ' {
			t.ptr = i + 1
			return
		}
	}
	t.ptr = 0
}

func (t *Textbox) movePrevWord() {
	t.pos = -1
	for i := t.ptr - 2; i >= 0; i-- {
		if t.text[i] == ' ' {
			t.ptr = i + 1
			return
		}
	}
	t.ptr = 0
}

func (t *Textbox) selectNextWord() {
	if t.pos == -1 {
		t.pos = t.ptr
	}
	for i := t.ptr + 2; i < len(t.text); i++ {
		if t.text[i] == ' ' {
			t.ptr = i - 1
			return
		}
	}
	t.ptr = len(t.text)
}

func (t *Textbox) moveNextWord() {
	t.pos = -1
	for i := t.ptr; i < len(t.text); i++ {
		if t.text[i] == ' ' {
			t.ptr = i + 1
			return
		}
	}
	t.ptr = len(t.text)
}

func (t *Textbox) selectAll() {
	t.ptr = len(t.text)
	t.pos = 0
}

func (t *Textbox) moveToEnd() {
	t.pos = -1
	for i := t.ptr; i < len(t.text); i++ {
		if t.text[i] == '\n' {
			t.ptr = i
			return
		}
	}
	t.ptr = len(t.text)
}

func (t *Textbox) push(r rune) {
	if t.pos != -1 {
		t.pop()
	}
	t.text = append(t.text[:t.ptr], append([]rune{r}, t.text[t.ptr:]...)...)
	t.moveRight()
}

func (t *Textbox) pop() {
	if t.ptr == 0 && t.pos == -1 {
		return
	}

	if t.pos == -1 { // normal backspace.
		t.text = append(t.text[:t.ptr-1], t.text[t.ptr:]...)
		t.moveLeft()
		return
	}

	start, end := t.cursorPos()
	if end == len(t.text) {
		end--
	}

	t.text = append(t.text[:start], t.text[end+1:]...)
	t.ptr = start
	t.pos = -1
}

func (t *Textbox) setLabel(label string) {
	t.label = []rune(label)
}

func (t *Textbox) setText(text string) {
	t.text = []rune(text)
}
