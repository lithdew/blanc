package layout

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitVertically(t *testing.T) {
	cases := []struct {
		input    Rect
		params   []Constraint
		expected []Rect
	}{
		{
			input:  Rect{X: 2, Y: 2, W: 10, H: 10},
			params: []Constraint{Length(5), Min(0)},
			expected: []Rect{
				{X: 2, Y: 2, W: 10, H: 5},
				{X: 2, Y: 7, W: 10, H: 5},
			},
		},
		{
			input:  Rect{X: 2, Y: 2, W: 10, H: 10},
			params: []Constraint{Ratio(1, 3), Ratio(2, 3)},
			expected: []Rect{
				{X: 2, Y: 2, W: 10, H: 3},
				{X: 2, Y: 5, W: 10, H: 7},
			},
		},
		{
			input:  Rect{X: 2, Y: 2, W: 10, H: 10},
			params: []Constraint{Percentage(40), Length(1), Min(0)},
			expected: []Rect{
				{X: 2, Y: 2, W: 10, H: 4},
				{X: 2, Y: 6, W: 10, H: 1},
				{X: 2, Y: 7, W: 10, H: 5},
			},
		},
		{
			input:  Rect{X: 2, Y: 2, W: 10, H: 10},
			params: []Constraint{Percentage(10), Max(5), Min(1)},
			expected: []Rect{
				{X: 2, Y: 2, W: 10, H: 1},
				{X: 2, Y: 3, W: 10, H: 5},
				{X: 2, Y: 8, W: 10, H: 4},
			},
		},
	}

	for _, test := range cases {
		actual, err := SplitVertically(test.input, test.params...)
		require.NoError(t, err, test)
		require.EqualValues(t, test.expected, actual, test)
	}
}
