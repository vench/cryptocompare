package mysql

import (
	"database/sql"
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"

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

func (s *Storage) StoreCurrency(currencies ...*entities.Currency) error {
	for _, currency := range currencies {
		sql := "INSERT INTO currency(`key`,`value`,`updated_at`) VALUES (?,?,now()) " +
			"ON DUPLICATE KEY UPDATE `value` = VALUES(`value`), `updated_at`= VALUES(`updated_at`)"

		data, err := jsoniter.Marshal(currency)
		if err != nil {
			return fmt.Errorf("failed to marshal currency: %w", err)
		}

		if _, err := s.conn.Exec(sql, key(currency), data); err != nil {
			return fmt.Errorf("failed to insert row: %w", err)
		}
	}
	return nil
}

func (s *Storage) GetCurrencyBy(fromSymbol, toSymbol []string) ([]*entities.Currency, error) {
	keys := make([]interface{}, len(fromSymbol)*len(toSymbol))
	result := make([]*entities.Currency, 0)
	if len(keys) == 0 {
		return result, nil
	}

	queryCase := strings.Repeat("?,", len(keys))
	queryCase = queryCase[0 : len(queryCase)-1]

	for i := range fromSymbol {
		for j := range toSymbol {
			keys[i+j] = fmt.Sprintf("%s%s", fromSymbol[i], toSymbol[j])
		}
	}

	res, err := s.conn.Query(
		fmt.Sprintf("SELECT `value` FROM currency "+
			"WHERE `updated_at` >  DATE_SUB(NOW(),INTERVAL 2 MINUTE) "+
			"AND `key` IN (%s)", queryCase),
		keys...)
	if err != nil {
		return nil, fmt.Errorf("failed to run query: %w", err)
	}
	defer res.Close()

	for res.Next() {
		var data []byte
		if err := res.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan data: %w", err)
		}

		var currency entities.Currency
		if err := jsoniter.Unmarshal(data, &currency); err != nil {
			return nil, fmt.Errorf("failed to unmarshal currency data: %w", err)
		}
		result = append(result, &currency)
	}

	return result, nil
}

func (s *Storage) Close() error {
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("failed to close mysql connection: %v", err)
	}
	return nil
}

func key(currency *entities.Currency) string {
	return fmt.Sprintf("%s%s", currency.FromSymbol, currency.ToSymbol)
}
