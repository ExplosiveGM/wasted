package admin

import (
	"database/sql"
	"fmt"

	"github.com/spf13/viper"
)

func Connect() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres?sslmode=%s",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_SSLMODE"),
	)

	return sql.Open("pgx", connStr)
}
