package scan

import (
	"database/sql"
	"errors"
	"reflect"
	"github.com/Snowlights/tool/parse"
)

const (
	defaultTagOption TagOption = "db"
)

var (
	ErrNilRows           = errors.New("[scanner]: rows is nil")
	ErrTargetNotSettable = errors.New("[scanner]: target need a pointer")
)

type (
	Scanner   struct{}
	TagOption string
)

var ScannerIns *Scanner

func init() {
	ScannerIns = NewScanner()
}

func NewScanner() *Scanner {
	return &Scanner{}
}

func (s *Scanner) Scan(rows *sql.Rows, target interface{}, opts ...TagOption) error {
	rowData, err := s.parseRows(rows)
	if err != nil {
		return err
	}

	targetObj := reflect.ValueOf(target)
	if !targetObj.Elem().CanSet() {
		return ErrTargetNotSettable
	}
	length := len(rowData)
	valueArr := reflect.MakeSlice(targetObj.Elem().Type(), 0, length)
	typeValue := valueArr.Type().Elem()

	if opts == nil {
		opts = []TagOption{defaultTagOption}
	}
	for i := 0; i < length; i++ {
		newValue := reflect.New(typeValue.Elem())
		err = parse.UnmarshalKV(rowData[i], newValue.Interface(), string(opts[0]))
		if nil != err {
			return err
		}
		valueArr = reflect.Append(valueArr, newValue)
	}
	targetObj.Elem().Set(valueArr)
	return nil
}

func (s *Scanner) parseRows(rows *sql.Rows) ([]map[string]string, error) {
	if rows == nil {
		return nil, ErrNilRows
	}
	columns, err := rows.Columns()
	if nil != err {
		return nil, err
	}
	length, values := len(columns), make([]interface{}, len(columns))
	var result []map[string]string
	for i := 0; i < length; i++ {
		values[i] = new(string)
	}

	for rows.Next() {
		err = rows.Scan(values...)
		if nil != err {
			return nil, err
		}
		res := make(map[string]string)
		for idx, name := range columns {
			res[name] = *(values[idx].(*string))
		}
		result = append(result, res)
	}

	return result, nil
}
