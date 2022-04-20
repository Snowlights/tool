package vredis

import (
	"context"
	"github.com/go-redis/redis"
)

// HDel deletes the specified fields from the hash stored at key.
// returns the number of fields that were removed from the hash stored at key.
func (c *RedisClient) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	cmd := c.client.HDel(key, fields...)
	c.setSpan(ctx, cmd)
	return cmd
}

// HExists returns if field is an existing field in the hash stored at key.
func (c *RedisClient) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	cmd := c.client.HExists(key, field)
	c.setSpan(ctx, cmd)
	return cmd
}

// HGet returns the value associated with field in the hash stored at key.
func (c *RedisClient) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	cmd := c.client.HGet(key, field)
	c.setSpan(ctx, cmd)
	return cmd
}

// HGetAll returns all fields and values of the hash stored at key.
func (c *RedisClient) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	cmd := c.client.HGetAll(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// HIncrBy increments the number stored at field in the hash stored at key by increment.
func (c *RedisClient) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	cmd := c.client.HIncrBy(key, field, incr)
	c.setSpan(ctx, cmd)
	return cmd
}

// HIncrByFloat increments the float64 value stored at field in the hash stored at key by increment.
func (c *RedisClient) HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	cmd := c.client.HIncrByFloat(key, field, incr)
	c.setSpan(ctx, cmd)
	return cmd
}

// HKeys returns all field names in the hash stored at key.
func (c *RedisClient) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	cmd := c.client.HKeys(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// HLen returns the number of fields contained in the hash stored at key.
func (c *RedisClient) HLen(ctx context.Context, key string) *redis.IntCmd {
	cmd := c.client.HLen(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// HMGet returns the values associated with the specified fields in the hash stored at key.
func (c *RedisClient) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	cmd := c.client.HMGet(key, fields...)
	c.setSpan(ctx, cmd)
	return cmd
}

// HMSet sets the specified fields to their respective values in the hash stored at key.
func (c *RedisClient) HMSet(ctx context.Context, key string, fields map[string]interface{}) *redis.StatusCmd {
	cmd := c.client.HMSet(key, fields)
	c.setSpan(ctx, cmd)
	return cmd
}

// HSet sets field in the hash stored at key to value.
func (c *RedisClient) HSet(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	cmd := c.client.HSet(key, field, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// HSetNX sets field in the hash stored at key to value, only if field does not yet exist.
func (c *RedisClient) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	cmd := c.client.HSetNX(key, field, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// HVals returns all values in the hash stored at key.
func (c *RedisClient) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	cmd := c.client.HVals(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// HScan returns all fields and values of the hash stored at key.
func (c *RedisClient) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	cmd := c.client.HScan(key, cursor, match, count)
	c.setSpan(ctx, cmd)
	return cmd
}
