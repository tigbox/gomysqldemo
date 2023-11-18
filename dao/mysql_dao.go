package dao

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

var mysqlDao *MysqlDao

type MysqlDao struct {
	db *sqlx.DB
}

func Init() {
	db, err := sqlx.Connect("mysql", "root:123456@(127.0.0.1:3306)/demo")
	if err != nil {
		log.Fatalln(err)
	}
	mysqlDao = &MysqlDao{db: db}
}

func NewMysqlDao(ctx context.Context) *MysqlDao {
	if mysqlDao == nil {
		Init()
	}
	return mysqlDao
}
