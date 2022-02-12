package mysql

import (
	"testing"

	"github.com/vench/cryptocompare/internal/entities"

	"github.com/stretchr/testify/require"
	"github.com/vench/cryptocompare/internal/config"
)

func TestStorage_GetCurrencyBy(t *testing.T) {
	t.Parallel()

}

func TestStorage_StoreCurrency(t *testing.T) {
	t.Parallel()

	s, err := New(&config.Mysql{
		ConnectionString: "root:admin@tcp(127.0.0.1:3306)/test",
	})
	require.NoError(t, err)

	val := &entities.Currency{
		FromSymbol: "ABC",
		ToSymbol:   "XYZ",

		PRICE:           1000.10,
		CHANGEPCT24HOUR: 1999.2,
		CHANGE24HOUR:    4000.1,
	}

	require.NoError(t, s.StoreCurrency(val))

	cur, err := s.GetCurrencyBy([]string{val.FromSymbol}, []string{val.ToSymbol})
	require.NoError(t, err)
	require.Equal(t, 1, len(cur))
	require.Equal(t, val, cur[0])
}
