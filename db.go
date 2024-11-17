package main

import (
	"net"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func connectDB() (*sqlx.DB, error) {
	conf := mysql.NewConfig()

	conf.Net = "tcp"
	conf.Addr = net.JoinHostPort("127.0.0.1", "3306")
	conf.User = "mysql"
	conf.Passwd = "password"
	conf.DBName = "jaeger"
	conf.ParseTime = true

	db, err := sqlx.Open("mysql", conf.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
