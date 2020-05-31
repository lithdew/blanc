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
			input:  Rect{x: 2, y: 2, w: 10, h: 10},
			params: []Constraint{Length(5), Min(0)},
			expected: []Rect{
				{x: 2, y: 2, w: 10, h: 5},
				{x: 2, y: 7, w: 10, h: 5},
			},
		},
		{
			input:  Rect{x: 2, y: 2, w: 10, h: 10},
			params: []Constraint{Ratio(1, 3), Ratio(2, 3)},
			expected: []Rect{
				{x: 2, y: 2, w: 10, h: 3},
				{x: 2, y: 5, w: 10, h: 7},
			},
		},
		{
			input:  Rect{x: 2, y: 2, w: 10, h: 10},
			params: []Constraint{Percentage(40), Length(1), Min(0)},
			expected: []Rect{
				{x: 2, y: 2, w: 10, h: 4},
				{x: 2, y: 6, w: 10, h: 1},
				{x: 2, y: 7, w: 10, h: 5},
			},
		},
		{
			input:  Rect{x: 2, y: 2, w: 10, h: 10},
			params: []Constraint{Percentage(10), Max(5), Min(1)},
			expected: []Rect{
				{x: 2, y: 2, w: 10, h: 1},
				{x: 2, y: 3, w: 10, h: 5},
				{x: 2, y: 8, w: 10, h: 4},
			},
		},
	}

	for _, test := range cases {
		actual, err := SplitVertically(test.input, test.params...)
		require.NoError(t, err, test)
		require.EqualValues(t, test.expected, actual, test)
	}
}
