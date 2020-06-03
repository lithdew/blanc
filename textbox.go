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

func (t *Textbox) CursorX(r layout.Rect) int {
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
		t.Pop()
		return true
	case tcell.KeyCtrlA:
		t.SelectAll()
		return true
	case tcell.KeyCtrlW:
		t.MoveToNextWord()
		return true
	case tcell.KeyCtrlU:
		t.MoveToEnd()
		return true
	case tcell.KeyLeft:
		switch {
		case ctrl && shift:
			t.SelectPrevWord()
		case ctrl:
			t.MoveToPrevWord()
		case shift:
			t.SelectLeft()
		default:
			t.MoveLeft()
		}
		return true
	case tcell.KeyRight:
		switch {
		case ctrl && shift:
			t.SelectNextWord()
		case ctrl:
			t.MoveToNextWord()
		case shift:
			t.SelectRight()
		default:
			t.MoveRight()
		}
		return true
	case tcell.KeyRune:
		if ev.Key() == tcell.KeyRune {
			t.Push(ev.Rune())
		}
		return true
	default:
		return false
	}
}

func (t Textbox) Text() string  { return string(t.buf) }
func (t Textbox) Label() string { return t.label.Text() }

func (t *Textbox) Selected() (start int, end int) {
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
			start, end := t.Selected()

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

func (t *Textbox) SelectLeft() {
	if t.pos != -1 && t.ptr == t.pos && t.dir == 1 {
		t.pos = -1
		return
	}
	if t.pos == -1 {
		t.pos = t.ptr - 1
		t.dir = -1
	}
	t.ShiftLeft()
}

func (t *Textbox) SelectRight() {
	if t.ptr == len(t.buf) {
		return
	}

	if t.pos != -1 && t.ptr == t.pos && t.dir == -1 {
		t.ShiftRight()
		t.pos = -1
		return
	}
	if t.pos == -1 {
		t.pos = t.ptr
		t.dir = 1
		return
	}

	t.ShiftRight()
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

func (t *Textbox) ShiftLeft() {
	t.ptr--
	if t.ptr < 0 {
		t.ptr = 0
	}
}

func (t *Textbox) ShiftRight() {
	t.ptr++
	if t.ptr > len(t.buf) {
		t.ptr = len(t.buf)
	}
}

func (t *Textbox) MoveLeft() {
	t.resetCursorLeft()
	t.ShiftLeft()
}

func (t *Textbox) MoveRight() {
	t.resetCursorRight()
	t.ShiftRight()
}

func (t *Textbox) MoveToNextWord() {
	t.resetCursorLeft()

	for t.ptr < len(t.buf) {
		if !unicode.IsSpace(t.buf[t.ptr]) {
			break
		}
		t.ShiftRight()
	}

	if t.ptr == len(t.buf) {
		return
	}

	class := t.getRuneClass(t.buf[t.ptr])
	t.ShiftRight()

	for t.ptr < len(t.buf) {
		if class != t.getRuneClass(t.buf[t.ptr]) {
			break
		}
		t.ShiftRight()
	}
}

func (t *Textbox) SelectPrevWord() {
	if t.dir == -1 || t.ptr == len(t.buf) {
		t.SelectLeft()
	}

	for t.ptr > 0 {
		if !unicode.IsSpace(t.buf[t.ptr]) {
			break
		}
		t.SelectLeft()
	}

	if t.ptr == 0 {
		return
	}

	class := t.getRuneClass(t.buf[t.ptr])

	for t.ptr > 0 {
		t.SelectLeft()

		if class != t.getRuneClass(t.buf[t.ptr]) {
			if t.dir == -1 {
				t.SelectRight()
			}
			break
		}
	}
}

func (t *Textbox) MoveToPrevWord() {
	t.resetCursorRight()

	for t.ptr > 0 {
		t.ShiftLeft()
		if !unicode.IsSpace(t.buf[t.ptr]) {
			break
		}
	}

	t.ShiftRight()

	if t.ptr == 0 {
		return
	}

	t.ShiftLeft()
	class := t.getRuneClass(t.buf[t.ptr])

	for t.ptr > 0 {
		t.ShiftLeft()
		if class != t.getRuneClass(t.buf[t.ptr]) {
			if t.dir == 1 {
				t.ShiftRight()
			}
			break
		}
	}
}

func (t *Textbox) SelectNextWord() {
	if t.ptr == len(t.buf) {
		return
	}

	t.SelectRight()

	for t.ptr < len(t.buf) {
		if !unicode.IsSpace(t.buf[t.ptr]) {
			break
		}
		t.SelectRight()
	}

	if t.ptr == len(t.buf) {
		return
	}

	class := t.getRuneClass(t.buf[t.ptr])

	for t.ptr < len(t.buf)-1 {
		t.SelectRight()

		if class != t.getRuneClass(t.buf[t.ptr]) {
			if t.dir == 1 {
				t.SelectLeft()
			}
			break
		}
	}
}

func (t *Textbox) SelectAll() {
	if len(t.buf) == 0 {
		return
	}
	t.ptr = len(t.buf)
	t.pos = 0
}

func (t *Textbox) MoveToEnd() {
	t.pos = -1
	for i := t.ptr; i < len(t.buf); i++ {
		if t.buf[i] == '\n' {
			t.ptr = i
			return
		}
	}
	t.ptr = len(t.buf)
}

func (t *Textbox) Insert(pos int, r rune) {
	t.size += runewidth.RuneWidth(r)
	t.buf = append(t.buf[:pos], append([]rune{r}, t.buf[pos:]...)...)
}

func (t *Textbox) Push(r rune) {
	if t.pos != -1 {
		t.Pop()
	}

	t.Insert(t.ptr, r)
	t.ptr++
}

func (t *Textbox) Pop() {
	if t.pos == -1 {
		if t.ptr == 0 {
			return
		}
		t.size -= runewidth.RuneWidth(t.buf[t.ptr-1])
		t.buf = append(t.buf[:t.ptr-1], t.buf[t.ptr:]...)
		t.ShiftLeft()
		return
	}

	start, end := t.Selected()
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

func (t *Textbox) SetLabel(label string) { t.label.SetText(label) }
func (t *Textbox) SetText(text string) {
	t.buf = []rune(text)

	t.size = 0
	for i := 0; i < len(t.buf); i++ {
		t.size += runewidth.RuneWidth(t.buf[i])
	}
}

func (t *Textbox) SetLabelStyleFunc(style StyleFunc) { t.label.SetStyleFunc(style) }
func (t *Textbox) SetLabelStyle(style tcell.Style)   { t.label.SetStyle(style) }

func (t *Textbox) SetTextStyleFunc(style StyleFunc) { t.style = style }
func (t *Textbox) SetTextStyle(style tcell.Style) {
	t.SetTextStyleFunc(func(int) tcell.Style { return style })
}

func (t *Textbox) SetSelectedStyle(style tcell.Style) { t.selected = style }
