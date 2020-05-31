package layout

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSplitHorizontally(t *testing.T) {
	cases := []struct {
		input    Rect
		params   []Constraint
		expected []Rect
	}{
		{
			input:  Rect{x: 0, y: 0, w: 9, h: 2},
			params: []Constraint{Ratio(1, 3), Ratio(2, 3)},
			expected: []Rect{
				{x: 0, y: 0, w: 3, h: 2},
				{x: 3, y: 0, w: 6, h: 2},
			},
		},
	}

	for _, test := range cases {
		actual, err := SplitHorizontally(test.input, test.params...)
		require.NoError(t, err, test)
		require.EqualValues(t, test.expected, actual, test)
	}
}
