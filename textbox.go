package main

import (
	"github.com/gdamore/tcell"
	"github.com/lithdew/blanc/layout"
	"github.com/mattn/go-runewidth"
	"unicode"
)

type Textbox struct {
	selected tcell.Style
	style    StyleFunc

	label Text
	text  Text

	buf  []rune
	size int

	ptr int // cursor start index
	pos int // cursor end index
	dir int // -1 is left, 1 is right
}

func newTextbox() *Textbox {
	return &Textbox{
		selected: tcell.StyleDefault.Reverse(true),

		pos: -1,
		dir: -1,
	}
}

func (t *Textbox) cursorX(r layout.Rect) int {
	x := r.X + t.label.width
	for i := 0; i < t.ptr; i++ {
		x += runewidth.RuneWidth(t.buf[i])
	}
	return x
}

func (t *Textbox) HandleEvent(e tcell.Event) bool {
	ev, ok := e.(*tcell.EventKey)
	if !ok {
		return false
	}

	m := ev.Modifiers()

	ctrl := m&tcell.ModCtrl != 0
	shift := m&tcell.ModShift != 0

	switch ev.Key() {
	case tcell.KeyDelete, tcell.KeyDEL:
		t.pop()
		return true
	case tcell.KeyCtrlA:
		t.selectAll()
		return true
	case tcell.KeyCtrlW:
		t.moveNextWord()
		return true
	case tcell.KeyCtrlU:
		t.moveToEnd()
		return true
	case tcell.KeyLeft:
		switch {
		case ctrl && shift:
			t.selectPrevWord()
		case ctrl:
			t.movePrevWord()
		case shift:
			t.selectLeft()
		default:
			t.moveLeft()
		}
		return true
	case tcell.KeyRight:
		switch {
		case ctrl && shift:
			t.selectNextWord()
		case ctrl:
			t.moveNextWord()
		case shift:
			t.selectRight()
		default:
			t.moveRight()
		}
		return true
	case tcell.KeyRune:
		if ev.Key() == tcell.KeyRune {
			t.push(ev.Rune())
		}
		return true
	default:
		return false
	}
}

func (t *Textbox) getText() string {
	return string(t.buf)
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

func (t *Textbox) Draw(s tcell.Screen, r layout.Rect) {
	if r.H < 1 {
		return
	}

	t.label.Draw(s, r)

	r = r.PadLeft(t.label.Width())

	var styleFunc StyleFunc

	if t.style != nil {
		if t.selected != tcell.StyleDefault {
			start, end := t.selectedArea()

			styleFunc = func(i int) tcell.Style {
				if t.pos != -1 && i >= start && i <= end {
					return t.selected
				}
				return t.style(i)
			}
		} else {
			styleFunc = t.style
		}
	}

	t.text.SetStyleFunc(styleFunc)

	t.text.SetText(string(t.buf))
	t.text.Draw(s, r)

	if t.pos == -1 {
		s.ShowCursor(r.X+t.ptr, r.Y)
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
	if t.ptr == len(t.buf) {
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
	if t.ptr > len(t.buf) {
		t.ptr = len(t.buf)
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
	for t.ptr < len(t.buf) {
		if !unicode.IsSpace(t.buf[t.ptr]) {
			break
		}
		t.shiftRight()
	}

	if t.ptr == len(t.buf) {
		return
	}

	class := t.getRuneClass(t.buf[t.ptr])
	t.shiftRight()

	for t.ptr < len(t.buf) {
		if class != t.getRuneClass(t.buf[t.ptr]) {
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
		if !unicode.IsSpace(t.buf[t.ptr]) {
			break
		}
	}

	t.shiftRight()

	if t.ptr == 0 {
		return
	}

	t.shiftLeft()
	class := t.getRuneClass(t.buf[t.ptr])

	for t.ptr > 0 {
		t.shiftLeft()
		if class != t.getRuneClass(t.buf[t.ptr]) {
			if t.dir == 1 {
				t.shiftRight()
			}
			break
		}
	}
}

func (t *Textbox) selectPrevWord() {
	if t.dir == -1 || t.ptr == len(t.buf) {
		t.selectLeft()
	}

	for t.ptr > 0 {
		if !unicode.IsSpace(t.buf[t.ptr]) {
			break
		}
		t.selectLeft()
	}

	if t.ptr == 0 {
		return
	}

	class := t.getRuneClass(t.buf[t.ptr])

	for t.ptr > 0 {
		t.selectLeft()

		if class != t.getRuneClass(t.buf[t.ptr]) {
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
	if t.ptr == len(t.buf) {
		return
	}

	t.selectRight()

	for t.ptr < len(t.buf) {
		if !unicode.IsSpace(t.buf[t.ptr]) {
			break
		}
		t.selectRight()
	}

	if t.ptr == len(t.buf) {
		return
	}

	class := t.getRuneClass(t.buf[t.ptr])

	for t.ptr < len(t.buf)-1 {
		t.selectRight()

		if class != t.getRuneClass(t.buf[t.ptr]) {
			if t.dir == 1 {
				t.selectLeft()
			}
			break
		}
	}
}

func (t *Textbox) selectAll() {
	if len(t.buf) == 0 {
		return
	}
	t.ptr = len(t.buf)
	t.pos = 0
}

func (t *Textbox) moveToEnd() {
	t.pos = -1
	for i := t.ptr; i < len(t.buf); i++ {
		if t.buf[i] == '\n' {
			t.ptr = i
			return
		}
	}
	t.ptr = len(t.buf)
}

func (t *Textbox) insert(pos int, r rune) {
	t.size += runewidth.RuneWidth(r)
	t.buf = append(t.buf[:pos], append([]rune{r}, t.buf[pos:]...)...)
}

func (t *Textbox) push(r rune) {
	if t.pos != -1 {
		t.pop()
	}

	t.insert(t.ptr, r)
	t.ptr++
}

func (t *Textbox) pop() {
	if t.pos == -1 {
		if t.ptr == 0 {
			return
		}
		t.size -= runewidth.RuneWidth(t.buf[t.ptr-1])
		t.buf = append(t.buf[:t.ptr-1], t.buf[t.ptr:]...)
		t.shiftLeft()
		return
	}

	start, end := t.selectedArea()
	if end+1 >= len(t.buf) {
		end = len(t.buf) - 1
	}
	for i := start; i < end+1; i++ {
		t.size -= runewidth.RuneWidth(t.buf[i])
	}
	t.buf = append(t.buf[:start], t.buf[end+1:]...)
	if t.ptr > t.pos {
		t.ptr = t.pos
	}
	t.pos = -1
}

func (t *Textbox) setLabel(label string) {
	t.label.SetText(label)
}

func (t *Textbox) setText(text string) {
	t.buf = []rune(text)

	t.size = 0
	for i := 0; i < len(t.buf); i++ {
		t.size += runewidth.RuneWidth(t.buf[i])
	}
}

func (t *Textbox) SetLabelStyleFunc(style StyleFunc) { t.label.SetStyleFunc(style) }
func (t *Textbox) SetLabelStyle(style tcell.Style)   { t.label.SetStyle(style) }

func (t *Textbox) SetTextStyle(style tcell.Style)     { t.style = func(int) tcell.Style { return style } }
func (t *Textbox) SetTextStyleFunc(style StyleFunc)   { t.style = style }
func (t *Textbox) SetSelectedStyle(style tcell.Style) { t.selected = style }
