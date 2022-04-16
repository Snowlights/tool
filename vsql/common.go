package vsql

import (
	"errors"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/opentracing/opentracing-go"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"strings"
	"vtool/vconfig"
	"vtool/vtrace"
)

const (
	dbBeginOperation = "sql.db.Begin"
	dbQueryOperation = "sql.db.QueryContext"
	dbExecOperation  = "sql.db.ExecContext"

	txQueryOperation    = "sql.tx.QueryContext"
	txExecOperation     = "sql.tx.ExecContext"
	txCommitOperation   = "sql.tx.Commit"
	txRollbackOperation = "sql.tx.Rollback"

	comma = ","
)

type ChangeType int64

const (
	Reopen ChangeType = 1
	Close  ChangeType = 2
	Reset  ChangeType = 3
)

type DBConfig struct {
	// cluster to db config
	Conf map[string]InstanceConfig `json:"conf" properties:"conf"`
}

type InstanceConfig struct {
	DBName       string `json:"db_name" properties:"db_name"`
	Host         string `json:"host" properties:"host"`
	Timeout      int64  `json:"timeout" properties:"timeout"`
	ReadTimeout  int64  `json:"read_timeout" properties:"read_timeout"`
	WriteTimeout int64  `json:"write_timeout" properties:"write_timeout"`
	MaxLifeTime  int64  `json:"max_life_time" properties:"max_life_time"`
	MaxIdleConn  int    `json:"max_idle_conn" properties:"max_idle_conn"`
	MaxOpenConn  int    `json:"max_open_conn" properties:"max_open_conn"`
	Username     string `json:"username" properties:"username"`
	Password     string `json:"password" properties:"password"`
}

func (ic InstanceConfig) buildInsCfgKey() string {
	return strings.Join([]string{ic.Host, ic.Username, ic.Password}, comma)
}

var (
	parserIns *parser.Parser

	NotInitManager  = errors.New("NotInitManager")
	NotFoundCluster = errors.New("NotFoundCluster")
)

func init() {
	parserIns = parser.New()
}

func setDBSpanTags(span opentracing.Span, cluster, schema, table, query string) {
	span.SetTag(vtrace.Cluster, cluster)
	span.SetTag(vtrace.Schema, schema)
	span.SetTag(vtrace.Table, table)
	span.LogKV(vtrace.Query, query)

	span.SetTag(vtrace.Component, vtrace.ComponentSQL)
	span.SetTag(vtrace.SpanKind, vtrace.SpanKindSQL)
}

type TableVisitor struct {
	table []string
}

func (t *TableVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	switch nn := n.(type) {
	case *ast.TableName:
		t.table = append(t.table, nn.Name.String())
	}
	return n, false
}

func (t *TableVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

func (t *TableVisitor) Table() []string {
	return t.table
}

func parseTable(sql string) string {
	stmt, err := parserIns.ParseOneStmt(sql, "", "")
	if err != nil {
		return ""
	}
	v := &TableVisitor{}
	stmt.Accept(v)
	return strings.Join(v.Table(), comma)
}

type MysqlListener struct {
	Change func()
}

func (cl *MysqlListener) OnChange(event *storage.ChangeEvent) {

}

func (cl *MysqlListener) OnNewestChange(event *storage.FullChangeEvent) {
	if event.Namespace != vconfig.ServerDB {
		return
	}
	cl.Change()
}
