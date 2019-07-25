package models

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var (
	Db *sql.DB
)

//InitMysql create connection to mysql server
func InitMysql(dsn string) {
	var err error
	if len(dsn) == 0 {
		logrus.Fatalln("mysql DSN is empty.")
	}

	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		logrus.Fatalln("Open mysql error: ", err)
	}

	err = Db.Ping()
	if err != nil {
		logrus.Fatalln("Ping mysql error: ", err)
	}

	logrus.Infoln("Connected to mysql")

	Db.SetMaxIdleConns(20)
	Db.SetMaxOpenConns(20)
}

//CloseMysql closes mysql connection
func CloseMysql() {
	Db.Close()
}
