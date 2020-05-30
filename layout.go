package main

import (
	"github.com/gdamore/tcell"
	"github.com/lithdew/casso"
)

type Layout struct {
	app *App

	// FIXME(kenta): add origin x
	// FIXME(kenta): add origin y

	x casso.Symbol // x
	y casso.Symbol // y
	w casso.Symbol // width
	h casso.Symbol // height

	pt casso.Symbol // padding top
	pb casso.Symbol // padding bottom
	pl casso.Symbol // padding left
	pr casso.Symbol // padding right

	tags []casso.Symbol

	xo bool
	yo bool
	wo bool
	ho bool
}

func NewLayout() *Layout {
	return &Layout{
		x: casso.New(),
		y: casso.New(),
		w: casso.New(),
		h: casso.New(),

		pt: casso.New(),
		pb: casso.New(),
		pl: casso.New(),
		pr: casso.New(),
	}
}

func (v *Layout) Render(screen tcell.Screen) {
	box(
		screen,
		int(v.app.solver.Val(v.x)),
		int(v.app.solver.Val(v.y)),
		int(v.app.solver.Val(v.x)+v.app.solver.Val(v.w)),
		int(v.app.solver.Val(v.y)+v.app.solver.Val(v.h)),
		tcell.StyleDefault,
		' ',
	)
}

func (v *Layout) constrain(priority casso.Priority, cell casso.Constraint) error {
	tag, err := v.app.solver.AddConstraintWithPriority(priority, cell)
	if err != nil {
		return err
	}
	v.tags = append(v.tags, tag)
	return nil
}

func (v *Layout) SetX(x int) error {
	if !v.xo {
		if err := v.app.solver.Edit(v.x, casso.Strong); err != nil {
			return err
		}
		v.xo = true
	}
	return v.app.solver.Suggest(v.x, float64(x))
}

func (v *Layout) SetY(y int) error {
	if !v.yo {
		if err := v.app.solver.Edit(v.y, casso.Strong); err != nil {
			return err
		}
		v.yo = true
	}
	return v.app.solver.Suggest(v.y, float64(y))
}

func (v *Layout) SetWidth(w int) error {
	if !v.wo {
		if err := v.app.solver.Edit(v.w, casso.Strong); err != nil {
			return err
		}
		v.wo = true
	}
	return v.app.solver.Suggest(v.w, float64(w))
}

func (v *Layout) SetHeight(h int) error {
	if !v.ho {
		if err := v.app.solver.Edit(v.h, casso.Strong); err != nil {
			return err
		}
		v.ho = true
	}
	return v.app.solver.Suggest(v.h, float64(h))
}

// w == screen-width
func (v *Layout) FillX() error {
	return v.constrain(casso.Medium, casso.NewConstraint(casso.EQ, 1, v.w.T(1), v.app.sw.T(-1)))
}

// h == screen_height
func (v *Layout) FillY() error {
	return v.constrain(casso.Medium, casso.NewConstraint(casso.EQ, 1, v.h.T(1), v.app.sh.T(-1)))
}

func (v *Layout) SetPaddingTop(h int) error    { return v.app.solver.Suggest(v.pt, float64(h)) }
func (v *Layout) SetPaddingBottom(h int) error { return v.app.solver.Suggest(v.pb, float64(h)) }
func (v *Layout) SetPaddingLeft(h int) error   { return v.app.solver.Suggest(v.pl, float64(h)) }
func (v *Layout) SetPaddingRight(h int) error  { return v.app.solver.Suggest(v.pr, float64(h)) }

func (v *Layout) Init(app *App) error {
	v.app = app

	if err := app.solver.Edit(v.pt, casso.Strong); err != nil {
		return err
	}
	if err := app.solver.Edit(v.pb, casso.Strong); err != nil {
		return err
	}
	if err := app.solver.Edit(v.pl, casso.Strong); err != nil {
		return err
	}
	if err := app.solver.Edit(v.pr, casso.Strong); err != nil {
		return err
	}

	// w >= 0
	// h >= 0
	// x >= padding_left (FIXME)
	// y >= padding_top (FIXME)
	// x + w <= screen_width - 1 - padding_right
	// y + h <= screen_height - 1 - padding_bottom

	constraints := []casso.Constraint{
		casso.NewConstraint(casso.GTE, 0, v.w.T(1)),
		casso.NewConstraint(casso.GTE, 0, v.h.T(1)),
		casso.NewConstraint(casso.GTE, 0, v.x.T(1), v.pl.T(-1)),
		casso.NewConstraint(casso.GTE, 0, v.y.T(1), v.pt.T(-1)),
		casso.NewConstraint(casso.LTE, 1, v.x.T(1), v.w.T(1), v.pr.T(1), app.sw.T(-1)),
		casso.NewConstraint(casso.LTE, 1, v.y.T(1), v.h.T(1), v.pb.T(1), app.sh.T(-1)),
	}

	for _, constraint := range constraints {
		if err := v.constrain(casso.Required, constraint); err != nil {
			return err
		}
	}

	return nil
}
