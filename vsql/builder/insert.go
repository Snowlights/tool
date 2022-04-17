package builder

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

const (
	// Insert is the insert statement
	insertFormat = "INSERT INTO %s (%s) VALUES %s"
	// Replace is the replace statement
	replaceInsert = "REPLACE INTO %s (%s) VALUES %s"

	leftBracket  = "("
	rightBracket = ")"
	oneField     = "?,"
	comma        = ","
)

var (
	// Insert is the insert statement
	InsertNoData = errors.New("no data to insert")
	// Replace is the replace statement
	ReplaceNoData = errors.New("no data to replace")
)

func BuildInsert(ctx context.Context, table string, dataList []map[string]interface{}) (string, []interface{}, error) {
	if len(dataList) < 1 {
		return "", nil, InsertNoData
	}

	dataZero := dataList[0]
	vals := make([]interface{}, 0, len(dataZero)*len(dataList))
	fields, fieldVal := make([]string, 0, len(dataZero)), make([]string, 0, len(dataList))
	for k := range dataZero {
		fields = append(fields, k)
	}
	oneDataVal := buildOneDataField(len(fields))

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

func BuildReplaceInsert(ctx context.Context, table string, dataList []map[string]interface{}) (string, []interface{}, error) {
	if len(dataList) < 1 {
		return "", nil, ReplaceNoData
	}

	dataZero := dataList[0]
	vals := make([]interface{}, 0, len(dataZero)*len(dataList))
	fields, fieldVal := make([]string, 0, len(dataZero)), make([]string, 0, len(dataList))
	for k := range dataZero {
		fields = append(fields, k)
	}
	oneDataVal := buildOneDataField(len(fields))

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

func buildOneDataField(length int) string {
	return leftBracket + strings.TrimRight(strings.Repeat(oneField, length), comma) + rightBracket
}
