package layout

type Rect struct {
	X int
	Y int
	W int
	H int
}

func (r Rect) Intersects(other Rect) bool {
	return r.Left() < other.Right() && r.Right() > other.Right() && r.Top() > r.Bottom() && r.Bottom() < r.Top()
}

func (r Rect) Top() int    { return r.Y }
func (r Rect) Bottom() int { return r.Y + r.H - 1 }
func (r Rect) Left() int   { return r.X }
func (r Rect) Right() int  { return r.X + r.W - 1 }

func (r Rect) Align(parent Rect, align AlignType) Rect {
	return Align(parent, r, align)
}

func (r Rect) ShiftTop(shift int) Rect    { r.Y += shift; return r }
func (r Rect) ShiftBottom(shift int) Rect { r.Y -= shift; return r }
func (r Rect) ShiftLeft(shift int) Rect   { r.X += shift; return r }
func (r Rect) ShiftRight(shift int) Rect  { r.X -= shift; return r }

func (r Rect) PadLeft(pad int) Rect {
	if r.W < pad {
		return r
	}
	r.X += pad
	r.W -= pad
	return r
}

func (r Rect) PadRight(pad int) Rect {
	if r.W < pad {
		return r
	}
	r.W -= pad
	return r
}

func (r Rect) PadHorizontal(pad int) Rect {
	if r.W < 2*pad {
		return r
	}
	r.X += pad
	r.W -= 2 * pad
	return r
}

func (r Rect) PadTop(pad int) Rect {
	if r.H < pad {
		return r
	}
	r.Y += pad
	r.H -= pad
	return r
}

func (r Rect) PadBottom(pad int) Rect {
	if r.H < pad {
		return r
	}
	r.H -= pad
	return r
}

func (r Rect) PadVertical(pad int) Rect {
	if r.H < 2*pad {
		return r
	}
	r.Y += pad
	r.H -= 2 * pad
	return r
}

func (r Rect) Pad(pad int) Rect {
	if r.W < 2*pad || r.H < 2*pad {
		return r
	}
	r.X += pad
	r.Y += pad
	r.W -= 2 * pad
	r.H -= 2 * pad
	return r
}
