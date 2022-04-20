package vredis

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

// BLPop BLPOP is a blocking list pop primitive.
// It is the blocking version of LPOP because it blocks the connection when
// there are no elements to pop from any of the given lists.
// values
// A nil multi-bulk when no element could be popped and the timeout expired.
// A two-element multi-bulk with the first element being the name of the key where an element
// was popped and the second element being the value of the popped element
func (c *RedisClient) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	cmd := c.client.BLPop(timeout, keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

func (c *RedisClient) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	cmd := c.client.BRPop(timeout, keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// BRPopLPush BRPOPLPUSH is the blocking variant of RPOPLPUSH.
// When source contains elements, this command behaves exactly like RPOPLPUSH.
// When source is empty, Redis will block the connection until another client pushes to it or
// until the timeout is reached.
func (c *RedisClient) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	cmd := c.client.BRPopLPush(source, destination, timeout)
	c.setSpan(ctx, cmd)
	return cmd
}

// LIndex GET an element from a list by its index.
func (c *RedisClient) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	cmd := c.client.LIndex(key, index)
	c.setSpan(ctx, cmd)
	return cmd
}

// LInsert LIST Insert an element before or after another element in a list.
func (c *RedisClient) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	cmd := c.client.LInsert(key, op, pivot, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// LLen Returns the length of the list stored at key. If key does not exist,
// it is interpreted as an empty list and 0 is returned.
// An error is returned when the value stored at key is not a list
func (c *RedisClient) LLen(ctx context.Context, key string) *redis.IntCmd {
	cmd := c.client.LLen(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// LPop Removes and returns the front elements of the list stored at key
func (c *RedisClient) LPop(ctx context.Context, key string) *redis.StringCmd {
	cmd := c.client.LPop(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// LPush Insert all the specified values at the head of the list stored at key.
// If key does not exist, it is created as empty list before performing the push operations.
// When key holds a value that is not a list, an error is returned
func (c *RedisClient) LPush(ctx context.Context, key string, value ...interface{}) *redis.IntCmd {
	cmd := c.client.LPush(key, value...)
	c.setSpan(ctx, cmd)
	return cmd
}

// LPushX Inserts value at the head of the list stored at key, only if key already exists and holds a list.
// In contrary to LPUSH, no operation will be performed when key does not yet exist.
func (c *RedisClient) LPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	cmd := c.client.LPushX(key, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// LRange Return the specified elements of the list stored at key.
// The offsets start and stop are zero-based indexes,
// with 0 being the first element of the list (the head of the list),
// 1 being the next element and so on.
func (c *RedisClient) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	cmd := c.client.LRange(key, start, stop)
	c.setSpan(ctx, cmd)
	return cmd
}

// LRem Remove elements from a list
func (c *RedisClient) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	cmd := c.client.LRem(key, count, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// LSet Assign a value to a given list index
func (c *RedisClient) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	cmd := c.client.LSet(key, index, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// LTrim Trim an existing list so that it will contain only the specified range of elements specified.
// Both start and stop are zero-based indexes, where 0 is the first element of the list (the head),
// 1 the next element and so on.
func (c *RedisClient) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	cmd := c.client.LTrim(key, start, stop)
	c.setSpan(ctx, cmd)
	return cmd
}

// RPop Removes and returns the last elements of the list stored at key
func (c *RedisClient) RPop(ctx context.Context, key string) *redis.StringCmd {
	cmd := c.client.RPop(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// RPopLPush RPOPLPUSH is equivalent to RPOP followed by LPUSH.
func (c *RedisClient) RPopLPush(ctx context.Context, source, destination string) *redis.StringCmd {
	cmd := c.client.RPopLPush(source, destination)
	c.setSpan(ctx, cmd)
	return cmd
}

// RPush Insert all the specified values at the tail of the list stored at key.
// If key does not exist, it is created as empty list before performing the push operations.
// When key holds a value that is not a list, an error is returned
func (c *RedisClient) RPush(ctx context.Context, key string, value ...interface{}) *redis.IntCmd {
	cmd := c.client.RPush(key, value...)
	c.setSpan(ctx, cmd)
	return cmd
}

// RPushX Inserts value at the tail of the list stored at key, only if key already exists and holds a list.
func (c *RedisClient) RPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	cmd := c.client.RPushX(key, value)
	c.setSpan(ctx, cmd)
	return cmd
}
