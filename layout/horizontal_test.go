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
			input:  Rect{X: 0, Y: 0, W: 9, H: 2},
			params: []Constraint{Ratio(1, 3), Ratio(2, 3)},
			expected: []Rect{
				{X: 0, Y: 0, W: 3, H: 2},
				{X: 3, Y: 0, W: 6, H: 2},
			},
		},
	}

	for _, test := range cases {
		actual, err := SplitHorizontally(test.input, test.params...)
		require.NoError(t, err, test)
		require.EqualValues(t, test.expected, actual, test)
	}
}
