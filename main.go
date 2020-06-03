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
	screen, err := tcell.NewScreen()
	check(err)

	check(screen.Init())
	defer screen.Fini()

	inputStyle := tcell.StyleDefault.Reverse(true)
	selectedStyle := tcell.StyleDefault.Background(tcell.ColorDarkGray).Foreground(tcell.ColorWhite)

	input := newTextbox()
	input.SetLabel(">>> ")
	input.SetTextStyle(inputStyle)
	input.SetLabelStyle(inputStyle)
	input.SetSelectedStyle(selectedStyle)

	handlers = append(handlers, input)

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

		headerRect := screenRect.Align(layout.Top | layout.Left).WidthOf(screenRect).Height(1)
		Clear(screen, titleStyle, headerRect)
		title.Draw(screen, layout.Text(title.Text()).AlignTo(headerRect, layout.Left).ShiftRight(1))

		// body

		bodyRect := screenRect.PadVertical(1)
		Clear(screen, contentStyle, bodyRect)

		content.Draw(screen, bodyRect.Pad(4))
		//graph.Draw(screen, bodyRect.Pad(4))

		// footer

		renderFooter(screen, screenRect, input)

		screen.Show()
	}
}

func renderFooter(s tcell.Screen, screenRect layout.Rect, input *Textbox) {
	footerRect := screenRect.Align(layout.Bottom | layout.Left).WidthOf(screenRect).Height(1)
	Clear(s, tcell.StyleDefault.Reverse(true), footerRect)

	inputRect := footerRect.PadLeft(1)
	input.Draw(s, inputRect)

	renderMenu(s, input, inputRect)
}

func renderMenu(s tcell.Screen, input *Textbox, inputRect layout.Rect) {
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

	Clear(s, style(-1), menuRect)

	for i := range items {
		itemRect := menuRect.Align(layout.Top | layout.Left).ShiftTop(i).WidthOf(menuRect).Height(1)
		itemStyle := style(i)

		Clear(s, itemStyle, itemRect)

		item := NewText(" " + items[i] + " ")
		item.SetStyle(itemStyle)
		item.Draw(s, itemRect)
	}
}
