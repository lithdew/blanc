package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/lithdew/casso"
	"log"
	"time"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

var solver = casso.NewSolver()
var parent = NewBox(solver)
var left = NewBox(solver)
var right = NewBox(solver)

func InitLayout() *Layout {
	a := New(solver)
	a.Required(Inside(parent, left, 2)...)
	a.Required(Inside(parent, right, 2)...)
	a.Medium(
		FillX(parent, left, 0.5),
		FillX(parent, right, 0.5),
		FillY(parent, left, 1),
		FillY(parent, right, 1),
	)
	a.Required(SpaceBetween(left, right, 2))
	a.Required(SameWidth(left, right))
	//a.Required(MinWidth(left, 50))
	//a.Required(MaxWidth(left, 50))
	check(a.Finalize())
	return &a
}

func main() {
	encoding.Register()

	screen, err := tcell.NewScreen()
	check(err)

	check(screen.Init())
	defer screen.Fini()

	parent.Fixed(casso.Medium)

	width, height := screen.Size()
	check(parent.SetX(0))
	check(parent.SetY(0))
	check(parent.SetW(float64(width)))
	check(parent.SetH(float64(height)))

	layout := InitLayout()

	ch := make(chan struct{})

	go func() {
		defer close(ch)
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlH:
					if layout != nil {
						layout.Destroy()
						layout = nil
					} else {
						layout = InitLayout()
					}
				case tcell.KeyCtrlC:
					return
				case tcell.KeyCtrlL:
					screen.Sync()

					width, height := screen.Size()
					check(parent.SetW(float64(width)))
					check(parent.SetH(float64(height)))
				}
			case *tcell.EventResize:
				screen.Sync()

				width, height := screen.Size()
				check(parent.SetW(float64(width)))
				check(parent.SetH(float64(height)))
			}
		}
	}()

	drawbox := func(b Box) {
		box(screen, int(b.X()), int(b.Y()), int(b.X()+b.W())-1, int(b.Y()+b.H())-1, tcell.StyleDefault, ' ')
	}

loop:
	for {
		select {
		case <-ch:
			break loop
		case <-time.After(12 * time.Millisecond):
		}

		screen.Clear()

		drawbox(parent)
		drawbox(left)
		drawbox(right)

		screen.Show()

		//puts(screen, tcell.StyleDefault, 0, 0, fmt.Sprintf("[W]: %d", app.Width()))
		//puts(screen, tcell.StyleDefault, 0, 1, fmt.Sprintf("[H]: %d", app.Height()))

		//app.Render(screen)
	}
}
