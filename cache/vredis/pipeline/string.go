package pipeline

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

// Append appends the value at the end of the string stored at key
// If key does not exist, it is created and set as an empty string,
func (p *Pipeline) Append(ctx context.Context, key, value string) *redis.IntCmd {
	cmd := p.Pipe.Append(key, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// BitCount counts the number of set bits (population counting) in a string.
// By default all the bytes contained in the string are examined. It is possible
// to specify the counting operation only in an interval passing the additional arguments start and end.
func (p *Pipeline) BitCount(ctx context.Context, key string, bitCount *redis.BitCount) *redis.IntCmd {
	cmd := p.Pipe.BitCount(key, bitCount)
	p.setSpan(ctx, cmd)
	return cmd
}

// BitOpAnd performs a bitwise operation between multiple keys (containing string values)
// and stores the result in the destination key.
func (p *Pipeline) BitOpAnd(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	cmd := p.Pipe.BitOpAnd(destKey, keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// BitOpOr performs a bitwise operation between multiple keys (containing string values)
// and stores the result in the destination key.
func (p *Pipeline) BitOpOr(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	cmd := p.Pipe.BitOpOr(destKey, keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// BitOpNot performs a bitwise operation between multiple keys (containing string values)
// and stores the result in the destination key.
func (p *Pipeline) BitOpNot(ctx context.Context, destKey string, key string) *redis.IntCmd {
	cmd := p.Pipe.BitOpNot(destKey, key)
	p.setSpan(ctx, cmd)
	return cmd
}

// BitOpXor performs a bitwise operation between multiple keys (containing string values)
// and stores the result in the destination key.
func (p *Pipeline) BitOpXor(ctx context.Context, destKey string, keys ...string) *redis.IntCmd {
	cmd := p.Pipe.BitOpXor(destKey, keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// Decr Decrements the number stored at key by one.
// If the key does not exist, it is set to 0 before performing the operation
func (p *Pipeline) Decr(ctx context.Context, key string) *redis.IntCmd {
	cmd := p.Pipe.Decr(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// DecrBy Decrements the number stored at key by decrement. If the key does not exist,
// it is set to 0 before performing the operation
func (p *Pipeline) DecrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	cmd := p.Pipe.DecrBy(key, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// Get retrieve a string value associated with the given key.
func (p *Pipeline) Get(ctx context.Context, key string) *redis.StringCmd {
	cmd := p.Pipe.Get(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// GetBit returns the bit value at offset in the string value stored at key.
// The offset is zero-based, and may be negative to index characters from the end of the string.
func (p *Pipeline) GetBit(ctx context.Context, key string, offset int64) *redis.IntCmd {
	cmd := p.Pipe.GetBit(key, offset)
	p.setSpan(ctx, cmd)
	return cmd
}

// GetRange returns the substring of the string value stored at key
// between the offsets start and stop.
func (p *Pipeline) GetRange(ctx context.Context, key string, start, end int64) *redis.StringCmd {
	cmd := p.Pipe.GetRange(key, start, end)
	p.setSpan(ctx, cmd)
	return cmd
}

// GetSet command sets a key to a new value, returning the old value as the result
func (p *Pipeline) GetSet(ctx context.Context, key string, value interface{}) *redis.StringCmd {
	cmd := p.Pipe.GetSet(key, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// Incr The INCR command parses the string value as an integer, increments it by one,
// and finally sets the obtained value as the new value. INCR is atomic
// Increments the number stored at key by one. If the key does not exist,
// it is set to 0 before performing the operation
func (p *Pipeline) Incr(ctx context.Context, key string) *redis.IntCmd {
	cmd := p.Pipe.Incr(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// IncrBy Increments the number stored at key by increment. If the key does not exist,
// it is set to 0 before performing the operation
func (p *Pipeline) IncrBy(ctx context.Context, key string, value int64) *redis.IntCmd {
	cmd := p.Pipe.IncrBy(key, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// IncrByFloat Increments the number stored at key by increment. If the key does not exist,
// it is set to 0 before performing the operation
func (p *Pipeline) IncrByFloat(ctx context.Context, key string, value float64) *redis.FloatCmd {
	cmd := p.Pipe.IncrByFloat(key, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// MGet MSet The ability to set or retrieve the value of multiple keys in
// a single command is also useful for reduced latency.
// For this reason there are the MSET and MGET
// When MGET is used, Redis returns an array of values
func (p *Pipeline) MGet(ctx context.Context, keys ...string) *redis.SliceCmd {
	cmd := p.Pipe.MGet(keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

func (p *Pipeline) MSet(ctx context.Context, keyValues ...interface{}) *redis.StatusCmd {
	cmd := p.Pipe.MSet(keyValues...)
	p.setSpan(ctx, cmd)
	return cmd
}

// MSetNX The MSETNX command is similar to MSET, with the only difference that
// it only sets the keys and values if none of the keys exist.
func (p *Pipeline) MSetNX(ctx context.Context, keyValues ...interface{}) *redis.BoolCmd {
	cmd := p.Pipe.MSetNX(keyValues...)
	p.setSpan(ctx, cmd)
	return cmd
}

// Set that SET will replace any existing value already stored into the key, in the case that the key already exists,
// even if the key is associated with a non-string value.
func (p *Pipeline) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	cmd := p.Pipe.Set(key, value, expiration)
	p.setSpan(ctx, cmd)
	return cmd
}

// SetBit sets the bit at offset in the string value stored at key.
// The offset is zero-based, and may be negative to count from the end of the string.
func (p *Pipeline) SetBit(ctx context.Context, key string, offset int64, value int) *redis.IntCmd {
	cmd := p.Pipe.SetBit(key, offset, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// SetRange sets the part of the string stored at key to value.
// The operation is atomip.
func (p *Pipeline) SetRange(ctx context.Context, key string, offset int64, value string) *redis.IntCmd {
	cmd := p.Pipe.SetRange(key, offset, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// SetNX it only succeed if the key do not exists
func (p *Pipeline) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	cmd := p.Pipe.SetNX(key, value, expiration)
	p.setSpan(ctx, cmd)
	return cmd
}

// SetXX it only succeed if the key already exists
func (p *Pipeline) SetXX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	cmd := p.Pipe.SetXX(key, value, expiration)
	p.setSpan(ctx, cmd)
	return cmd
}

// StrLen returns the length of the string value stored at key.
func (p *Pipeline) StrLen(ctx context.Context, key string) *redis.IntCmd {
	cmd := p.Pipe.StrLen(key)
	p.setSpan(ctx, cmd)
	return cmd
}
