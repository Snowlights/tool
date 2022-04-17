package builder

import "context"

func BuildSelect(ctx context.Context, table string, cond map[string]interface{}) (string, []interface{}) {
	var (
		sql  string
		args []interface{}
	)
	sql = "SELECT * FROM " + table
	if len(cond) > 0 {
		sql += " WHERE "
		for k, v := range cond {
			sql += k + " = ? AND "
			args = append(args, v)
		}
		sql = sql[:len(sql)-5]
	}
	return sql, args
}
