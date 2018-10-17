package library

import (
	"github.com/revel/config"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"fmt"
	"time"
)

var Db *sql.DB

func InitMysql(config *config.Context) {
	config.SetSection("mysql")
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&timeout=10s",
		config.StringDefault("user", "root"),
		config.StringDefault("passwd", "root"),
		config.StringDefault("host", "127.0.0.1"),
		config.StringDefault("port", "3306"),
		config.StringDefault("dbname", "test"),
		config.StringDefault("charset", "utf8"),
		)
	var err error
	Db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		panic("mysql connect error" + err.Error())
		return;
	}
	Db.SetMaxIdleConns(2)
	Db.SetMaxOpenConns(config.IntDefault("maxOpenConns", 1000))
	config.SetSection("DEFAULT")
	Db.SetConnMaxLifetime(time.Duration(60) * time.Second)
}
