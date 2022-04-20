package builder

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type cond interface {
	Build() ([]string, []interface{})
}

type (
	condType  string
	operation func(map[string]interface{}) (cond, error)
)

const (
	equal        condType = "="
	notEqual     condType = "!="
	notEqualV2   condType = "<>"
	less         condType = "<"
	lessEqual    condType = "<="
	greater      condType = ">"
	greaterEqual condType = ">="

	like       condType = "like"
	notLike    condType = "not like"
	in         condType = "in"
	notIn      condType = "not in"
	between    condType = "between"
	notBetween condType = "not between"
	isNull     condType = "is null"
	isNotNull  condType = "is not null"

	OrderByKey  = "_order_by"
	OrderByCond = "order by"
	HavingKey   = "_having"
	HavingCond  = "having"
	GroupByKey  = "_group_by"
	GroupByCond = "group by"
	LimitKey    = "_limit"
	LimitCond   = "limit"
)

var condOperationMap = map[condType]operation{
	equal: func(m map[string]interface{}) (cond, error) {
		return equalCond(m), nil
	},
	notEqual: func(m map[string]interface{}) (cond, error) {
		return notEqualCond(m), nil
	},
	notEqualV2: func(m map[string]interface{}) (cond, error) {
		return notEqualCond(m), nil
	},
	less: func(m map[string]interface{}) (cond, error) {
		return lessCond(m), nil
	},
	lessEqual: func(m map[string]interface{}) (cond, error) {
		return lessEqualCond(m), nil
	},
	greater: func(m map[string]interface{}) (cond, error) {
		return greaterCond(m), nil
	},
	greaterEqual: func(m map[string]interface{}) (cond, error) {
		return greaterEqualCond(m), nil
	},
	like: func(m map[string]interface{}) (cond, error) {
		return likeCond(m), nil
	},
	notLike: func(m map[string]interface{}) (cond, error) {
		return notLikeCond(m), nil

	},
	in: func(m map[string]interface{}) (cond, error) {
		wp, err := convertWhereMapToWhereMapSlice(m)
		if nil != err {
			return nil, err
		}
		return inCond(wp), nil
	},
	notIn: func(m map[string]interface{}) (cond, error) {
		wp, err := convertWhereMapToWhereMapSlice(m)
		if nil != err {
			return nil, err
		}
		return notInCond(wp), nil
	},
	between: func(m map[string]interface{}) (cond, error) {
		wp, err := convertWhereMapToWhereMapSlice(m)
		if nil != err {
			return nil, err
		}
		return betweenCond(wp), nil
	},
	notBetween: func(m map[string]interface{}) (cond, error) {
		wp, err := convertWhereMapToWhereMapSlice(m)
		if nil != err {
			return nil, err
		}
		return notBetweenCond(wp), nil
	},
	isNull: func(m map[string]interface{}) (cond, error) {
		return nullCond(m), nil
	},
	isNotNull: func(m map[string]interface{}) (cond, error) {
		return notNullCond(m), nil
	},
}

func convertWhereMapToWhereMapSlice(where map[string]interface{}) (map[string][]interface{}, error) {
	result := make(map[string][]interface{})
	for key, val := range where {
		vals, ok := convertInterfaceToMap(val)
		if !ok {
			return nil, errWhereInType
		}
		if 0 == len(vals) {
			return nil, errEmptyINCondition
		}
		result[key] = vals
	}
	return result, nil
}

func convertInterfaceToMap(val interface{}) ([]interface{}, bool) {
	s := reflect.ValueOf(val)
	if s.Kind() != reflect.Slice {
		return nil, false
	}
	interfaceSlice := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		interfaceSlice[i] = s.Index(i).Interface()
	}
	return interfaceSlice, true
}

type (
	equalCond        map[string]interface{}
	notEqualCond     map[string]interface{}
	lessCond         map[string]interface{}
	lessEqualCond    map[string]interface{}
	greaterCond      map[string]interface{}
	greaterEqualCond map[string]interface{}

	likeCond       map[string]interface{}
	notLikeCond    map[string]interface{}
	inCond         map[string][]interface{}
	notInCond      map[string][]interface{}
	betweenCond    map[string][]interface{}
	notBetweenCond map[string][]interface{}
	nullCond       map[string]interface{}
	notNullCond    map[string]interface{}
)

func (eq equalCond) Build() ([]string, []interface{}) {
	return build(eq, string(equal))
}

func (neq notEqualCond) Build() ([]string, []interface{}) {
	return build(neq, string(notEqual))
}

func (l lessCond) Build() ([]string, []interface{}) {
	return build(l, string(less))
}

func (leq lessEqualCond) Build() ([]string, []interface{}) {
	return build(leq, string(lessEqual))
}

func (g greaterCond) Build() ([]string, []interface{}) {
	return build(g, string(greater))
}

func (geq greaterEqualCond) Build() ([]string, []interface{}) {
	return build(geq, string(greaterEqual))
}

func (l likeCond) Build() ([]string, []interface{}) {
	if l == nil || len(l) == 0 {
		return nil, nil
	}
	var conds []string
	var vals []interface{}
	for k := range l {
		conds = append(conds, k)
	}
	for j := 0; j < len(conds); j++ {
		val := l[conds[j]]
		conds[j] = conds[j] + space + string(like) + space + placeholder
		vals = append(vals, val)
	}
	return conds, vals
}

