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

func renderFooter(s tcell.Screen, scr layout.Rect, in *Textbox) {
	style := tcell.StyleDefault.Reverse(true)

	ftr := layout.Rect{W: scr.W, H: 1}.Align(scr, layout.Bottom|layout.Left)
	clear(s, style, ftr.X, ftr.Y, ftr.X+ftr.W-1, ftr.Y+ftr.H-1)

	in.render(s, style, ftr.PadLeft(1))

	//menuRect := layout.Rect{X: textRect.Right() + 1, Y: textRect.Y - 10, W: 30, H: 10}
	//clear(s, tcell.StyleDefault.Background(tcell.ColorBlue), menuRect.Left(), menuRect.Top(), menuRect.Right(), menuRect.Bottom())
}
