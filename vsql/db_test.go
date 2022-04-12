package vsql

import (
	"fmt"
	"testing"
)

func TestDB(t *testing.T) {

	fmt.Println(utf8mb4DSN())

}

func TestParseTable(t *testing.T) {

	// sql := "select * from `user`, `data` where `id` = 1"
	// sql := "insert into `user` (`id`, `name`) values (1, 'test')"
	sql := "update `user` set `name` = 'test' where `id` = 1"

	fmt.Println(parseTable(sql))

}
