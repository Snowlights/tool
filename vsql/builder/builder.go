package builder

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	placeholder = "?"
	space       = " "

	setField   = "=?,"
	whereCond  = " WHERE "
	whereField = "=? "
	andField   = " AND "

	leftBracket  = "("
	rightBracket = ")"
	oneField     = "?,"
	comma        = ","
)

var (
	ErrNotFoundOperator = errors.New("[builder] not found operator")
	errSplitEmptyKey    = errors.New("[builder] couldn't split a empty string")
	errWhereInType      = errors.New(`[builder] the value of "xxx in" must be of []interface{} type`)
	errEmptyINCondition = errors.New(`[builder] the value of "in" must contain at least one element`)
	errLimitValueLength = errors.New(`[builder] the value of "_limit" must contain two uint elements`)
	errLimitValueType   = errors.New(`[builder] the value of "_limit" must be of []uint、[]uint8、[]uint16、[]uint32、[]uint64 type`)
)

var (
	// Insert is the insert statement
	InsertNoData = errors.New("[builder] no data to insert")
	// Replace is the replace statement
	ReplaceNoData = errors.New("[builder] no data to replace")
)

var BuilderIns *Builder

type Builder struct {
}

func init() {
	BuilderIns = NewBuilder()
}

func NewBuilder() *Builder {
	return &Builder{}
}

const (
	// Insert is the insert statement
	insertFormat = "INSERT INTO %s (%s) VALUES %s"
	// Replace is the replace statement
	replaceInsert = "REPLACE INTO %s (%s) VALUES %s"

	// Update
	updateFormat = "UPDATE %s SET "

	// Delete
	deleteFormat = "DELETE FROM %s "

	// Select
	selectFormat = "SELECT * FROM %s "
)

func (b *Builder) BuildInsert(ctx context.Context, table string, dataList []map[string]interface{}) (string, []interface{}, error) {
	if len(dataList) < 1 {
		return "", nil, InsertNoData
	}

	dataZero := dataList[0]
	vals := make([]interface{}, 0, len(dataZero)*len(dataList))
	fields, fieldVal := make([]string, 0, len(dataZero)), make([]string, 0, len(dataList))
	for k := range dataZero {
		fields = append(fields, k)
	}
	oneDataVal := b.buildOneDataField(len(fields))

	for _, data := range dataList {
		fieldVal = append(fieldVal, oneDataVal)
		for _, field := range fields {
			val, ok := data[field]
			if !ok {
				return "", nil, fmt.Errorf("field %s not found in data", field)
			}
			vals = append(vals, val)
		}
	}
	return fmt.Sprintf(insertFormat, table, strings.Join(fields, comma), strings.Join(fieldVal, comma)), vals, nil
}

func (b *Builder) BuildReplaceInsert(ctx context.Context, table string, dataList []map[string]interface{}) (string, []interface{}, error) {
	if len(dataList) < 1 {
		return "", nil, ReplaceNoData
	}

	dataZero := dataList[0]
	vals := make([]interface{}, 0, len(dataZero)*len(dataList))
	fields, fieldVal := make([]string, 0, len(dataZero)), make([]string, 0, len(dataList))
	for k := range dataZero {
		fields = append(fields, k)
	}
	oneDataVal := b.buildOneDataField(len(fields))

	for _, data := range dataList {
		fieldVal = append(fieldVal, oneDataVal)
		for _, field := range fields {
			val, ok := data[field]
			if !ok {
				return "", nil, fmt.Errorf("field %s not found in data", field)
			}
			vals = append(vals, val)
		}
	}
	return fmt.Sprintf(replaceInsert, table, strings.Join(fields, comma), strings.Join(fieldVal, comma)), vals, nil
}

func (b *Builder) BuildUpdate(ctx context.Context, table string, cond, data map[string]interface{}) (string, []interface{}, error) {
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
		whereSql, whereArgs, err := b.buildWhereCond(cond)
		if err != nil {
			return "", nil, err
		}
		sql += whereSql
		args = append(args, whereArgs...)
	}
	return sql, args, nil
}

func (b *Builder) BuildDelete(ctx context.Context, table string, cond map[string]interface{}) (string, []interface{}, error) {
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

func (b *Builder) BuildSelect(ctx context.Context, table string, cond map[string]interface{}) (string, []interface{}, error) {
	var (
		sql  string
		args []interface{}
	)
	sql = fmt.Sprintf(selectFormat, table)
	if len(cond) > 0 {
		whereSql, whereArgs, err := b.buildWhereCond(cond)
		if err != nil {
			return "", nil, err
		}
		sql += whereSql
		args = append(args, whereArgs...)
	}
	return sql, args, nil
}

func (b *Builder) buildOneDataField(length int) string {
	return leftBracket + strings.TrimRight(strings.Repeat(oneField, length), comma) + rightBracket
}

func (b *Builder) buildWhereCond(conds map[string]interface{}) (string, []interface{}, error) {
	if len(conds) == 0 {
		return "", nil, nil
	}

	var groupBy, having, orderBy, limit string
	var limitValue []interface{}
	var whereConds []string
	var whereValue []interface{}
	for key, val := range conds {
		key = strings.Trim(key, space)
		switch key {
		case OrderByKey:
			orderBy = OrderBy(val.(string))
		case GroupByKey:
			groupBy = GroupBy(val.(string))
		case HavingKey:
			having = Having(val.(string))
		case LimitKey:
			limitRes, limitVals, err := Limit(val)
			if err != nil {
				return "", nil, err
			}
			limit = limitRes
			limitValue = append(limitValue, limitVals...)
		default:
			field, operator, err := b.splitKeyAndGetCondType(key)
			if err != nil {
				return "", nil, err
			}
			operationIns, ok := condOperationMap[operator]
			if !ok {
				return "", nil, fmt.Errorf("%s key is %s", ErrNotFoundOperator, key)
			}
			condIns, err := operationIns(map[string]interface{}{field: val})
			if err != nil {
				return "", nil, err
			}
			where, whereVals := condIns.Build()
			whereConds = append(whereConds, where...)
			whereValue = append(whereValue, whereVals...)
		}
	}

	whereSql := whereCond + strings.Join(whereConds, andField)
	if orderBy != "" {
		whereSql = whereSql + space + orderBy
	}
	if groupBy != "" {
		whereSql = whereSql + space + groupBy
	}
	if having != "" {
		whereSql = whereSql + space + having
	}
	if limit != "" {
		whereSql = whereSql + space + limit
		whereValue = append(whereValue, limitValue...)
	}

	return whereSql, whereValue, nil
}

func (b *Builder) splitKeyAndGetCondType(key string) (field string, operator condType, err error) {
	key = strings.Trim(key, space)
	if key == "" {
		err = errSplitEmptyKey
		return
	}
	idx := strings.IndexByte(key, ' ')
	if idx == -1 {
		field = key
		operator = equal
	} else {
		field = key[:idx]
		operator = condType(strings.Trim(key[idx+1:], space))
	}
	if operator == "" {
		err = ErrNotFoundOperator
		return
	}
	return
}
