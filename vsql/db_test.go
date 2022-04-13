package vsql

import (
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"testing"
)

func TestGetDB(t *testing.T) {

	// sql := "select t1.id from t1, t3, t2 where t2.id = 1 and t2.id = t1.id and t1.name = t3.name"
	sql := "update `user` set `name` = 'test' where `id` = 1 and `name` in (select `name` from `data`)"

	fmt.Println(parseTable(sql))
}

func TestDB(t *testing.T) {

	parse := parser.New()

	// sql := ""select t1.id from t1, t3, t2 where t2.id = 1 and t2.id = t1.id and t1.name = t3.name""
	// sql := "select * from `user`, `data` where `id` = 1"
	// sql := "insert into `user` (`id`, `name`) values (1, 'test')"
	sql := "update `user` set `name` = 'test' where `id` = 1 and `name` in (select `name` from `data`)"

	stmt, err := parse.ParseOneStmt(sql, "", "")
	if err != nil {
		fmt.Println(err)
	}

	var table TableNameVisitor

	fmt.Println(stmt.Accept(&table))
}

type TableNameVisitor struct {
	table []string
}

func (f *TableNameVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	switch nn := n.(type) {
	case *ast.TableName:
		fmt.Println(nn.Name.String())
	}

	// fmt.Printf("%T\n", n)
	return n, false
}

func (f *TableNameVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

func (f *TableNameVisitor) TableName() []string {
	return f.table
}

func TestParseTable(t *testing.T) {

	// sql := "select * from `user`, `data` where `id` = 1"
	// sql := "insert into `user` (`id`, `name`) values (1, 'test')"
	sql := "update `user` set `name` = 'test' where `id` = 1"

	fmt.Println(parseTable(sql))

}
