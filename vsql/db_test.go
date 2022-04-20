package vsql

import (
	"fmt"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	_ "github.com/pingcap/tidb/types/parser_driver"
	driver "github.com/pingcap/tidb/types/parser_driver"
	"strings"
	"testing"
	"time"
)

func TestGetDB(t *testing.T) {

	sql := "update `user` set `name` = 'test' where `id` = 1 and `name` in (select `name` from `data`)"

	// sql := "select t1.id from t1, t3, t2 where t2.id = 1 and t2.id = t1.id and t1.name = t3.name"
	// sql := "update `user` set `name` = 'test' where `id` = 1 and `name` in (select `name` from `data`)"

	fmt.Println(parseTable(sql))
}

func TestDB(t *testing.T) {

	parse := parser.New()

	//sql := "select t1.id from t1, t3, t2 where t2.id = 1 and t2.id = t1.id and t1.name = t3.name"
	sql := "select * from `user`, `data` where `name` in ('test', 'test2') and `id` in(1,2,3) "
	// sql := "insert into `user` (`id`, `name`) values (1, 'test')"
	// sql := "update `user` set `name` = 'test' where `id` = 1 and `name` in (select `name` from `data`)"

	stmt, err := parse.ParseOneStmt(sql, "", "")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(ExtractAstVisit(stmt))

	//var table TableNameVisitor
	//
	//fmt.Println(stmt.Accept(&table))
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

type AstVisitor struct {
	sqlFeature string
}

func ExtractAstVisit(stmt ast.StmtNode) (*AstVisitor, error) {
	visitor := &AstVisitor{}

	stmt.Accept(visitor)

	sb := strings.Builder{}
	if err := stmt.Restore(format.NewRestoreCtx(format.DefaultRestoreFlags, &sb)); err != nil {
		return nil, err
	}
	visitor.sqlFeature = sb.String()

	return visitor, nil
}

func (f *AstVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	switch nn := n.(type) {
	case *ast.PatternInExpr:
		if len(nn.List) == 0 {
			return nn, false
		}
		if _, ok := nn.List[0].(*driver.ValueExpr); ok {
			nn.List = nn.List[:1]
		}
	case *driver.ValueExpr:
		nn.SetValue("?")
	}
	return n, false
}

func (f *AstVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

func (f *AstVisitor) SqlFeature() string {
	return f.sqlFeature
}

func TestParseTable(t *testing.T) {

	// sql := "select * from `user`, `data` where `id` = 1"
	// sql := "insert into `user` (`id`, `name`) values (1, 'test')"
	//sql := "update `user` set `name` = 'test' where `id` = 1"
	//
	//fmt.Println(parseTable(sql))

	tt := time.Second * time.Duration(10)

	fmt.Println(fmt.Sprintf(dsnFormat, "timeout", tt))

}
