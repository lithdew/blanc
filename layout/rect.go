package layout

type Rect struct {
	X int
	Y int
	W int
	H int
}

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
