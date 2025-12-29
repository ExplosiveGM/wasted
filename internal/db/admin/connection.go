package admin

import (
	"database/sql"
	"fmt"

	"github.com/ExplosiveGM/wasted/config"
)

func Connect(dbConfig *config.DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/postgres?sslmode=%s",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.SslMode,
	)

	return sql.Open("pgx", connStr)
}
