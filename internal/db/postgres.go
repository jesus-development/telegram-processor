package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/lib/pq"
	"log/slog"
	"strconv"
	"telegram-processor/internal/config"
)

type Vector32Float pq.Float32Array

func NewDatabase(cfg *config.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s %s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.Additional)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open -> %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping -> %w", err)
	}

	slog.Info("Database connected successfully.")

	return db, nil
}

// Value копия pq.Float32Array но с квадратными скобками
func (v Vector32Float) Value() (driver.Value, error) {
	if v == nil {
		return nil, nil
	}

	if n := len(v); n > 0 {
		// There will be at least two curly brackets, N bytes of values,
		// and N-1 bytes of delimiters.
		b := make([]byte, 1, 1+2*n)
		b[0] = '['

		b = strconv.AppendFloat(b, float64(v[0]), 'f', -1, 32)
		for i := 1; i < n; i++ {
			b = append(b, ',')
			b = strconv.AppendFloat(b, float64(v[i]), 'f', -1, 32)
		}

		return string(append(b, ']')), nil
	}

	return "[]", nil
}
