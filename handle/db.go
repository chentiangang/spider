package handle

import (
	"database/sql"
	"spider/config"

	_ "github.com/go-sql-driver/mysql"
)

func NewConn(cfg config.StorageConfig) *sql.DB {
	//dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
	//	cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	//// fmt.Println(dns)
	//d, err := sql.Open("mysql", dns)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = d.Ping()
	//if err != nil {
	//	xlog.Debug("%s", dns)
	//	xlog.Error("数据库连接失败")
	//	d.Close()
	//	return nil
	//}
	//return d
	return nil
}
