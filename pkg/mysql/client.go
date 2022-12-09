package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func NewClient(username, password, ip, port, database string) *sql.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, ip, port, database)
	db, err := sql.Open("mysql", url)

	if err != nil {
		fmt.Print(err)
		//todo logger err
	}

	return db
}
