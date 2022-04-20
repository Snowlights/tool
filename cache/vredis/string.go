package vredis

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

// Append appends the value at the end of the string stored at key
// If key does not exist, it is created and set as an empty string,
func (c *RedisClient) Append(ctx context.Context, key, value string) *redis.IntCmd {
	cmd := c.client.Append(key, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// BitCount counts the number of set bits (population counting) in a string.
// By default all the bytes contained in the string are examined. It is possible
// to specify the counting operation only in an interval passing the additional arguments start and end.
func (c *RedisClient) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	cmd := c.client.BitCount(key, bitCount)
	c.setSpan(ctx, cmd)
	return cmd
}

// BitOpAnd performs a bitwise operation between multiple keys (containing string values)
// and stores the result in the destination key.
func (c *RedisClient) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	cmd := c.client.BitOpAnd(destKey, keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// BitOpOr performs a bitwise operation between multiple keys (containing string values)
// and stores the result in the destination key.
func (c *RedisClient) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	cmd := c.client.BitOpOr(destKey, keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// BitOpNot performs a bitwise operation between multiple keys (containing string values)
// and stores the result in the destination key.
func (c *RedisClient) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {
	cmd := c.client.BitOpNot(destKey, key)
	c.setSpan(ctx, cmd)
	return cmd
}

// BitOpXor performs a bitwise operation between multiple keys (containing string values)
// and stores the result in the destination key.
func (c *RedisClient) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	cmd := c.client.BitOpXor(destKey, keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// Decr Decrements the number stored at key by one.
// If the key does not exist, it is set to 0 before performing the operation
func (c *RedisClient) Decr(ctx context.Context, key string) *redis.IntCmd {
	cmd := c.client.Decr(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// DecrBy Decrements the number stored at key by decrement. If the key does not exist,
// it is set to 0 before performing the operation
func (c *RedisClient) DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	cmd := c.client.DecrBy(key, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// Get retrieve a string value associated with the given key.
func (c *RedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := c.client.Get(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// GetBit returns the bit value at offset in the string value stored at key.
// The offset is zero-based, and may be negative to index characters from the end of the string.
func (c *RedisClient) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	cmd := c.client.GetBit(key, offset)
	c.setSpan(ctx, cmd)
	return cmd
}

// GetRange returns the substring of the string value stored at key
// between the offsets start and stop.
func (c *RedisClient) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	cmd := c.client.GetRange(key, start, end)
	c.setSpan(ctx, cmd)
	return cmd
}

// GetSet command sets a key to a new value, returning the old value as the result
func (c *RedisClient) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	cmd := c.client.GetSet(key, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// Incr The INCR command parses the string value as an integer, increments it by one,
// and finally sets the obtained value as the new value. INCR is atomic
// Increments the number stored at key by one. If the key does not exist,
// it is set to 0 before performing the operation
func (c *RedisClient) Incr(ctx context.Context, key string) *redis.IntCmd {
	cmd := c.client.Incr(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// IncrBy Increments the number stored at key by increment. If the key does not exist,
// it is set to 0 before performing the operation
func (c *RedisClient) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	cmd := c.client.IncrBy(key, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// IncrByFloat Increments the number stored at key by increment. If the key does not exist,
// it is set to 0 before performing the operation
func (c *RedisClient) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	cmd := c.client.IncrByFloat(key, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// MGet MSet The ability to set or retrieve the value of multiple keys in
// a single command is also useful for reduced latency.
// For this reason there are the MSET and MGET
// When MGET is used, Redis returns an array of values
func (c *RedisClient) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	cmd := c.client.MGet(keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

func (c *RedisClient) MSet(ctx context.Context, keyValues ...interface{}) *redis.StatusCmd {
	cmd := c.client.MSet(keyValues...)
	c.setSpan(ctx, cmd)
	return cmd
}

// MSetNX The MSETNX command is similar to MSET, with the only difference that
// it only sets the keys and values if none of the keys exist.
func (c *RedisClient) MSetNX(ctx context.Context, keyValues ...interface{}) *redis.BoolCmd {
	cmd := c.client.MSetNX(keyValues...)
	c.setSpan(ctx, cmd)
	return cmd
}

// Set that SET will replace any existing value already stored into the key, in the case that the key already exists,
// even if the key is associated with a non-string value.
func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := c.client.Set(key, value, expiration)
	c.setSpan(ctx, cmd)
	return cmd
}

// SetBit sets the bit at offset in the string value stored at key.
// The offset is zero-based, and may be negative to count from the end of the string.
func (c *RedisClient) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	cmd := c.client.SetBit(key, offset, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// SetRange sets the part of the string stored at key to value.
// The operation is atomic.
func (c *RedisClient) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	cmd := c.client.SetRange(key, offset, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// SetNX it only succeed if the key do not exists
func (c *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	cmd := c.client.SetNX(key, value, expiration)
	c.setSpan(ctx, cmd)
	return cmd
}

// SetXX it only succeed if the key already exists
func (c *RedisClient) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	cmd := c.client.SetXX(key, value, expiration)
	c.setSpan(ctx, cmd)
	return cmd
}

// StrLen returns the length of the string value stored at key.
func (c *RedisClient) StrLen(ctx context.Context, key string) *redis.IntCmd {
	cmd := c.client.StrLen(key)
	c.setSpan(ctx, cmd)
	return cmd
}
