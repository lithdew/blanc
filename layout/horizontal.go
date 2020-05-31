package layout

import "github.com/lithdew/casso"

func SplitHorizontally(r Rect, constraints ...Constraint) ([]Rect, error) {
	type Element struct {
		x casso.Symbol
		w casso.Symbol
	}

	results := make([]Rect, len(constraints))
	elements := make([]Element, len(constraints))

	solver := casso.NewSolver()
	layout := New(solver)

	for i := 0; i < len(constraints); i++ {
		// w >= 0
		// x >= r.x
		// x + w >= r.x + r.w

		// first.x + first.w == second.x

		// first.x == r.x
		// last.x + last.w == r.x + r.w

		// apply constraint

		elements[i].x = casso.New()
		elements[i].w = casso.New()

		layout.Required(GTE(0, elements[i].w.T(1)))
		layout.Required(GTE(-float64(r.x), elements[i].x.T(1)))
		layout.Required(LTE(-float64(r.x+r.w), elements[i].x.T(1), elements[i].w.T(1)))

		if i > 0 {
			layout.Required(EQ(0, elements[i-1].x.T(1), elements[i-1].w.T(1), elements[i].x.T(-1)))
		}

		if i == 0 {
			layout.Required(EQ(-float64(r.x), elements[0].x.T(1)))
		}

		if i == len(constraints)-1 {
			layout.Required(EQ(-float64(r.x+r.w), elements[len(elements)-1].x.T(1), elements[len(elements)-1].w.T(1)))
		}

		layout.Weak(constraints[i](r.w, elements[i].w))
	}

	if err := layout.Finalize(); err != nil {
		return nil, err
	}

	for i := 0; i < len(elements); i++ {
		results[i].x = int(solver.Val(elements[i].x))
		results[i].y = r.y
		results[i].w = int(solver.Val(elements[i].w))
		results[i].h = r.h
	}

	return results, nil
}
