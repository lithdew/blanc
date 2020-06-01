package main

import (
	"github.com/gdamore/tcell"
	"github.com/lithdew/asciigraph"
	"github.com/lithdew/blanc/layout"
	"log"
	"strings"
	"time"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	s, err := tcell.NewScreen()
	check(err)

	check(s.Init())
	defer s.Fini()

	ch := make(chan struct{})

	in := newTextbox()
	in.setLabel(">>> ")

	go func() {
		defer close(ch)
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlC:
					return
				case tcell.KeyCtrlL:
					s.Sync()
				case tcell.KeyDelete, tcell.KeyDEL:
					in.pop()
				case tcell.KeyCtrlW:
					//m := ev.Modifiers()
					//if m & tcell.ModShift!= 0 && m & tcell.ModAlt != 0 {
					//
					//}
					in.moveNextWord()
				case tcell.KeyCtrlU:
					in.moveToEnd()
				case tcell.KeyLeft:
					m := ev.Modifiers()
					if m&tcell.ModCtrl != 0 {
						in.movePrevWord()
					} else if m&tcell.ModShift != 0 {
						in.selectLeft()
					} else {
						in.moveLeft()
					}
				case tcell.KeyRight:
					m := ev.Modifiers()
					if m&tcell.ModCtrl != 0 {
						in.moveNextWord()
					} else if m&tcell.ModShift != 0 {
						in.selectRight()
					} else {
						in.moveRight()
					}
				case tcell.KeyRune:
					if ev.Key() == tcell.KeyRune {
						in.push(ev.Rune())
					}
				}
			case *tcell.EventResize:
				s.Sync()
			}
		}
	}()

	sh := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	sb := sh.Reverse(true)

loop:
	for {
		select {
		case <-ch:
			break loop
		case <-time.After(12 * time.Millisecond):
		}

		w, h := s.Size()

		scr := layout.Rect{X: 0, Y: 0, W: w, H: h}

		// header

		hdr := layout.Rect{W: w, H: 1}.Align(scr, layout.Top|layout.Left)
		clear(s, sh, hdr.X, hdr.Y, hdr.X+hdr.W-1, hdr.Y+hdr.H-1)

		txt := "flatend."
		rect := layout.Text(txt).Align(hdr, layout.Left).ShiftLeft(1)

		puts(s, sh.Bold(true), rect.X, rect.Y, txt)

		// body

		bdy := scr.PadVertical(1)
		clear(s, sb, bdy.X, bdy.Y, bdy.X+bdy.W-1, bdy.Y+bdy.H-1)

		graph := bdy.Pad(4)

		data := []float64{3, 4, 9, 6, 2, 4, 5, 8, 5, 10, 2, 7, 2, 5, 6}
		txt = asciigraph.Plot(data, asciigraph.Width(graph.W), asciigraph.Height(graph.H))
		for i, c := range strings.Split(txt, "\n") {
			puts(s, sb.Reverse(false), graph.X, graph.Y+i, c)
		}

		// footer

		renderFooter(s, scr, in)

		s.Show()
	}
}

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
	cs := tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite) // cursor style

	j := 0
	for i := 0; i < len(t.label) && r.X+j < r.X+r.W-1; i, j = i+1, j+1 {
		s.SetContent(r.X+j, r.Y, t.label[i], nil, style)
	}

	start, end := t.cursorPos()

	for i := 0; i < len(t.text) && r.X+j < r.X+r.W-1; i, j = i+1, j+1 {
		if i >= start && i <= end {
			s.SetContent(r.X+j, r.Y, t.text[i], nil, cs)
		} else {
			s.SetContent(r.X+j, r.Y, t.text[i], nil, style)
		}
	}

	if t.ptr == len(t.text) && t.pos == -1 { // render cursor
		s.SetContent(r.X+j, r.Y, tcell.RuneBlock, nil, cs.Reverse(true))
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
	if t.ptr == len(t.text)-1 {
		if t.pos == -1 {
			t.pos = t.ptr
			t.ptr = len(t.text)
		}
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
			t.ptr = t.pos - 1
		}
		t.pos = -1
	}
	if t.ptr == len(t.text) {
		return
	}
	t.ptr++
}

func (t *Textbox) movePrevWord() {
	t.pos = -1

	i := t.ptr - 2
	if i == len(t.text) {
		i--
	}
	for ; i >= 0; i-- {
		if t.text[i] == ' ' {
			t.ptr = i + 1
			return
		}
	}
	t.ptr = 0
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

	// todo
}

func (t *Textbox) setLabel(label string) {
	t.label = []rune(label)
}

func (t *Textbox) setText(text string) {
	t.text = []rune(text)
}

func renderFooter(s tcell.Screen, scr layout.Rect, in *Textbox) {
	style := tcell.StyleDefault.Reverse(true)

	ftr := layout.Rect{W: scr.W, H: 1}.Align(scr, layout.Bottom|layout.Left)
	clear(s, style, ftr.X, ftr.Y, ftr.X+ftr.W-1, ftr.Y+ftr.H-1)

	in.render(s, style, ftr.PadLeft(1))

	//menuRect := layout.Rect{X: textRect.Right() + 1, Y: textRect.Y - 10, W: 30, H: 10}
	//clear(s, tcell.StyleDefault.Background(tcell.ColorBlue), menuRect.Left(), menuRect.Top(), menuRect.Right(), menuRect.Bottom())
}
