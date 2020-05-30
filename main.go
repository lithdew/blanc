package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
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

	app, err := NewApp(screen)
	check(err)

	v := NewLayout()
	check(app.Add(v))

	check(v.SetY(5))
	check(v.SetPaddingTop(2))
	check(v.SetPaddingBottom(2))
	check(v.SetPaddingLeft(2))
	check(v.SetPaddingRight(2))
	check(v.FillX())
	check(v.FillY())

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
					sw, sh := screen.Size()
					check(app.Resize(sw, sh))
				}
			case *tcell.EventResize:
				screen.Sync()
				sw, sh := screen.Size()
				check(app.Resize(sw, sh))
			}
		}
	}()

loop:
	for {
		select {
		case <-ch:
			break loop
		case <-time.After(24 * time.Millisecond):
		}

		screen.Clear()

		puts(screen, tcell.StyleDefault, 0, 0, fmt.Sprintf("[W]: %d", app.Width()))
		puts(screen, tcell.StyleDefault, 0, 1, fmt.Sprintf("[H]: %d", app.Height()))

		app.Render(screen)
	}
}
