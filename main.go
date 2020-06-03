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

var inputs []InputListener

func eventLoop(s tcell.Screen, ch chan<- struct{}) {
	defer close(ch)
	for {
		ev := s.PollEvent()

		handled := false
		for _, input := range inputs {
			handled = input.HandleEvent(ev)
			if handled {
				break
			}
		}

		if handled {
			continue
		}

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				return
			case tcell.KeyCtrlL:
				s.Sync()
			}
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func main() {
	s, err := tcell.NewScreen()
	check(err)

	check(s.Init())
	defer s.Fini()

	ch := make(chan struct{})

	input := newTextbox()
	input.setLabel(">>> ")
	inputs = append(inputs, input)

	go eventLoop(s, ch)

	titleStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)

	titleText := "flatend."
	title := NewText(titleText)
	title.SetStyle(titleStyle.Bold(true))

	contentStyle := titleStyle.Reverse(true)

	graph := NewASCIIGraph()
	graph.SetStyle(contentStyle)

	content := NewText(
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit." +
			" Praesent sollicitudin augue nisi, vel euismod mi eleifend et. Nulla maximus " +
			"magna id ex malesuada vestibulum semper nec dui. Duis sagittis scelerisque augue" +
			" et eleifend. Nam quis est urna. Suspendisse non sapien pellentesque, porta dui" +
			" quis, hendrerit ex. Vestibulum tempor efficitur nisi quis accumsan. Vestibulum" +
			" nisl magna, dignissim at eros ac, maximus scelerisque mauris. Vivamus consequat" +
			" metus justo, eget venenatis urna finibus quis. Curabitur congue feugiat ipsum, " +
			"sed lacinia turpis aliquam eu. Mauris rhoncus lectus id erat luctus ultricies. " +
			"Fusce sodales urna eu purus ornare consectetur. In vitae leo dignissim, tincidunt" +
			" velit ut, viverra velit. Quisque vel nibh nec mi bibendum tempor sit amet vitae nisl." +
			" In maximus odio eget tristique imperdiet. Fusce id nunc ut arcu ultrices convallis." +
			" Pellentesque.",
	)
	content.SetWrap(true)
	content.SetStyle(contentStyle)

loop:
	for {
		select {
		case <-ch:
			break loop
		case <-time.After(40 * time.Millisecond):
		}

		w, h := s.Size()

		screenRect := layout.Rect{X: 0, Y: 0, W: w, H: h}

		// header

		headerRect := layout.Rect{W: w, H: 1}.Align(screenRect, layout.Top|layout.Left)
		clear(s, titleStyle, headerRect.X, headerRect.Y, headerRect.X+headerRect.W-1, headerRect.Y+headerRect.H-1)
		title.Draw(s, layout.Text(title.Text()).Align(headerRect, layout.Left).ShiftLeft(1))

		// body

		bodyRect := screenRect.PadVertical(1)
		clear(s, contentStyle, bodyRect.X, bodyRect.Y, bodyRect.X+bodyRect.W-1, bodyRect.Y+bodyRect.H-1)

		content.Draw(s, bodyRect.Pad(4))
		//graph.Draw(s, bodyRect.Pad(4))

		// footer

		renderFooter(s, screenRect, input)

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
		items := []string{"hello", "world", "testing"}

		menuRect := layout.Rect{X: in.cursorX(inRect) + 1, Y: ftr.Y - len(items), W: 30, H: len(items)}
		clear(s, tcell.StyleDefault.Background(tcell.ColorGray), menuRect.Left(), menuRect.Top(), menuRect.Right(), menuRect.Bottom())

		first := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
		second := tcell.StyleDefault.Background(tcell.ColorDimGray).Foreground(tcell.ColorWhite)

		for i := range items {
			itemRect := layout.Rect{X: menuRect.X, Y: menuRect.Y + i, W: menuRect.W, H: 1}

			view := NewText(" " + items[i])

			if i%2 == 1 {
				view.SetStyle(first)
			} else {
				view.SetStyle(second)
			}

			view.Draw(s, itemRect)
		}
	}
}
