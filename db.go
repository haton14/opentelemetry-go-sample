package main

import (
	"net"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func connectDB() (*sqlx.DB, error) {
	conf := mysql.NewConfig()

	conf.Net = "tcp"
	conf.Addr = net.JoinHostPort("127.0.0.1", "3306")
	conf.User = "mysql"
	conf.Passwd = "password"
	conf.DBName = "jaeger"
	conf.ParseTime = true

	dsn := conf.FormatDSN()

	// db, err := sqlx.Open("mysql", conf.FormatDSN())
	db, err := otelsqlx.Open("mysql", dsn,
		otelsql.WithAttributes(semconv.DBSystemMySQL),
		otelsql.WithAttributes(attribute.KeyValue{Key: "service.name", Value: attribute.StringValue("jaeger_db")}),
		otelsql.WithDBName("jaeger"),
	)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
