package vsql

import (
	"context"
	"database/sql"
	"github.com/opentracing/opentracing-go"
)

// only support for one schema

type DBContent interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type DB struct {
	db              *sql.DB
	cluster, schema string
}

func (db *DB) Begin(ctx context.Context) (*TX, error) {

	span, ctx := opentracing.StartSpanFromContext(ctx, dbBeginOperation)
	defer span.Finish()
	setDBSpanTags(span, db.cluster, db.schema, "", "")

	tx, err := db.db.Begin()
	if err != nil {
		return nil, err
	}

	return &TX{
		tx:      tx,
		cluster: db.cluster,
		schema:  db.schema,
	}, nil
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, dbQueryOperation)
	defer span.Finish()
	setDBSpanTags(span, db.cluster, db.schema, parseTable(query), query)

	return db.db.QueryContext(ctx, query, args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, dbExecOperation)
	defer span.Finish()
	setDBSpanTags(span, db.cluster, db.schema, parseTable(query), query)

	return db.db.ExecContext(ctx, query, args...)
}

type TX struct {
	tx              *sql.Tx
	cluster, schema string
}

func (tx *TX) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, txQueryOperation)
	defer span.Finish()
	setDBSpanTags(span, tx.cluster, tx.schema, parseTable(query), query)

	return tx.tx.QueryContext(ctx, query, args...)
}

func (tx *TX) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, txExecOperation)
	defer span.Finish()
	setDBSpanTags(span, tx.cluster, tx.schema, parseTable(query), query)

	return tx.tx.ExecContext(ctx, query, args...)
}

func (tx *TX) Commit(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, txCommitOperation)
	defer span.Finish()
	setDBSpanTags(span, tx.cluster, tx.schema, "", "")

	return tx.tx.Commit()
}

func (tx *TX) Rollback(ctx context.Context) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, txRollbackOperation)
	defer span.Finish()
	setDBSpanTags(span, tx.cluster, tx.schema, "", "")

	return tx.tx.Rollback()
}
