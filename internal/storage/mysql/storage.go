package mysql

import (
	"database/sql"
	"fmt"

	"github.com/vench/cryptocompare/internal/entities"

	_ "github.com/go-sql-driver/mysql"
	"github.com/vench/cryptocompare/internal/config"
)

type Storage struct {
	conn *sql.DB
}

func New(conf *config.Mysql) (*Storage, error) {
	db, err := sql.Open("mysql", conf.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql connection: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql: %v", err)
	}

	return &Storage{
		conn: db,
	}, nil
}

func (s *Storage) StoreCurrency(currency *entities.Currency) error {
	return nil
}

func (s *Storage) GetCurrencyBy(fromSymbol, toSymbol string) (*entities.Currency, error) {
	return nil, nil
}

func (s *Storage) Close() error {
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("failed to close mysql connection: %v", err)
	}
	return nil
}
