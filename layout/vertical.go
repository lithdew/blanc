package layout

import "github.com/lithdew/casso"

func SplitVertically(r Rect, constraints ...Constraint) ([]Rect, error) {
	type Element struct {
		y casso.Symbol
		h casso.Symbol
	}

	results := make([]Rect, len(constraints))
	elements := make([]Element, len(constraints))

	solver := casso.NewSolver()
	layout := New(solver)

	for i := 0; i < len(constraints); i++ {
		// h >= 0
		// y >= r.y
		// y + h >= r.y + r.h

		// first.y + first.h == second.y

		// first.y == r.y
		// last.y + last.h == r.y + r.h

		// apply constraint

		elements[i].y = casso.New()
		elements[i].h = casso.New()

		layout.Required(GTE(0, elements[i].h.T(1)))
		layout.Required(GTE(-float64(r.y), elements[i].y.T(1)))
		layout.Required(LTE(-float64(r.y+r.h), elements[i].y.T(1), elements[i].h.T(1)))

		if i > 0 {
			layout.Required(EQ(0, elements[i-1].y.T(1), elements[i-1].h.T(1), elements[i].y.T(-1)))
		}

		if i == 0 {
			layout.Required(EQ(-float64(r.y), elements[0].y.T(1)))
		}

		if i == len(constraints)-1 {
			layout.Required(EQ(-float64(r.y+r.h), elements[len(elements)-1].y.T(1), elements[len(elements)-1].h.T(1)))
		}

		layout.Weak(constraints[i](r.h, elements[i].h))
	}

	if err := layout.Finalize(); err != nil {
		return nil, err
	}

	for i := 0; i < len(elements); i++ {
		results[i].x = r.x
		results[i].y = int(solver.Val(elements[i].y))
		results[i].w = r.w
		results[i].h = int(solver.Val(elements[i].h))
	}

	return results, nil
}
