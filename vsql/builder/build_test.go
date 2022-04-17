package builder

import (
	"context"
	"fmt"
	"testing"
)

func TestBuildInsert(t *testing.T) {

	query, vals, err := BuildInsert(context.Background(), "users", []map[string]interface{}{
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

	query, vals, err := BuildUpdate(context.Background(), "users",
		map[string]interface{}{
			"name": "John",
			"age":  30,
		},
		map[string]interface{}{
			"name": "Jane",
		})

	fmt.Println(query, vals, err)
}

func TestBuildDelete(t *testing.T) {

	query, vals, err := BuildDelete(context.Background(), "users",
		map[string]interface{}{
			"name": "John",
			"age":  30,
		})

	fmt.Println(query, vals, err)

}

func TestBuildSelect(t *testing.T) {

}
