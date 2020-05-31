package layout

import "github.com/lithdew/casso"

type Constraint func(a int, b casso.Symbol) casso.Constraint

func Percentage(percent int) Constraint {
	if percent < 0 || percent > 100 {
		panic("BUG: Percentage(percent int) may only be called with number in range [0, 100].")
	}
	return func(a int, b casso.Symbol) casso.Constraint {
		return EQ(-float64(a*percent/100), b.T(1))
	}
}

func Ratio(num, den int) Constraint {
	return func(a int, b casso.Symbol) casso.Constraint {
		return EQ(-float64(a*num/den), b.T(1))
	}
}

func Min(min int) Constraint {
	return func(_ int, b casso.Symbol) casso.Constraint {
		return GTE(-float64(min), b.T(1))
	}
}

func Max(max int) Constraint {
	return func(_ int, b casso.Symbol) casso.Constraint {
		return LTE(-float64(max), b.T(1))
	}
}

func Length(val int) Constraint {
	return func(_ int, b casso.Symbol) casso.Constraint {
		return EQ(-float64(val), b.T(1))
	}
}

func EQ(constant float64, terms ...casso.Term) casso.Constraint {
	return casso.NewConstraint(casso.EQ, constant, terms...)
}

func GTE(constant float64, terms ...casso.Term) casso.Constraint {
	return casso.NewConstraint(casso.GTE, constant, terms...)
}

func LTE(constant float64, terms ...casso.Term) casso.Constraint {
	return casso.NewConstraint(casso.LTE, constant, terms...)
}
