package dao

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
//本番
func InitDB(dbUser, dbPwd, dbName, instanceConnectionName string) {
	socketDir, isSet := os.LookupEnv("DB_SOCKET_DIR")
    if !isSet {
        socketDir = "/cloudsql"
    }
	dsn := fmt.Sprintf("%s:%s@unix(%s/%s)/%s?parseTime=true", dbUser, dbPwd, socketDir, instanceConnectionName, dbName)
	_db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("fail: sql.Open, %v\n", err)
	}

	if err := _db.Ping(); err != nil {
		log.Fatalf("fail: _db.Ping, %v\n", err)
	}
	db = _db
}
func DB() *sql.DB {
	return db
}