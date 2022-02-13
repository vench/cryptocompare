package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vench/cryptocompare/internal/entities"

	"github.com/vench/cryptocompare/internal/config"
	"go.uber.org/zap"
)

//nolint:lll
const responseJSON = "{\n\"RAW\": {\n\"BTC\": {\n\"USD\": {\n\"TYPE\": \"5\",\n\"MARKET\": \"CCCAGG\",\n\"FROMSYMBOL\": \"BTC\",\n\"TOSYMBOL\": \"USD\",\n\"FLAGS\": \"2052\",\n\"PRICE\": 42371.93,\n\"LASTUPDATE\": 1644683102,\n\"MEDIAN\": 42377.53,\n\"LASTVOLUME\": 0.02071615,\n\"LASTVOLUMETO\": 877.6382446195,\n\"LASTTRADEID\": \"280324221\",\n\"VOLUMEDAY\": 10561.68326764645,\n\"VOLUMEDAYTO\": 445839153.842952,\n\"VOLUME24HOUR\": 26367.795692919997,\n\"VOLUME24HOURTO\": 1120081383.2888746,\n\"OPENDAY\": 42399.36,\n\"HIGHDAY\": 42495.54,\n\"LOWDAY\": 41786.2,\n\"OPEN24HOUR\": 43613.7,\n\"HIGH24HOUR\": 43755.05,\n\"LOW24HOUR\": 41767.34,\n\"LASTMARKET\": \"Coinbase\",\n\"VOLUMEHOUR\": 723.350194739971,\n\"VOLUMEHOURTO\": 30602731.764161117,\n\"OPENHOUR\": 42139.18,\n\"HIGHHOUR\": 42457.15,\n\"LOWHOUR\": 42022.17,\n\"TOPTIERVOLUME24HOUR\": 26367.690653919995,\n\"TOPTIERVOLUME24HOURTO\": 1120076771.2158,\n\"CHANGE24HOUR\": -1241.7699999999968,\n\"CHANGEPCT24HOUR\": -2.8472016820402692,\n\"CHANGEDAY\": -27.43000000000029,\n\"CHANGEPCTDAY\": -0.06469437274524967,\n\"CHANGEHOUR\": 232.75,\n\"CHANGEPCTHOUR\": 0.5523363292783581,\n\"CONVERSIONTYPE\": \"direct\",\n\"CONVERSIONSYMBOL\": \"\",\n\"SUPPLY\": 18955975,\n\"MKTCAP\": 803201245781.75,\n\"MKTCAPPENALTY\": 0,\n\"CIRCULATINGSUPPLY\": 18955975,\n\"CIRCULATINGSUPPLYMKTCAP\": 803201245781.75,\n\"TOTALVOLUME24H\": 143859.43952671083,\n\"TOTALVOLUME24HTO\": 6098429091.399192,\n\"TOTALTOPTIERVOLUME24H\": 143614.4319693793,\n\"TOTALTOPTIERVOLUME24HTO\": 6088047486.96255,\n\"IMAGEURL\": \"/media/37746251/btc.png\"\n}\n}\n},\n\"DISPLAY\": {\n\"BTC\": {\n\"USD\": {\n\"FROMSYMBOL\": \"Ƀ\",\n\"TOSYMBOL\": \"$\",\n\"MARKET\": \"CryptoCompare Index\",\n\"PRICE\": \"$ 42,371.9\",\n\"LASTUPDATE\": \"Just now\",\n\"LASTVOLUME\": \"Ƀ 0.02072\",\n\"LASTVOLUMETO\": \"$ 877.64\",\n\"LASTTRADEID\": \"280324221\",\n\"VOLUMEDAY\": \"Ƀ 10,561.7\",\n\"VOLUMEDAYTO\": \"$ 445,839,153.8\",\n\"VOLUME24HOUR\": \"Ƀ 26,367.8\",\n\"VOLUME24HOURTO\": \"$ 1,120,081,383.3\",\n\"OPENDAY\": \"$ 42,399.4\",\n\"HIGHDAY\": \"$ 42,495.5\",\n\"LOWDAY\": \"$ 41,786.2\",\n\"OPEN24HOUR\": \"$ 43,613.7\",\n\"HIGH24HOUR\": \"$ 43,755.1\",\n\"LOW24HOUR\": \"$ 41,767.3\",\n\"LASTMARKET\": \"Coinbase\",\n\"VOLUMEHOUR\": \"Ƀ 723.35\",\n\"VOLUMEHOURTO\": \"$ 30,602,731.8\",\n\"OPENHOUR\": \"$ 42,139.2\",\n\"HIGHHOUR\": \"$ 42,457.2\",\n\"LOWHOUR\": \"$ 42,022.2\",\n\"TOPTIERVOLUME24HOUR\": \"Ƀ 26,367.7\",\n\"TOPTIERVOLUME24HOURTO\": \"$ 1,120,076,771.2\",\n\"CHANGE24HOUR\": \"$ -1,241.77\",\n\"CHANGEPCT24HOUR\": \"-2.85\",\n\"CHANGEDAY\": \"$ -27.43\",\n\"CHANGEPCTDAY\": \"-0.06\",\n\"CHANGEHOUR\": \"$ 232.75\",\n\"CHANGEPCTHOUR\": \"0.55\",\n\"CONVERSIONTYPE\": \"direct\",\n\"CONVERSIONSYMBOL\": \"\",\n\"SUPPLY\": \"Ƀ 18,955,975.0\",\n\"MKTCAP\": \"$ 803.20 B\",\n\"MKTCAPPENALTY\": \"0 %\",\n\"CIRCULATINGSUPPLY\": \"Ƀ 18,955,975.0\",\n\"CIRCULATINGSUPPLYMKTCAP\": \"$ 803.20 B\",\n\"TOTALVOLUME24H\": \"Ƀ 143.86 K\",\n\"TOTALVOLUME24HTO\": \"$ 6.10 B\",\n\"TOTALTOPTIERVOLUME24H\": \"Ƀ 143.61 K\",\n\"TOTALTOPTIERVOLUME24HTO\": \"$ 6.09 B\",\n\"IMAGEURL\": \"/media/37746251/btc.png\"\n}\n}\n}\n}"

