package main

import (
	"github.com/gdamore/tcell"
	"github.com/lithdew/blanc/layout"
	"log"
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
				case tcell.KeyCtrlA:
					in.selectAll()
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

					ctrl := m&tcell.ModCtrl != 0
					shift := m&tcell.ModShift != 0

					if ctrl && shift {
						in.selectPrevWord()
					} else if ctrl {
						in.movePrevWord()
					} else if shift {
						in.selectLeft()
					} else {
						in.moveLeft()
					}
				case tcell.KeyRight:
					m := ev.Modifiers()

					ctrl := m&tcell.ModCtrl != 0
					shift := m&tcell.ModShift != 0

					if ctrl && shift {
						in.selectNextWord()
					} else if ctrl {
						in.moveNextWord()
					} else if shift {
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

		//data := []float64{3, 4, 9, 6, 2, 4, 5, 8, 5, 10, 2, 7, 2, 5, 6}
		//txt = asciigraph.Plot(data, asciigraph.Width(graph.W), asciigraph.Height(graph.H))
		//for i, c := range strings.Split(txt, "\n") {
		//	puts(s, sb.Reverse(true), graph.X, graph.Y+i, c)
		//}

		sentence := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent sollicitudin augue nisi, vel euismod mi eleifend et. Nulla maximus magna id ex malesuada vestibulum semper nec dui. Duis sagittis scelerisque augue et eleifend. Nam quis est urna. Suspendisse non sapien pellentesque, porta dui quis, hendrerit ex. Vestibulum tempor efficitur nisi quis accumsan. Vestibulum nisl magna, dignissim at eros ac, maximus scelerisque mauris. Vivamus consequat metus justo, eget venenatis urna finibus quis. Curabitur congue feugiat ipsum, sed lacinia turpis aliquam eu. Mauris rhoncus lectus id erat luctus ultricies. Fusce sodales urna eu purus ornare consectetur. In vitae leo dignissim, tincidunt velit ut, viverra velit. Quisque vel nibh nec mi bibendum tempor sit amet vitae nisl. In maximus odio eget tristique imperdiet. Fusce id nunc ut arcu ultrices convallis. Pellentesque."
		tst := NewText(sentence)
		tst.SetWrap(true)
		tst.Draw(s, graph)

		// footer

		renderFooter(s, scr, in)

		s.Show()
	}
}

func renderFooter(s tcell.Screen, scr layout.Rect, in *Textbox) {
	style := tcell.StyleDefault.Reverse(true)

	ftr := layout.Rect{W: scr.W, H: 1}.Align(scr, layout.Bottom|layout.Left)
	clear(s, style, ftr.X, ftr.Y, ftr.X+ftr.W-1, ftr.Y+ftr.H-1)

	inRect := ftr.PadLeft(1)
	in.render(s, style, inRect)

	if len(in.getText()) > 0 {
		menuRect := layout.Rect{X: in.cursorX(inRect) + 1, Y: ftr.Y - 10, W: 30, H: 10}
		clear(s, tcell.StyleDefault.Background(tcell.ColorBlue), menuRect.Left(), menuRect.Top(), menuRect.Right(), menuRect.Bottom())
	}
}
