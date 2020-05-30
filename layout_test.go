package main

import (
	"fmt"
	"github.com/lithdew/casso"
	"github.com/stretchr/testify/require"
	"testing"
)

type Box struct {
	x casso.Symbol
	y casso.Symbol
	w casso.Symbol
	h casso.Symbol

	solver *casso.Solver
	tags   []casso.Symbol
}

func NewBox(solver *casso.Solver) Box {
	return Box{solver: solver, x: casso.New(), y: casso.New(), w: casso.New(), h: casso.New()}
}

func (b Box) X() float64 { return b.solver.Val(b.x) }
func (b Box) Y() float64 { return b.solver.Val(b.y) }
func (b Box) W() float64 { return b.solver.Val(b.w) }
func (b Box) H() float64 { return b.solver.Val(b.h) }

func (b *Box) SetX(x float64) error { return b.solver.Suggest(b.x, x) }
func (b *Box) SetY(y float64) error { return b.solver.Suggest(b.y, y) }
func (b *Box) SetW(w float64) error { return b.solver.Suggest(b.w, w) }
func (b *Box) SetH(h float64) error { return b.solver.Suggest(b.h, h) }

func (b Box) Editable() {
	_ = b.solver.Edit(b.x, casso.Medium)
	_ = b.solver.Edit(b.y, casso.Medium)
	_ = b.solver.Edit(b.w, casso.Medium)
	_ = b.solver.Edit(b.h, casso.Medium)
}

func Inside(parent, child Box, padding float64) []casso.Constraint {
	// child.x >= parent.x + padding
	// child.y >= parent.y + padding
	// child.width <= parent.width - 2 * padding
	// child.height <= parent.height - 2 * padding
	return []casso.Constraint{
		casso.NewConstraint(casso.GTE, -padding, child.x.T(1), parent.x.T(-1)),
		casso.NewConstraint(casso.GTE, -padding, child.y.T(1), parent.y.T(-1)),
		casso.NewConstraint(casso.LTE, 2*padding, child.w.T(1), parent.w.T(-1)),
		casso.NewConstraint(casso.LTE, 2*padding, child.h.T(1), parent.h.T(-1)),
	}
}

func Fill(parent, child Box) []casso.Constraint {
	// child.width == parent.width
	// child.height == parent.height
	return []casso.Constraint{
		casso.NewConstraint(casso.EQ, 0, child.w.T(1), parent.w.T(-1)),
		casso.NewConstraint(casso.EQ, 0, child.h.T(1), parent.h.T(-1)),
	}
}

func Apply(solver *casso.Solver, priority casso.Priority, constraints ...casso.Constraint) ([]casso.Symbol, error) {
	tags := make([]casso.Symbol, 0, len(constraints))
	for _, constraint := range constraints {
		tag, err := solver.AddConstraintWithPriority(priority, constraint)
		if err != nil {
			for _, tag := range tags {
				_ = solver.RemoveConstraint(tag)
			}
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func TestBoxInside(t *testing.T) {
	solver := casso.NewSolver()

	parent := NewBox(solver)
	child := NewBox(solver)

	parent.Editable()
	//child.Editable()

	require.NoError(t, parent.SetX(0))
	require.NoError(t, parent.SetY(0))
	require.NoError(t, parent.SetW(800))
	require.NoError(t, parent.SetH(600))

	_, err := Apply(solver, casso.Required, Inside(parent, child, 10)...)
	require.NoError(t, err)

	_, err = Apply(solver, casso.Medium, Fill(parent, child)...)
	require.NoError(t, err)

	fmt.Println(parent.X(), parent.Y(), parent.W(), parent.H())
	fmt.Println(child.X(), child.Y(), child.W(), child.H())
}
