package vredis

import (
	"context"
	"github.com/go-redis/redis"
)

func (c *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	cmd := c.client.SAdd(key, members...)
	c.setSpan(ctx, cmd)
	return cmd
}
