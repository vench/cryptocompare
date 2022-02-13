package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFormatMoney(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name      string
		input     float64
		out       string
		sym       string
		precision int
	}{
		{
			name:      "empty",
			input:     0,
			out:       "0.00",
			precision: 2,
		},
		{
			name:      "ok 1",
			input:     18955975,
			out:       "Ƀ 18,955,975.0",
			sym:       "Ƀ",
			precision: 1,
		},
		{
			name:      "ok 2",
			input:     803201245781.75,
			out:       "$ 803.20 B",
			sym:       "$",
			precision: 2,
		},
		{
			name:      "ok 3",
			input:     -1241.7699999999968,
			out:       "$ -1,241.77",
			sym:       "$",
			precision: 2,
		},
		{
			name:      " ok 4",
			input:     -2.8472016820402692,
			out:       "-2.85",
			precision: 2,
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			out := moneyFormat(tc.input, tc.sym, tc.precision)
			require.Equal(t, tc.out, out)
		})
	}
}
