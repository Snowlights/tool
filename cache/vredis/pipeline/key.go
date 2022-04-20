package pipeline

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

// Del Removes the specified keys. A key is ignored if it does not exist.
func (p *Pipeline) Del(ctx context.Context, key ...string) *redis.IntCmd {
	cmd := p.Pipe.Del(key...)
	p.setSpan(ctx, cmd)
	return cmd
}

// Dump Returns the serialized form of the value stored at the specified key.
// If key does not exist a nil value is returned.
func (p *Pipeline) Dump(ctx context.Context, key string) *redis.StringCmd {
	cmd := p.Pipe.Dump(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// Exists Returns if key exists.
//  it will be counted multiple times. So if somekey exists, EXISTS somekey somekey will return 2
func (p *Pipeline) Exists(ctx context.Context, key string) *redis.IntCmd {
	cmd := p.Pipe.Exists(key)
	p.setSpan(ctx, cmd)
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
func (p *Pipeline) Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	cmd := p.Pipe.Expire(key, expiration)
	p.setSpan(ctx, cmd)
	return cmd
}

// ExpireAt Set the expiration for a key as a UNIX timestamp.
func (p *Pipeline) ExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	cmd := p.Pipe.ExpireAt(key, tm)
	p.setSpan(ctx, cmd)
	return cmd
}

// Keys Returns all keys matching pattern.
// If no pattern is specified, the command returns all keys.
// The command sorts the keys alphabetically.
func (p *Pipeline) Keys(ctx context.Context, pattern string) *redis.StringSliceCmd {
	cmd := p.Pipe.Keys(pattern)
	p.setSpan(ctx, cmd)
	return cmd
}

// Migrate Moves a key to another database.
// If the source database is not specified,
// the key is always copied to the destination database.
func (p *Pipeline) Migrate(ctx context.Context, host, port, key string, db int64, timeout time.Duration) *redis.StatusCmd {
	cmd := p.Pipe.Migrate(host, port, key, db, timeout)
	p.setSpan(ctx, cmd)
	return cmd
}

// Move Moves a key to another database.
// If the source database is not specified,
// the key is always copied to the destination database.
func (p *Pipeline) Move(ctx context.Context, key string, db int64) *redis.BoolCmd {
	cmd := p.Pipe.Move(key, db)
	p.setSpan(ctx, cmd)
	return cmd
}

// ObjectEncoding Returns the object encoding for the specified key.
// Object encoding is the internal representation of the value stored at a key.
func (p *Pipeline) ObjectEncoding(ctx context.Context, key string) *redis.StringCmd {
	cmd := p.Pipe.ObjectEncoding(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// ObjectIdleTime Returns the idle time information for the specified key.
// The idle time is the number of seconds since the last access to the key.
// The access is either a read or write operation.
func (p *Pipeline) ObjectIdleTime(ctx context.Context, key string) *redis.DurationCmd {
	cmd := p.Pipe.ObjectIdleTime(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// ObjectRefCount Returns the number of references of the value associated with the specified key.
// The command returns -1 when the key does not exist.
func (p *Pipeline) ObjectRefCount(ctx context.Context, key string) *redis.IntCmd {
	cmd := p.Pipe.ObjectRefCount(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// Persist Remove the existing timeout on key,
// turning the key from volatile (a key with an expire set) to persistent
// (a key that will never expire as no timeout is associated
// value:
// 1: if the timeout was removed.
// 0: if the key does not exist or does not have an associated timeout.
func (p *Pipeline) Persist(ctx context.Context, key string) *redis.BoolCmd {
	cmd := p.Pipe.Persist(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// PExpire Set a timeout on key.
// After the timeout has expired, the key will automatically be deleted.
// A key with an associated timeout cannot be changed as the timeout
func (p *Pipeline) PExpire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd {
	cmd := p.Pipe.PExpire(key, expiration)
	p.setSpan(ctx, cmd)
	return cmd
}

// PExpireAt Set the expiration for a key as a UNIX timestamp.
// After the specified timestamp, the key will be deleted.
// A key with an associated timeout cannot be changed as the timeout
func (p *Pipeline) PExpireAt(ctx context.Context, key string, tm time.Time) *redis.BoolCmd {
	cmd := p.Pipe.PExpireAt(key, tm)
	p.setSpan(ctx, cmd)
	return cmd
}

// PTTL Like TTL this command returns the remaining time to live of a key that has an expire set,
// with the sole difference that TTL returns the amount of remaining time in seconds
// while PTTL returns it in milliseconds
// value:
// The command returns -2 if the key does not exist.
// The command returns -1 if the key exists but has no associated expire.
func (p *Pipeline) PTTL(ctx context.Context, key string) *redis.DurationCmd {
	cmd := p.Pipe.PTTL(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// RandomKey Returns a random key from the currently selected database.
// The command returns a nil when the database is empty.
func (p *Pipeline) RandomKey(ctx context.Context) *redis.StringCmd {
	cmd := p.Pipe.RandomKey()
	p.setSpan(ctx, cmd)
	return cmd
}

// Rename key to newkey.
// It returns an error when key does not exist.
// It does not rename newkey if it already exists.
func (p *Pipeline) Rename(ctx context.Context, key, newkey string) *redis.StatusCmd {
	cmd := p.Pipe.Rename(key, newkey)
	p.setSpan(ctx, cmd)
	return cmd
}

// RenameNX Rename key to newkey if newkey does not yet exist.
// It returns an error when key does not exist.
func (p *Pipeline) RenameNX(ctx context.Context, key, newkey string) *redis.BoolCmd {
	cmd := p.Pipe.RenameNX(key, newkey)
	p.setSpan(ctx, cmd)
	return cmd
}

// Restore key from a ttl encoded string value stored at source.
// It returns an error when the key already exists.
func (p *Pipeline) Restore(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	cmd := p.Pipe.Restore(key, ttl, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// RestoreReplace key from a ttl encoded string value stored at source.
// It returns an error when the key already exists.
func (p *Pipeline) RestoreReplace(ctx context.Context, key string, ttl time.Duration, value string) *redis.StatusCmd {
	cmd := p.Pipe.RestoreReplace(key, ttl, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// Sort sorts the elements in a list, set or sorted set.
// By default sorting is numeric with elements being compared as double precision floating point numbers.
// It is possible to change the comparison algorithm by adding an optional parameter:
func (p *Pipeline) Sort(ctx context.Context, key string, sort *redis.Sort) *redis.StringSliceCmd {
	cmd := p.Pipe.Sort(key, sort)
	p.setSpan(ctx, cmd)
	return cmd
}

// TTL Returns the remaining time to live of a key that has a timeout. This introspection capability
// allows a Redis client to check how many seconds a given key will continue to be part of the dataset
// value:
// The command returns -2 if the key does not exist.
// The command returns -1 if the key exists but has no associated expire.
func (p *Pipeline) TTL(ctx context.Context, key string) *redis.DurationCmd {
	cmd := p.Pipe.TTL(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// Type returns the string representation of the type of the value stored at key.
// The different types that can be returned are: string, list, set, zset and hash.
func (p *Pipeline) Type(ctx context.Context, key string) *redis.StatusCmd {
	cmd := p.Pipe.Type(key)
	p.setSpan(ctx, cmd)
	return cmd
}

func (p *Pipeline) Scan(ctx context.Context, cursor uint64, match string, count int64) *redis.ScanCmd {
	cmd := p.Pipe.Scan(cursor, match, count)
	p.setSpan(ctx, cmd)
	return cmd
}
