package vredis

import (
	"context"
	"testing"
)

func TestRedisClient_Set(t *testing.T) {

	ctx := context.Background()
	client, err := NewRedisClient(context.Background(), &RedisConfig{
		Host: "127.0.0.1:6379",
		Auth: "",
		Db:   0,
	})

	if err != nil {
		t.Error(err)
		return
	}

	cmd := client.Set(ctx, "test", "test", 0)
	res, err := cmd.Result()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(res)
}

func TestNewRedisClient(t *testing.T) {

	ctx := context.Background()
	client, err := NewRedisClient(ctx, &RedisConfig{
		Host: "127.0.0.1:6379",
		Auth: "",
		Db:   0,
	})

	if err != nil {
		t.Error(err)
		return
	}

	pipe := client.Pipeline()
	cmd := pipe.Set(ctx, "test1", "test", 0)
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(cmd, cmds)

}
