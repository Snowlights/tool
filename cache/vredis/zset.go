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

// ZCard gets the number of elements in the sorted set at key.
func (c *RedisClient) ZCard(ctx context.Context, key string) *redis.IntCmd {
	cmd := c.client.ZCard(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZCount gets the number of elements in the sorted set at key with a score between min and max.
func (c *RedisClient) ZCount(ctx context.Context, key string, min, max string) *redis.IntCmd {
	cmd := c.client.ZCount(key, min, max)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZIncrBy increments the score of member in the sorted set stored at key by increment.
func (c *RedisClient) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	cmd := c.client.ZIncrBy(key, increment, member)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZRange gets the specified range of elements in the sorted set stored at key.
func (c *RedisClient) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	cmd := c.client.ZRange(key, start, stop)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZRangeWithScores gets the specified range of elements in the sorted set stored at key.
func (c *RedisClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	cmd := c.client.ZRangeWithScores(key, start, stop)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZRank gets the rank of member in the sorted set stored at key.
func (c *RedisClient) ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	cmd := c.client.ZRank(key, member)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZRem removes the specified members from the sorted set stored at key.
func (c *RedisClient) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	cmd := c.client.ZRem(key, members...)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZRemRangeByRank removes all elements in the sorted set stored at key with rank between start and stop.
func (c *RedisClient) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	cmd := c.client.ZRemRangeByRank(key, start, stop)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZRemRangeByScore removes all elements in the sorted set stored at key with a score between min and max.
func (c *RedisClient) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	cmd := c.client.ZRemRangeByScore(key, min, max)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZRevRange gets the specified range of elements in the sorted set stored at key.
func (c *RedisClient) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	cmd := c.client.ZRevRange(key, start, stop)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZRevRangeWithScores gets the specified range of elements in the sorted set stored at key.
func (c *RedisClient) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	cmd := c.client.ZRevRangeWithScores(key, start, stop)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZRevRank gets the rank of member in the sorted set stored at key.
func (c *RedisClient) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	cmd := c.client.ZRevRank(key, member)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZScore gets the score of member in the sorted set at key.
func (c *RedisClient) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	cmd := c.client.ZScore(key, member)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZUnionStore is a sorted set union.
func (c *RedisClient) ZUnionStore(ctx context.Context, dest string, store redis.ZStore, keys ...string) *redis.IntCmd {
	cmd := c.client.ZUnionStore(dest, store, keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZInterStore is a sorted set intersection.
func (c *RedisClient) ZInterStore(ctx context.Context, dest string, store redis.ZStore, keys ...string) *redis.IntCmd {
	cmd := c.client.ZInterStore(dest, store, keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// ZScan is a sorted set scan.
func (c *RedisClient) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	cmd := c.client.ZScan(key, cursor, match, count)
	c.setSpan(ctx, cmd)
	return cmd
}
