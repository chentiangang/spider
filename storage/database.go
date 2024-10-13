package storage

import "database/sql"

// 数据库配置
type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"database"`
}

func NewDatabaseWithName(name string) *sql.DB {

}
