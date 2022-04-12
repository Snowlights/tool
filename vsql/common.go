package vsql

import (
	"github.com/opentracing/opentracing-go"
	"github.com/xwb1989/sqlparser"
	"strings"
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

func setDBSpanTags(span opentracing.Span, cluster, schema, table, query string) {
	span.SetTag(vtrace.Cluster, cluster)
	span.SetTag(vtrace.Schema, schema)
	span.SetTag(vtrace.Table, table)
	span.LogKV(vtrace.Query, query)

	span.SetTag(vtrace.Component, vtrace.ComponentSQL)
	span.SetTag(vtrace.SpanKind, vtrace.SpanKindSQL)
}

func parseTable(query string) string {

	stem, err := sqlparser.Parse(query)
	if err != nil {
		return ""
	}

	switch node := stem.(type) {
	case *sqlparser.Select:
		tables := []string{}
		for _, from := range node.From {
			switch from := from.(type) {
			case *sqlparser.AliasedTableExpr:
				tables = append(tables, from.Expr.(sqlparser.TableName).Name.String())
			}
		}
		return strings.Join(tables, comma)
	case *sqlparser.Insert:
		return node.Table.Name.String()
	case *sqlparser.Update:
		return node.TableExprs[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
	case *sqlparser.Delete:
		return node.TableExprs[0].(*sqlparser.AliasedTableExpr).Expr.(sqlparser.TableName).Name.String()
	default:
		return ""
	}

}
