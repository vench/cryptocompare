package scheduler

import "fmt"

func (s *Scheduler) parseApiCryptocompare() error {
	s.logger.Debug("begin parse api cryptocompare")

	result, err := s.currencyReader.GetCurrencyBy(
		s.conf.CryptoCompare.FromSymbols,
		s.conf.CryptoCompare.ToSymbols,
	)
	if err != nil {
		return fmt.Errorf("failed to get currency: %w", err)
	}

	if err := s.currencyWriter.StoreCurrency(result...); err != nil {
		return fmt.Errorf("failed to store currency: %w", err)
	}

	return nil
}
