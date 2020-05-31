package layout

import (
	"github.com/mattn/go-runewidth"
)

type AlignType uint8

const (
	Center AlignType = 1 << iota
	Top
	Bottom
	Left
	Right
)

func (a AlignType) Is(align AlignType) bool {
	return a&align != 0
}

func (a AlignType) Valid() bool {
	return !(a.Is(Top) && a.Is(Bottom) || a.Is(Left) && a.Is(Right))
}

func Align(parent, child Rect, a AlignType) Rect {
	switch {
	case a.Is(Left):
		child.X = parent.X
	case a.Is(Right):
		child.X = parent.X + (parent.W - 1) - child.W
	default:
		child.X = parent.X + (parent.W-1)/2 - child.W/2
	}
	switch {
	case a.Is(Top):
		child.Y = parent.Y + child.H/2
	case a.Is(Bottom):
		child.Y = parent.Y + (parent.H - 1) - child.H/2
	default:
		child.Y = parent.Y + (parent.H-1)/2
	}
	return child
}

func TextBounds(text string) Rect {
	return Rect{X: 0, Y: 0, W: runewidth.StringWidth(text), H: 1}
}
