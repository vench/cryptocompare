package api

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vench/cryptocompare/internal/config"
	"go.uber.org/zap"
)

func TestStorage_GetCurrencyBy(t *testing.T) {
	t.Parallel()

	s := Storage{
		logger: zap.NewNop(),
		conf: &config.CryptoCompare{
			Url: "https://min-api.cryptocompare.com/data/pricemultifull",
		},
	}

	result, err := s.GetCurrencyBy([]string{"BTC"}, []string{"USD"})
	require.NoError(t, err)
	_ = result
}
