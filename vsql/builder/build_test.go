package builder

import (
	"context"
	"fmt"
	"testing"
)

func TestBuildInsert(t *testing.T) {

	query, vals, err := BuilderIns.BuildInsert(context.Background(), "users", []map[string]interface{}{
		{
			"name": "John",
			"age":  30,
		},
		{
			"name": "Jane",
			"age":  25,
		},
	})

	fmt.Println(query, vals, err)
}

func TestBuildUpdate(t *testing.T) {

	query, vals, err := BuilderIns.BuildUpdate(context.Background(), "users",
		map[string]interface{}{
			"name":  "John",
			"age":   30,
			"id in": []int{1, 2, 3, 4, 5},
			"id !=": 3,
			"id >":  5,
		},
		map[string]interface{}{
			"name": "Jane",
		})

	fmt.Println(query, vals, err)
}

func TestBuildDelete(t *testing.T) {

	query, vals, err := BuilderIns.BuildDelete(context.Background(), "users",
		map[string]interface{}{
			"name": "John",
			"age":  30,
		})

	fmt.Println(query, vals, err)

}

func TestBuildSelect(t *testing.T) {

	query, args, err := BuilderIns.BuildSelect(context.Background(), "users",
		map[string]interface{}{
			"name":           "John",
			"age between":    []int{30, 25},
			"id in":          []int{1, 2, 3},
			"id not in":      []int{4, 5, 6},
			"id":             1,
			"id !=":          2,
			"id <>":          3,
			"id <=":          4,
			"id >=":          5,
			"id <":           6,
			"id >":           7,
			"id is null":     true,
			"id is not null": true,
			OrderByKey:       "id desc",
			LimitKey:         []uint64{1, 2},
			GroupByKey:       "name",
		})

	fmt.Println(query, args, err)
}
