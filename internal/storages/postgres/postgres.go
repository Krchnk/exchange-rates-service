package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
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
}) (*Storage, error) {
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

	return &Storage{DB: db}, nil
}
