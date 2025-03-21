package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(cfg struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}, logger *logrus.Logger) (*Storage, error) {
	connStr := "host=" + cfg.Host + " port=" + cfg.Port + " user=" + cfg.User +
		" password=" + cfg.Password + " dbname=" + cfg.DBName + " sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Инициализация таблицы exchange_rates
	queries := []string{
		`CREATE TABLE IF NOT EXISTS exchange_rates (
            id SERIAL PRIMARY KEY,
            from_currency VARCHAR(3) NOT NULL,
            to_currency VARCHAR(3) NOT NULL,
            rate DECIMAL(15,6) NOT NULL,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(from_currency, to_currency)
        )`,
		`INSERT INTO exchange_rates (from_currency, to_currency, rate)
         SELECT 'USD', 'EUR', 0.85 WHERE NOT EXISTS (
             SELECT 1 FROM exchange_rates WHERE from_currency = 'USD' AND to_currency = 'EUR'
         )`,
		`INSERT INTO exchange_rates (from_currency, to_currency, rate)
         SELECT 'EUR', 'USD', 1.18 WHERE NOT EXISTS (
             SELECT 1 FROM exchange_rates WHERE from_currency = 'EUR' AND to_currency = 'USD'
         )`,
		`INSERT INTO exchange_rates (from_currency, to_currency, rate)
         SELECT 'USD', 'RUB', 90.00 WHERE NOT EXISTS (
             SELECT 1 FROM exchange_rates WHERE from_currency = 'USD' AND to_currency = 'RUB'
         )`,
	}

	for _, query := range queries {
		_, err = db.Exec(query)
		if err != nil {
			logger.WithError(err).Error("failed to initialize exchange_rates table")
			return nil, err
		}
	}

	logger.Info("exchange_rates table initialized")
	return &Storage{DB: db}, nil
}
