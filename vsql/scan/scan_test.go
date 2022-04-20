package scan

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type Data struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	CreateAt string `json:"create_at" db:"create_at"`
}

func TestNewScanner(t *testing.T) {

	db, err := sql.Open("mysql", "root:woaini12@tcp(127.0.0.1:3306)/test_table")
	if err != nil {
		fmt.Println(err)
		return
	}

	rows, err := db.QueryContext(context.Background(), "select * from test_table")
	if err != nil {
		fmt.Println(err)
		return
	}

	var d []*Data

	fmt.Println(ScannerIns.Scan(rows, &d))
	for _, val := range d {
		fmt.Println(fmt.Sprintf("%+v", val))
	}

}
