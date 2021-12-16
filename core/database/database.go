package database

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"bpm/core/config"
)

var db *sqlx.DB

func ConfigMysql() {
	host := config.ReadConfig("database.host")
	user := config.ReadConfig("database.user")
	password := config.ReadConfig("database.password")
	port := config.ReadConfig("database.port")
	dbname := config.ReadConfig("database.dbname")
	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?parseTime=true&loc=Local"
	mysqldb, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	mysqldb.SetMaxOpenConns(200)
	mysqldb.SetMaxIdleConns(10)
	db = mysqldb
}

func InitMySQL() *sqlx.DB {
	return db
}
