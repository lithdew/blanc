package main

import (
	"github.com/gdamore/tcell"
	"github.com/lithdew/blanc/layout"
	"github.com/mattn/go-runewidth"
	"unicode"
)

type Textbox struct {
	label []rune
	text  []rune

	ptr int // cursor start index
	pos int // cursor end index
	dir int // -1 is left, 1 is right
}

func newTextbox() *Textbox {
	return &Textbox{pos: -1, dir: -1}
}

func (t *Textbox) selectedArea() (start int, end int) {
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
	for i := 0; i < len(t.label) && r.X+j < r.X+r.W-1; i, j = i+1, j+runewidth.RuneWidth(t.label[i]) {
		s.SetContent(r.X+j, r.Y, t.label[i], nil, style)
	}
	k := j

	start, end := t.selectedArea()

	for i := 0; i < len(t.text) && r.X+j < r.X+r.W-1; i, j = i+1, j+runewidth.RuneWidth(t.text[i]) {
		if t.pos != -1 && i >= start && i <= end {
			s.SetContent(r.X+j, r.Y, t.text[i], nil, cursor)
		} else {
			s.SetContent(r.X+j, r.Y, t.text[i], nil, style)
		}
	}

	if t.pos == -1 {
		s.ShowCursor(r.X+k+t.ptr, r.Y)
	} else {
		s.HideCursor()
	}
}

func (t *Textbox) selectLeft() {
	if t.pos != -1 && t.ptr == t.pos && t.dir == 1 {
		t.pos = -1
		return
	}
	if t.pos == -1 {
		t.pos = t.ptr - 1
		t.dir = -1
	}
	t.shiftLeft()
}

func (t *Textbox) selectRight() {
	if t.ptr == len(t.text) {
		return
	}

	if t.pos != -1 && t.ptr == t.pos && t.dir == -1 {
		t.shiftRight()
		t.pos = -1
		return
	}
	if t.pos == -1 {
		t.pos = t.ptr
		t.dir = 1
		return
	}

	t.shiftRight()
}

func (t *Textbox) resetCursorLeft() {
	if t.pos == -1 {
		return
	}
	if t.ptr >= t.pos {
		t.ptr = t.pos
	}
	t.ptr++
	t.pos = -1
}

func (t *Textbox) resetCursorRight() {
	if t.pos == -1 {
		return
	}
	if t.ptr <= t.pos {
		t.ptr = t.pos
	}
	t.pos = -1
}

func (t *Textbox) shiftLeft() {
	t.ptr--
	if t.ptr < 0 {
		t.ptr = 0
	}
}

func (t *Textbox) moveLeft() {
	t.resetCursorLeft()
	t.shiftLeft()
}

func (t *Textbox) shiftRight() {
	t.ptr++
	if t.ptr > len(t.text) {
		t.ptr = len(t.text)
	}
}

func (t *Textbox) moveRight() {
	t.resetCursorRight()
	t.shiftRight()
}

func (t *Textbox) getRuneClass(r rune) int {
	switch {
	case unicode.IsSpace(r):
		return 0
	case unicode.IsPunct(r):
		return 1
	default:
		return 2
	}
}

func (t *Textbox) shiftNextWord() {
	for t.ptr < len(t.text) {
		if !unicode.IsSpace(t.text[t.ptr]) {
			break
		}
		t.shiftRight()
	}

	if t.ptr == len(t.text) {
		return
	}

	class := t.getRuneClass(t.text[t.ptr])
	t.shiftRight()

	for t.ptr < len(t.text) {
		if class != t.getRuneClass(t.text[t.ptr]) {
			break
		}
		t.shiftRight()
	}
}

func (t *Textbox) moveNextWord() {
	t.resetCursorLeft()
	t.shiftNextWord()
}

func (t *Textbox) shiftPrevWord() {
	for t.ptr > 0 {
		t.shiftLeft()
		if !unicode.IsSpace(t.text[t.ptr]) {
			break
		}
	}

	t.shiftRight()

	if t.ptr == 0 {
		return
	}

	t.shiftLeft()
	class := t.getRuneClass(t.text[t.ptr])

	for t.ptr > 0 {
		t.shiftLeft()
		if class != t.getRuneClass(t.text[t.ptr]) {
			if t.dir == 1 {
				t.shiftRight()
			}
			break
		}
	}
}

func (t *Textbox) selectPrevWord() {
	if t.dir == -1 || t.ptr == len(t.text) {
		t.selectLeft()
	}

	for t.ptr > 0 {
		if !unicode.IsSpace(t.text[t.ptr]) {
			break
		}
		t.selectLeft()
	}

	if t.ptr == 0 {
		return
	}

	class := t.getRuneClass(t.text[t.ptr])

	for t.ptr > 0 {
		t.selectLeft()

		if class != t.getRuneClass(t.text[t.ptr]) {
			if t.dir == -1 {
				t.selectRight()
			}
			break
		}
	}
}

func (t *Textbox) movePrevWord() {
	t.resetCursorRight()
	t.shiftPrevWord()
}

func (t *Textbox) selectNextWord() {
	if t.ptr == len(t.text) {
		return
	}

	t.selectRight()

	for t.ptr < len(t.text) {
		if !unicode.IsSpace(t.text[t.ptr]) {
			break
		}
		t.selectRight()
	}

	if t.ptr == len(t.text) {
		return
	}

	class := t.getRuneClass(t.text[t.ptr])

	for t.ptr < len(t.text)-1 {
		t.selectRight()

		if class != t.getRuneClass(t.text[t.ptr]) {
			if t.dir == 1 {
				t.selectLeft()
			}
			break
		}
	}
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
	t.ptr++
}

func (t *Textbox) pop() {
	if t.pos == -1 {
		if t.ptr == 0 {
			return
		}
		t.text = append(t.text[:t.ptr-1], t.text[t.ptr:]...)
		t.shiftLeft()
		return
	}

	start, end := t.selectedArea()
	if end+1 >= len(t.text) {
		end = len(t.text) - 1
	}
	t.text = append(t.text[:start], t.text[end+1:]...)
	if t.ptr > t.pos {
		t.ptr = t.pos
	}
	t.pos = -1
}

func (t *Textbox) setLabel(label string) {
	t.label = []rune(label)
}

func (t *Textbox) setText(text string) {
	t.text = []rune(text)
}
