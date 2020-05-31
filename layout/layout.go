package layout

import "github.com/lithdew/casso"

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
