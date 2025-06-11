

package dao

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB(dbUser, dbPwd, dbName, instanceConnectionName string) {
	godotenv.Load()

	dsn := fmt.Sprintf("%s:%s@unix(%s)/%s?parseTime=true", dbUser, dbPwd, instanceConnectionName, dbName)
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