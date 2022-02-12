package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/vench/cryptocompare/internal/config"
	"github.com/vench/cryptocompare/internal/entities"
	"go.uber.org/zap"
)

type Storage struct {
	logger *zap.Logger
	conf   *config.CryptoCompare
}

func New(logger *zap.Logger, conf *config.CryptoCompare) (*Storage, error) {
	return &Storage{
		logger: logger,
		conf:   conf,
	}, nil
}

type currencyResponse struct {
	PRICE           string `json:"PRICE"`
	VOLUME24HOUR    string `json:"VOLUME24HOUR"`
	VOLUME24HOURTO  string `json:"VOLUME24HOURTO"`
	OPEN24HOUR      string `json:"OPEN24HOUR"`
	HIGH24HOUR      string `json:"HIGH24HOUR"`
	LOW24HOUR       string `json:"LOW24HOUR"`
	CHANGE24HOUR    string `json:"CHANGE24HOUR"`
	CHANGEPCT24HOUR string `json:"CHANGEPCT24HOUR"`
	SUPPLY          string `json:"SUPPLY"`
	MKTCAP          string `json:"MKTCAP"`
}

type response struct {
	Display map[string]map[string]currencyResponse `json:"DISPLAY"`
}

func (s *Storage) GetCurrencyBy(fromSymbol, toSymbol []string) ([]*entities.Currency, error) {
	url := fmt.Sprintf("%s?fsyms=%s&tsyms=%s",
		s.conf.Url,
		strings.Join(fromSymbol, ","),
		strings.Join(toSymbol, ","),
	)

	s.logger.Debug("url to cryptocompare", zap.String("url", url))

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to new request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	// unmarshal
	var data response
	if err := jsoniter.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	s.logger.Debug("data", zap.Reflect("data", data))

	result := make([]*entities.Currency, 0)
	for from, item := range data.Display {
		for to, val := range item {
			result = append(result, &entities.Currency{
				FromSymbol: from,
				ToSymbol:   to,

				PRICE:           val.PRICE,
				VOLUME24HOUR:    val.VOLUME24HOUR,
				VOLUME24HOURTO:  val.VOLUME24HOURTO,
				OPEN24HOUR:      val.OPEN24HOUR,
				HIGH24HOUR:      val.HIGH24HOUR,
				LOW24HOUR:       val.LOW24HOUR,
				CHANGE24HOUR:    val.CHANGE24HOUR,
				CHANGEPCT24HOUR: val.CHANGEPCT24HOUR,
				SUPPLY:          val.SUPPLY,
				MKTCAP:          val.MKTCAP,
			})

		}
	}

	return result, nil
}