func (nl notLikeCond) Build() ([]string, []interface{}) {
	if nl == nil || len(nl) == 0 {
		return nil, nil
	}
	var conds []string
	var vals []interface{}
	for k := range nl {
		conds = append(conds, k)
	}
	for j := 0; j < len(conds); j++ {
		val := nl[conds[j]]
		conds[j] = conds[j] + space + string(notLike) + space + placeholder
		vals = append(vals, val)
	}
	return conds, vals
}

func (i inCond) Build() ([]string, []interface{}) {
	if i == nil || len(i) == 0 {
		return nil, nil
	}
	var conds []string
	var vals []interface{}
	for k := range i {
		conds = append(conds, k)
	}
	for j := 0; j < len(conds); j++ {
		val := i[conds[j]]
		conds[j] = buildIn(conds[j], string(in), val)
		vals = append(vals, val...)
	}
	return conds, vals
}

func (ni notInCond) Build() ([]string, []interface{}) {
	if ni == nil || len(ni) == 0 {
		return nil, nil
	}
	var conds []string
	var vals []interface{}
	for k := range ni {
		conds = append(conds, k)
	}
	for j := 0; j < len(conds); j++ {
		val := ni[conds[j]]
		conds[j] = buildIn(conds[j], string(notIn), val)
		vals = append(vals, val...)
	}
	return conds, vals
}

func (b betweenCond) Build() ([]string, []interface{}) {
	if b == nil || len(b) == 0 {
		return nil, nil
	}
	var conds []string
	var vals []interface{}
	for k := range b {
		conds = append(conds, k)
	}
	for j := 0; j < len(conds); j++ {
		val := b[conds[j]]
		condJ, err := buildBetween(conds[j], string(between), val)
		if nil != err {
			continue
		}
		conds[j] = condJ
		vals = append(vals, val...)
	}
	return conds, vals
}

func (nb notBetweenCond) Build() ([]string, []interface{}) {
	if nb == nil || len(nb) == 0 {
		return nil, nil
	}
	var conds []string
	var vals []interface{}
	for k := range nb {
		conds = append(conds, k)
	}
	for j := 0; j < len(conds); j++ {
		val := nb[conds[j]]
		condJ, err := buildBetween(conds[j], string(notBetween), val)
		if nil != err {
			continue
		}
		conds[j] = condJ
		vals = append(vals, val...)
	}
	return conds, vals
}

func (n nullCond) Build() ([]string, []interface{}) {
	if n == nil || len(n) == 0 {
		return nil, nil
	}
	var conds []string
	var vals []interface{}
	for k := range n {
		conds = append(conds, k)
	}
	for j := 0; j < len(conds); j++ {
		conds[j] = conds[j] + space + string(isNull)
	}
	return conds, vals
}

func (nn notNullCond) Build() ([]string, []interface{}) {
	if nn == nil || len(nn) == 0 {
		return nil, nil
	}
	var conds []string
	var vals []interface{}
	for k := range nn {
		conds = append(conds, k)
	}
	for j := 0; j < len(conds); j++ {
		conds[j] = conds[j] + space + string(isNotNull)
	}
	return conds, vals
}

func buildBetween(key, op string, vals []interface{}) (string, error) {
	if len(vals) != 2 {
		return "", errors.New("vals of between must be a slice with two elements")
	}
	return fmt.Sprintf("(%s %s ? AND ?)", key, op), nil
}

func buildIn(field, op string, vals []interface{}) (conds string) {
	conds = strings.TrimRight(strings.Repeat(placeholder+comma, len(vals)), comma)
	conds = fmt.Sprintf("%s %s (%s)", field, op, conds)
	return
}

func build(m map[string]interface{}, op string) ([]string, []interface{}) {
	if nil == m || 0 == len(m) {
		return nil, nil
	}
	length := len(m)
	conds := make([]string, length)
	vals := make([]interface{}, length)
	var i int
	for key := range m {
		conds[i] = key
		i++
	}
	for i = 0; i < length; i++ {
		vals[i] = m[conds[i]]
		conds[i] = conds[i] + op + placeholder
	}
	return conds, vals
}

func OrderBy(orderBy string) string {
	return fmt.Sprintf("%s %s", OrderByCond, orderBy)
}

func Having(having string) string {
	return fmt.Sprintf("%s %s", HavingCond, having)
}

func GroupBy(groupBy string) string {
	return fmt.Sprintf("%s %s", GroupByCond, groupBy)
}

// limit not support negative number
func Limit(val interface{}) (string, []interface{}, error) {
	v := reflect.ValueOf(val)
	if v.Kind() != reflect.Slice {
		return "", nil, errLimitValueType
	}
	if v.Len() != 2 {
		return "", nil, errLimitValueLength
	}
	arr := make([]interface{}, 0, 2)
	switch vals := val.(type) {
	case []uint64:
		arr = append(arr, vals[0], vals[1])
	case []uint:
		arr = append(arr, vals[0], vals[1])
	case []uint8:
		arr = append(arr, vals[0], vals[1])
	case []uint16:
		arr = append(arr, vals[0], vals[1])
	case []uint32:
		arr = append(arr, vals[0], vals[1])
	default:
		return "", nil, errLimitValueType
	}
	return fmt.Sprintf("%s %s,%s", LimitCond, placeholder, placeholder), arr, nil
}
