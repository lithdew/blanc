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

func (a AlignType) Is(align AlignType) bool { return a&align != 0 }

func (a AlignType) Valid() bool { return !(a.Is(Top) && a.Is(Bottom) || a.Is(Left) && a.Is(Right)) }

// Align aligns a child _against_ the parent (outside the parent).
func Align(parent, child Rect, a AlignType) Rect {
	switch {
	case a.Is(Left):
		child.X = parent.Left() - child.W
	case a.Is(Right):
		child.X = parent.Right() + 1
	case a.Is(Center):
		child.X = parent.CenterX()
	}

	switch {
	case a.Is(Top):
		child.Y = parent.Top() - 1
	case a.Is(Bottom):
		child.Y = parent.Bottom() + 1
	case a.Is(Center):
		child.Y = parent.CenterY()
	}
	return child
}

// Position positions a child _inside_ the parent (inside the parent).
func Position(parent, child Rect, a AlignType) Rect {
	switch {
	case a.Is(Left):
		child.X = parent.Left()
	case a.Is(Right):
		child.X = parent.Right() - child.W
	case a.Is(Center):
		child.X = parent.CenterX() - child.W/2
	}
	switch {
	case a.Is(Top):
		child.Y = parent.Top()
	case a.Is(Bottom):
		child.Y = parent.Bottom() - child.H
	case a.Is(Center):
		child.Y = parent.CenterY() - child.H/2
	}
	return child
}

func Text(text string) Rect { return Rect{W: runewidth.StringWidth(text), H: 1} }
