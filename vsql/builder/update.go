package builder

import (
	"context"
	"fmt"
	"strings"
)

const (
	// Update
	updateFormat = "UPDATE %s SET "
	setField     = "=?,"
	whereCond    = " WHERE "
	whereField   = "=? "
	andField     = "AND "

	// Delete
	deleteFormat = "DELETE FROM %s "
)

func BuildUpdate(ctx context.Context, table string, cond, data map[string]interface{}) (string, []interface{}, error) {
	var (
		sql  string
		args []interface{}
	)
	sql = fmt.Sprintf(updateFormat, table)
	for k, v := range data {
		sql += k + setField
		args = append(args, v)
	}
	sql = strings.TrimRight(sql, comma)
	if cond != nil {
		// todo: add support for multiple conditions
		sql += whereCond
		for k, v := range cond {
			sql += k + whereField + andField
			args = append(args, v)
		}
		sql = strings.TrimRight(sql, andField)
	}
	return sql, args, nil
}

func BuildDelete(ctx context.Context, table string, cond map[string]interface{}) (string, []interface{}, error) {
	var (
		sql  string
		args []interface{}
	)
	sql = fmt.Sprintf(deleteFormat, table)
	if cond != nil {
		sql += whereCond
		for k, v := range cond {
			sql += k + whereField + andField
			args = append(args, v)
		}
		sql = strings.TrimRight(sql, andField)
	}
	return sql, args, nil
}
