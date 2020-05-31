package main

import (
	"fmt"
	"github.com/lithdew/casso"
)

type Layout struct {
	solver *casso.Solver
	tags   []casso.Symbol
	err    error
}

func New(solver *casso.Solver) Layout {
	return Layout{solver: solver}
}

func (a *Layout) Required(constraints ...casso.Constraint) { a.Apply(casso.Required, constraints...) }
func (a *Layout) Strong(constraints ...casso.Constraint)   { a.Apply(casso.Strong, constraints...) }
func (a *Layout) Medium(constraints ...casso.Constraint)   { a.Apply(casso.Medium, constraints...) }
func (a *Layout) Weak(constraints ...casso.Constraint)     { a.Apply(casso.Weak, constraints...) }

func (a *Layout) Apply(priority casso.Priority, constraints ...casso.Constraint) {
	if a.err != nil {
		return
	}
	for _, constraint := range constraints {
		tag, err := a.solver.AddConstraintWithPriority(priority, constraint)
		if err != nil {
			a.err = err
			return
		}
		a.tags = append(a.tags, tag)
	}
}

func (a *Layout) Finalize() error {
	err := a.err
	if err != nil {
		a.Destroy()
	}
	return err
}

func (a *Layout) Destroy() {
	for _, tag := range a.tags {
		_ = a.solver.RemoveConstraint(tag)
	}
	a.tags = a.tags[:0]
}

type Box struct {
	x casso.Symbol
	y casso.Symbol
	w casso.Symbol
	h casso.Symbol

	solver *casso.Solver
	tags   []casso.Symbol
}

func NewBox(solver *casso.Solver) Box {
	box := Box{solver: solver, x: casso.New(), y: casso.New(), w: casso.New(), h: casso.New()}

	constraints := []casso.Constraint{
		Nonzero(box.x),
		Nonzero(box.y),
		Nonzero(box.w),
		Nonzero(box.h),
	}

	for _, constraint := range constraints {
		if _, err := solver.AddConstraint(constraint); err != nil {
			panic(fmt.Errorf("failed to add basic constraints to box: %w", err))
		}
	}

	return box
}

func (b Box) X() float64 { return b.solver.Val(b.x) }
func (b Box) Y() float64 { return b.solver.Val(b.y) }
func (b Box) W() float64 { return b.solver.Val(b.w) }
func (b Box) H() float64 { return b.solver.Val(b.h) }

func (b *Box) SetX(x float64) error { return b.solver.Suggest(b.x, x) }
func (b *Box) SetY(y float64) error { return b.solver.Suggest(b.y, y) }
func (b *Box) SetW(w float64) error { return b.solver.Suggest(b.w, w) }
func (b *Box) SetH(h float64) error { return b.solver.Suggest(b.h, h) }

func (b Box) Fixed(priority casso.Priority) {
	_ = b.solver.Edit(b.x, priority)
	_ = b.solver.Edit(b.y, priority)
	_ = b.solver.Edit(b.w, priority)
	_ = b.solver.Edit(b.h, priority)
}

// child.x >= parent.x + padding
// child.y >= parent.y + padding
// child.x + child.w == parent.x + parent.w - padding
// childy + child.h == parent.y + parent.h - padding
func Inside(parent, child Box, padding float64) []casso.Constraint {
	return []casso.Constraint{
		casso.NewConstraint(casso.GTE, -padding, child.x.T(1), parent.x.T(-1)),
		casso.NewConstraint(casso.GTE, -padding, child.y.T(1), parent.y.T(-1)),
		casso.NewConstraint(casso.LTE, padding, child.x.T(1), child.w.T(1), parent.x.T(-1), parent.w.T(-1)),
		casso.NewConstraint(casso.LTE, padding, child.y.T(1), child.h.T(1), parent.y.T(-1), parent.h.T(-1)),
	}
}

// symbol >= 0
func Nonzero(symbol casso.Symbol) casso.Constraint {
	return casso.NewConstraint(casso.GTE, 0, symbol.T(1))
}

// child.width == ratio * parent.width
func FillX(parent, child Box, ratio float64) casso.Constraint {
	return casso.NewConstraint(casso.EQ, 0, child.w.T(1), parent.w.T(-ratio))
}

// child.width == ratio * parent.width
func FillY(parent, child Box, ratio float64) casso.Constraint {
	return casso.NewConstraint(casso.EQ, 0, child.h.T(1), parent.h.T(-ratio))
}

// a.x + a.width + spacing == b.x
func SpaceBetween(a, b Box, spacing float64) casso.Constraint {
	return casso.NewConstraint(casso.EQ, spacing, a.x.T(1), a.w.T(1), b.x.T(-1))
}

// a.width == b.width
func SameWidth(a, b Box) casso.Constraint {
	return casso.NewConstraint(casso.EQ, 0, a.w.T(1), b.w.T(-1))
}
