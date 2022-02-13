package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

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

type numberResponse float64

func (f numberResponse) Float() float64 {
	return float64(f)
}

func (f numberResponse) Int() int64 {
	return int64(f)
}

func (f *numberResponse) Unmarshal(b []byte) error {
	v, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		vi, err := strconv.ParseInt(string(b), 10, 64)
		if err != nil {
			return err
		}
		v = float64(vi)
	}
	*f = numberResponse(v)
	return nil
}

type currencyResponse struct {
	PRICE           numberResponse `json:"PRICE"`
	VOLUME24HOUR    numberResponse `json:"VOLUME24HOUR"`
	VOLUME24HOURTO  numberResponse `json:"VOLUME24HOURTO"`
	OPEN24HOUR      numberResponse `json:"OPEN24HOUR"`
	HIGH24HOUR      numberResponse `json:"HIGH24HOUR"`
	LOW24HOUR       numberResponse `json:"LOW24HOUR"`
	CHANGE24HOUR    numberResponse `json:"CHANGE24HOUR"`
	CHANGEPCT24HOUR numberResponse `json:"CHANGEPCT24HOUR"`
	SUPPLY          numberResponse `json:"SUPPLY"`
	MKTCAP          numberResponse `json:"MKTCAP"`
}

type response struct {
	Raw map[string]map[string]currencyResponse `json:"RAW"`
}

func (s *Storage) GetCurrencyBy(fromSymbol, toSymbol []string) ([]*entities.Currency, error) {
	url := fmt.Sprintf("%s?fsyms=%s&tsyms=%s",
		s.conf.URL,
		strings.Join(fromSymbol, ","),
		strings.Join(toSymbol, ","),
	)

	s.logger.Debug("url to cryptocompare", zap.String("url", url))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
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

	var data response
	if err := jsoniter.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal body: %w", err)
	}

	s.logger.Debug("data", zap.Reflect("data", data))

	result := make([]*entities.Currency, 0)
	for from, item := range data.Raw {
		for to, val := range item {
			result = append(result, &entities.Currency{
				FromSymbol: from,
				ToSymbol:   to,

				PRICE:           val.PRICE.Float(),
				VOLUME24HOUR:    val.VOLUME24HOUR.Float(),
				VOLUME24HOURTO:  val.VOLUME24HOURTO.Float(),
				OPEN24HOUR:      val.OPEN24HOUR.Float(),
				HIGH24HOUR:      val.HIGH24HOUR.Float(),
				LOW24HOUR:       val.LOW24HOUR.Float(),
				CHANGE24HOUR:    val.CHANGE24HOUR.Float(),
				CHANGEPCT24HOUR: val.CHANGEPCT24HOUR.Float(),
				MKTCAP:          val.MKTCAP.Float(),
				SUPPLY:          val.SUPPLY.Float(),
			})
		}
	}

	return result, nil
}
