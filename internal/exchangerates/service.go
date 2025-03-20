package exchangerates

import (
	"context"
	"database/sql"

	"github.com/Krchnk/currency-wallet-proto/exchangerates"
	"github.com/Krchnk/exchange-rates-service/internal/storages/postgres"
	"github.com/sirupsen/logrus"
)

type Service struct {
	exchangerates.UnimplementedExchangeRatesServiceServer
	db     *sql.DB
	logger *logrus.Logger
}

func NewService(store *postgres.Storage, logger *logrus.Logger) *Service {
	return &Service{
		db:     store.DB,
		logger: logger,
	}
}

func (s *Service) GetExchangeRates(ctx context.Context, req *exchangerates.GetExchangeRatesRequest) (*exchangerates.GetExchangeRatesResponse, error) {
	s.logger.Info("Received GetExchangeRates request")

	rows, err := s.db.Query("SELECT from_currency, to_currency, rate FROM exchange_rates")
	if err != nil {
		s.logger.WithError(err).Error("failed to query exchange rates")
		return nil, err
	}
	defer rows.Close()

	var rates []*exchangerates.ExchangeRate
	for rows.Next() {
		var fromCurrency, toCurrency string
		var rate float64
		if err := rows.Scan(&fromCurrency, &toCurrency, &rate); err != nil {
			s.logger.WithError(err).Error("failed to scan exchange rate")
			return nil, err
		}
		rates = append(rates, &exchangerates.ExchangeRate{
			FromCurrency: fromCurrency,
			ToCurrency:   toCurrency,
			Rate:         rate,
		})
	}

	return &exchangerates.GetExchangeRatesResponse{Rates: rates}, nil
}