func TestStorage_GetCurrencyBy(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		checkResult func(*testing.T, []*entities.Currency, error)
		handler     func(writer http.ResponseWriter, request *http.Request)
	}{
		{
			name: "ok empty",
			handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte("{}")) //nolint:errcheck
			},
			checkResult: func(t *testing.T, list []*entities.Currency, err error) {
				require.NoError(t, err)
				require.Equal(t, 0, len(list))
			},
		},
		{
			name: "error service",
			handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusInternalServerError)
			},
			checkResult: func(t *testing.T, list []*entities.Currency, err error) {
				require.Error(t, err)
			},
		},
		{
			name: "ok no empty",
			handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(responseJSON)) //nolint:errcheck
			},
			checkResult: func(t *testing.T, list []*entities.Currency, err error) {
				require.NoError(t, err)
				require.Equal(t, 1, len(list))
				require.Equal(t, &entities.Currency{
					FromSymbol:      "BTC",
					ToSymbol:        "USD",
					CHANGE24HOUR:    -1241.7699999999968,
					CHANGEPCT24HOUR: -2.8472016820402692,
					OPEN24HOUR:      43613.7,
					VOLUME24HOUR:    26367.795692919997,
					VOLUME24HOURTO:  1.1200813832888746e+09,
					LOW24HOUR:       41767.34,
					HIGH24HOUR:      43755.05,
					PRICE:           42371.93,
					MKTCAP:          8.0320124578175e+11,
					SUPPLY:          18955975,
				}, list[0])
			},
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.handler == nil {
				tc.handler = func(writer http.ResponseWriter, request *http.Request) {
					writer.WriteHeader(http.StatusOK)
				}
			}

			serv := httptest.NewServer(http.HandlerFunc(tc.handler))
			defer serv.Close()

			s := Storage{
				logger: zap.NewNop(),
				conf: &config.CryptoCompare{
					URL: serv.URL,
				},
			}

			result, err := s.GetCurrencyBy([]string{"BTC"}, []string{"USD"})
			if tc.checkResult != nil {
				tc.checkResult(t, result, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
