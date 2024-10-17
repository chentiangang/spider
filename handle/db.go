package handle

import (
	"database/sql"
	"spider/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewConn(cfg config.StorageConfig) *sql.DB {
	return nil
}
