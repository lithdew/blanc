package main

import (
	"fmt"
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

	var (
		container layout.Rect
		panels    []layout.Rect
	)

	resize := func() {
		w, h := s.Size()

		container = layout.Rect{W: w, H: h}

		panels, err = layout.SplitHorizontally(
			container,
			layout.Length(25),
			layout.Min(1),
			layout.Length(25),
		)
		check(err)

		middle, err := layout.SplitVertically(
			panels[1],
			layout.Length(4),
			layout.Ratio(1, 3),
			layout.Ratio(1, 3),
		)
		check(err)

		panels = append(panels, middle...)
	}

	ch := make(chan struct{})

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
					resize()
				}
			case *tcell.EventResize:
				s.Sync()
				resize()
			}
		}
	}()

	style := tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)

	R := func(r layout.Rect) {
		box(s, r.X, r.Y, r.X+r.W-1, r.Y+r.H-1, style, TcellFrame)
	}

loop:
	for {
		select {
		case <-ch:
			break loop
		case <-time.After(12 * time.Millisecond):
		}

		s.Clear()

		R(container)
		for _, rect := range panels {
			R(rect)
		}

		width, height := s.Size()

		txt := fmt.Sprintf("[W]: %d", width)
		rect := layout.Align(panels[0], layout.Text(txt), layout.Center).PadLeft(1)
		puts(s, style, rect.X, rect.Y, txt)

		txt = fmt.Sprintf("[H]: %d", height)
		rect = layout.Align(panels[2], layout.Text(txt), layout.Center).PadLeft(1)
		puts(s, style, rect.X, rect.Y, txt)

		s.Show()
	}
}
