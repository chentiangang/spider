package handle

import (
	"database/sql"
	"fmt"
	"spider/config"

	"github.com/chentiangang/xlog"
	_ "github.com/go-sql-driver/mysql"
)

func NewConn(cfg config.MySQLConfig) *sql.DB {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	// fmt.Println(dns)
	d, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}

	err = d.Ping()
	if err != nil {
		xlog.Error("数据库连接失败,%s", dns)
		return nil
	}
	return d
}
