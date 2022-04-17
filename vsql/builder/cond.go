package builder

type cond interface {
	Build() (string, []interface{})
}

type condType string

const (
	equal        condType = "="
	notEqual     condType = "!="
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

	OrderByCond = "ORDER BY"
	HavingCond  = "HAVING"
	GroupByCond = "GROUP BY"
	LimitCond   = "LIMIT"
)
