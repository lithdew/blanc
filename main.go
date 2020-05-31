package main

import (
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

	// layout ui

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
					//resize()
				}
			case *tcell.EventResize:
				screen.Sync()
				//resize()
			}
		}
	}()

loop:
	for {
		select {
		case <-ch:
			break loop
		case <-time.After(12 * time.Millisecond):
		}

		screen.Clear()

		//puts(screen, tcell.StyleDefault, 0, 0, fmt.Sprintf("[W]: %d", app.Width()))
		//puts(screen, tcell.StyleDefault, 0, 1, fmt.Sprintf("[H]: %d", app.Height()))

		screen.Show()
	}
}
