package layout

type Rect struct {
	x int
	y int
	w int
	h int
}

func (r Rect) PadLeft(pad int) Rect {
	if r.w < pad {
		return r
	}
	r.x += pad
	r.w -= pad
	return r
}

func (r Rect) PadRight(pad int) Rect {
	if r.w < pad {
		return r
	}
	r.w -= pad
	return r
}

func (r Rect) PadHorizontal(pad int) Rect {
	if r.w < 2*pad {
		return r
	}
	r.x += pad
	r.w -= 2 * pad
	return r
}

func (r Rect) PadTop(pad int) Rect {
	if r.h < pad {
		return r
	}
	r.y += pad
	r.h -= pad
	return r
}

func (r Rect) PadBottom(pad int) Rect {
	if r.h < pad {
		return r
	}
	r.h -= pad
	return r
}

func (r Rect) PadVertical(pad int) Rect {
	if r.h < 2*pad {
		return r
	}
	r.y += pad
	r.h -= 2 * pad
	return r
}

func (r Rect) Pad(pad int) Rect {
	if r.w < 2*pad || r.h < 2*pad {
		return r
	}
	r.x += pad
	r.y += pad
	r.w -= 2 * pad
	r.h -= 2 * pad
	return r
}
