package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
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
	encoding.Register()

	screen, err := tcell.NewScreen()
	check(err)

	check(screen.Init())
	defer screen.Fini()

	// layout ui

	var container layout.Rect
	var rects []layout.Rect

	resize := func() {
		width, height := screen.Size()
		container = layout.Rect{W: width, H: height}

		rects, err = layout.SplitHorizontally(
			container,
			layout.Ratio(1, 3),
			layout.Ratio(1, 3),
			layout.Ratio(1, 3),
		)
		check(err)

		smaller, err := layout.SplitVertically(
			rects[1],
			layout.Ratio(1, 4),
			layout.Ratio(2, 4),
			layout.Ratio(1, 4),
		)
		check(err)

		rects = append(rects, smaller...)
	}

	ch := make(chan struct{})

	go func() {
		defer close(ch)
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlC:
					return
				case tcell.KeyCtrlL:
					screen.Sync()
					resize()
				}
			case *tcell.EventResize:
				screen.Sync()
				resize()
			}
		}
	}()

	drawrect := func(r layout.Rect) {
		box(screen, r.X, r.Y, r.X+r.W-1, r.Y+r.H-1, tcell.StyleDefault, ' ')
	}

loop:
	for {
		select {
		case <-ch:
			break loop
		case <-time.After(12 * time.Millisecond):
		}

		screen.Clear()

		drawrect(container)
		for _, rect := range rects {
			drawrect(rect)
		}

		//puts(screen, tcell.StyleDefault, 0, 0, fmt.Sprintf("[W]: %d", app.Width()))
		//puts(screen, tcell.StyleDefault, 0, 1, fmt.Sprintf("[H]: %d", app.Height()))

		screen.Show()
	}
}
