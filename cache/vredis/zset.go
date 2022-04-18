package vredis

import (
	"context"
	"github.com/go-redis/redis"
)

func (c *RedisClient) ZAdd(ctx context.Context, key string, zList ...redis.Z) *redis.IntCmd {
	cmd := c.client.ZAdd(key, zList...)
	c.setSpan(ctx, cmd)
	return cmd
}
