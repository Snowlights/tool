package vredis

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

// Del Removes the specified keys. A key is ignored if it does not exist.
func (c *RedisClient) Del(ctx context.Context, key ...string) *redis.IntCmd {
	cmd := c.client.Del(key...)
	c.setSpan(ctx, cmd)
	return cmd
}

// Dump Returns the serialized form of the value stored at the specified key.
// If key does not exist a nil value is returned.
func (c *RedisClient) Dump(ctx context.Context, key string) *redis.StringCmd {
	cmd := c.client.Dump(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// Exists Returns if key exists.
//  it will be counted multiple times. So if somekey exists, EXISTS somekey somekey will return 2
func (c *RedisClient) Exists(ctx context.Context, key string) *redis.IntCmd {
	cmd := c.client.Exists(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// Expire
// The EXPIRE command supports a set of options:
// Options
// NX -- Set expiry only when the key has no expiry
// XX -- Set expiry only when the key has an existing expiry
// GT -- Set expiry only when the new expiry is greater than current one
// LT -- Set expiry only when the new expiry is less than current one
// A non-volatile key is treated as an infinite TTL for the purpose of
// GT and LT. The GT, LT and NX options are mutually exclusive
// Return Value
// 1: if the timeout was set.
// 0: if the timeout was not set. e.g. key doesn't exist, or operation skipped due to the provided arguments.
func (c *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	cmd := c.client.Expire(key, expiration)
	c.setSpan(ctx, cmd)
	return cmd
}

// ExpireAt Set the expiration for a key as a UNIX timestamp.
func (c *RedisClient) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	cmd := c.client.ExpireAt(key, tm)
	c.setSpan(ctx, cmd)
	return cmd
}

// Keys Returns all keys matching pattern.
// If no pattern is specified, the command returns all keys.
// The command sorts the keys alphabetically.
func (c *RedisClient) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	cmd := c.client.Keys(pattern)
	c.setSpan(ctx, cmd)
	return cmd
}

// Migrate Moves a key to another database.
// If the source database is not specified,
// the key is always copied to the destination database.
func (c *RedisClient) Migrate(ctx context.Context, host, port, key string, db int64, timeout time.Duration) *redis.StatusCmd {
	cmd := c.client.Migrate(host, port, key, db, timeout)
	c.setSpan(ctx, cmd)
	return cmd
}

// Move Moves a key to another database.
// If the source database is not specified,
// the key is always copied to the destination database.
func (c *RedisClient) Move(ctx context.Context, key string, db int64) *redis.BoolCmd {
	cmd := c.client.Move(key, db)
	c.setSpan(ctx, cmd)
	return cmd
}

// ObjectEncoding Returns the object encoding for the specified key.
// Object encoding is the internal representation of the value stored at a key.
func (c *RedisClient) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd {
	cmd := c.client.ObjectEncoding(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// ObjectIdleTime Returns the idle time information for the specified key.
// The idle time is the number of seconds since the last access to the key.
// The access is either a read or write operation.
func (c *RedisClient) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd {
	cmd := c.client.ObjectIdleTime(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// ObjectRefCount Returns the number of references of the value associated with the specified key.
// The command returns -1 when the key does not exist.
func (c *RedisClient) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd {
	cmd := c.client.ObjectRefCount(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// Persist Remove the existing timeout on key,
// turning the key from volatile (a key with an expire set) to persistent
// (a key that will never expire as no timeout is associated
// value:
// 1: if the timeout was removed.
// 0: if the key does not exist or does not have an associated timeout.
func (c *RedisClient) Persist(ctx context.Context, key string) *redis.BoolCmd {
	cmd := c.client.Persist(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// PExpire Set a timeout on key.
// After the timeout has expired, the key will automatically be deleted.
// A key with an associated timeout cannot be changed as the timeout
func (c *RedisClient) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	cmd := c.client.PExpire(key, expiration)
	c.setSpan(ctx, cmd)
	return cmd
}

// PExpireAt Set the expiration for a key as a UNIX timestamp.
// After the specified timestamp, the key will be deleted.
// A key with an associated timeout cannot be changed as the timeout
func (c *RedisClient) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	cmd := c.client.PExpireAt(key, tm)
	c.setSpan(ctx, cmd)
	return cmd
}

// PTTL Like TTL this command returns the remaining time to live of a key that has an expire set,
// with the sole difference that TTL returns the amount of remaining time in seconds
// while PTTL returns it in milliseconds
// value:
// The command returns -2 if the key does not exist.
// The command returns -1 if the key exists but has no associated expire.
func (c *RedisClient) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	cmd := c.client.PTTL(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// RandomKey Returns a random key from the currently selected database.
// The command returns a nil when the database is empty.
func (c *RedisClient) RandomKey(ctx context.Context) *redis.StringCmd {
	cmd := c.client.RandomKey()
	c.setSpan(ctx, cmd)
	return cmd
}

// Rename key to newkey.
// It returns an error when key does not exist.
// It does not rename newkey if it already exists.
func (c *RedisClient) Rename(ctx context.Context, key, newkey string) *redis.StatusCmd {
	cmd := c.client.Rename(key, newkey)
	c.setSpan(ctx, cmd)
	return cmd
}

// RenameNX Rename key to newkey if newkey does not yet exist.
// It returns an error when key does not exist.
func (c *RedisClient) RenameNX(ctx context.Context, key, newkey string) *redis.BoolCmd {
	cmd := c.client.RenameNX(key, newkey)
	c.setSpan(ctx, cmd)
	return cmd
}

// Restore key from a ttl encoded string value stored at source.
// It returns an error when the key already exists.
func (c *RedisClient) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	cmd := c.client.Restore(key, ttl, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// RestoreReplace key from a ttl encoded string value stored at source.
// It returns an error when the key already exists.
func (c *RedisClient) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	cmd := c.client.RestoreReplace(key, ttl, value)
	c.setSpan(ctx, cmd)
	return cmd
}

// Sort sorts the elements in a list, set or sorted set.
// By default sorting is numeric with elements being compared as double precision floating point numbers.
// It is possible to change the comparison algorithm by adding an optional parameter:
func (c *RedisClient) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	cmd := c.client.Sort(key, sort)
	c.setSpan(ctx, cmd)
	return cmd
}

// TTL Returns the remaining time to live of a key that has a timeout. This introspection capability
// allows a Redis client to check how many seconds a given key will continue to be part of the dataset
// value:
// The command returns -2 if the key does not exist.
// The command returns -1 if the key exists but has no associated expire.
func (c *RedisClient) TTL(ctx context.Context, key string) *redis.DurationCmd {
	cmd := c.client.TTL(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// Type returns the string representation of the type of the value stored at key.
// The different types that can be returned are: string, list, set, zset and hash.
func (c *RedisClient) Type(ctx context.Context, key string) *redis.StatusCmd {
	cmd := c.client.Type(key)
	c.setSpan(ctx, cmd)
	return cmd
}

func (c *RedisClient) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	cmd := c.client.Scan(cursor, match, count)
	c.setSpan(ctx, cmd)
	return cmd
}
