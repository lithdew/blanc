package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/lithdew/blanc"
	"github.com/lithdew/blanc/layout"
	"log"
	"time"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

var handlers []tcell.EventHandler

func handleEvent(ev tcell.Event) bool {
	for _, handler := range handlers {
		if handler.HandleEvent(ev) {
			return true
		}
	}
	return false
}

func eventLoop(screen tcell.Screen, ch chan<- struct{}) {
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
			}
		case *tcell.EventResize:
			screen.Sync()
		}

		handleEvent(ev)
	}
}

func main() {
	encoding.Register()

	screen, err := tcell.NewScreen()
	check(err)

	check(screen.Init())
	defer screen.Fini()

	inputStyle := tcell.StyleDefault.Reverse(true)
	selectedStyle := tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite)

	input := blanc.NewTextbox()
	input.SetLabel(">>> ")
	input.SetTextStyle(inputStyle)
	input.SetLabelStyle(inputStyle)
	input.SetSelectedStyle(selectedStyle)

	handlers = append(handlers, input)

	titleStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)

	title := blanc.NewText("flatend.")
	title.SetStyle(titleStyle)

	contentStyle := titleStyle.Reverse(true)

	//graph := NewASCIIGraph()
	//graph.SetStyle(contentStyle)

	//content := NewText(
	//	"Lorem ipsum dolor sit amet, consectetur adipiscing elit." +
	//		" Praesent sollicitudin augue nisi, vel euismod mi eleifend et. Nulla maximus " +
	//		"magna id ex malesuada vestibulum semper nec dui. Duis sagittis scelerisque augue" +
	//		" et eleifend. Nam quis est urna. Suspendisse non sapien pellentesque, porta dui" +
	//		" quis, hendrerit ex. Vestibulum tempor efficitur nisi quis accumsan. Vestibulum" +
	//		" nisl magna, dignissim at eros ac, maximus scelerisque mauris. Vivamus consequat" +
	//		" metus justo, eget venenatis urna finibus quis. Curabitur congue feugiat ipsum, " +
	//		"sed lacinia turpis aliquam eu. Mauris rhoncus lectus id erat luctus ultricies. " +
	//		"Fusce sodales urna eu purus ornare consectetur. In vitae leo dignissim, tincidunt" +
	//		" velit ut, viverra velit. Quisque vel nibh nec mi bibendum tempor sit amet vitae nisl." +
	//		" In maximus odio eget tristique imperdiet. Fusce id nunc ut arcu ultrices convallis." +
	//		" Pellentesque.",
	//)
	//content.SetWrap(true)
	//content.SetStyle(contentStyle)

	ch := make(chan struct{})
	go eventLoop(screen, ch)

	for {
		select {
		case <-ch:
			return
		case <-time.After(40 * time.Millisecond):
		}

		w, h := screen.Size()

		screenRect := layout.Rect{X: 0, Y: 0, W: w, H: h}

		// header

		headerRect := screenRect.Position(layout.Top | layout.Left).WidthOf(screenRect).Height(1)
		blanc.Clear(screen, titleStyle, headerRect)
		title.Draw(screen, layout.Text(title.Text()).PositionTo(headerRect, layout.Left).MoveRight(1))

		tabRect := headerRect.Align(layout.Top | layout.Center).Width(headerRect.W / 2)

		tab1 := blanc.NewText(" Backend ")
		tab3 := blanc.NewText(" Metrics ")

		tab1Rect := layout.Text(tab1.Text()).AlignTo(tabRect, layout.Left)
		tab3Rect := layout.Text(tab3.Text()).AlignTo(tab1Rect, layout.Right)

		tab1.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
		tab3.SetStyle(tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack))

		tab1.Draw(screen, tab1Rect)
		tab3.Draw(screen, tab3Rect)

		// body

		bodyRect := screenRect.PadVertical(1)
		blanc.Clear(screen, contentStyle, bodyRect)

		//content.Draw(screen, bodyRect.Pad(4))
		//graph.Draw(screen, bodyRect.Pad(4))

		items := [][]string{
			{"id", "type", "status", "params"},
			{"0", "http", "ready", "9000"},
			{"1", "post", "ready", "/post/new"},
			{"2", "get", "ready", "/post/:id"},
			{"3", "jsonrpc", "ready", "hello_world"},
			{"4", "cron", "ready", "hello_world @every 10s"},
			{"5", "sql", "ready", "select * from posts where id = :id"},
			{"6", "csv", "ready", "/home/kenta/Desktop/database.csv"},
		}

		for i := 0; i < len(items); i++ {
			itemRect := layout.Rect{H: 1}.MoveDown(2 + i).WidthOf(screenRect)
			blanc.Clear(screen, contentStyle, itemRect)

			rects, err := layout.SplitHorizontally(itemRect.MoveRight(1),
				layout.Ratio(1, 12),
				layout.Ratio(2, 12),
				layout.Ratio(2, 12),
				layout.Percentage(80),
			)
			check(err)

			for j := 0; j < len(items[i]); j++ {
				t := blanc.NewText(items[i][j])
				if i == 0 {
					t.SetStyle(contentStyle.Underline(true))
				} else {
					t.SetStyle(contentStyle)
				}
				t.Draw(screen, rects[j])
			}

			for j := 0; j < len(items[i]); j++ {

			}
		}

		// footer

		renderFooter(screen, screenRect, input)

		screen.Show()
	}
}

func renderFooter(s tcell.Screen, screenRect layout.Rect, input *blanc.Textbox) {
	footerRect := screenRect.Position(layout.Bottom | layout.Left).WidthOf(screenRect).Height(1)
	blanc.Clear(s, tcell.StyleDefault.Reverse(true), footerRect)

	inputRect := footerRect.PadLeft(1)
	input.Draw(s, inputRect)

	renderMenu(s, input, inputRect)
}

func renderMenu(s tcell.Screen, input *blanc.Textbox, inputRect layout.Rect) {
	if len(input.Text()) == 0 {
		return
	}

	items := []string{"hello", "world", "testing"}

	bg := tcell.StyleDefault.Background(tcell.ColorGray)
	first := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)
	second := tcell.StyleDefault.Background(tcell.ColorDimGray).Foreground(tcell.ColorWhite)

	style := func(i int) tcell.Style {
		if i == -1 {
			return bg
		}
		if i%2 == 1 {
			return first
		}
		return second
	}

	menuRect := layout.Rect{
		X: input.CursorX(inputRect) + 1,
		Y: inputRect.Top() - len(items),
		W: 30,
		H: len(items),
	}

	blanc.Clear(s, style(-1), menuRect)

	for i := range items {
		itemRect := menuRect.Position(layout.Top | layout.Left).MoveDown(i).WidthOf(menuRect).Height(1)
		itemStyle := style(i)

		blanc.Clear(s, itemStyle, itemRect)

		item := blanc.NewText(" " + items[i] + " ")
		item.SetStyle(itemStyle)
		item.Draw(s, itemRect)
	}
}
