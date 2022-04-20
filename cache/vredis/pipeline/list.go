package pipeline

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
func (p *Pipeline) BLPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	cmd := p.Pipe.BLPop(timeout, keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

func (p *Pipeline) BRPop(ctx context.Context, timeout time.Duration, keys ...string) *redis.StringSliceCmd {
	cmd := p.Pipe.BRPop(timeout, keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// BRPopLPush BRPOPLPUSH is the blocking variant of RPOPLPUSH.
// When source contains elements, this command behaves exactly like RPOPLPUSH.
// When source is empty, Redis will block the connection until another client pushes to it or
// until the timeout is reached.
func (p *Pipeline) BRPopLPush(ctx context.Context, source, destination string, timeout time.Duration) *redis.StringCmd {
	cmd := p.Pipe.BRPopLPush(source, destination, timeout)
	p.setSpan(ctx, cmd)
	return cmd
}

// LIndex GET an element from a list by its index.
func (p *Pipeline) LIndex(ctx context.Context, key string, index int64) *redis.StringCmd {
	cmd := p.Pipe.LIndex(key, index)
	p.setSpan(ctx, cmd)
	return cmd
}

// LInsert LIST Insert an element before or after another element in a list.
func (p *Pipeline) LInsert(ctx context.Context, key, op string, pivot, value interface{}) *redis.IntCmd {
	cmd := p.Pipe.LInsert(key, op, pivot, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// LLen Returns the length of the list stored at key. If key does not exist,
// it is interpreted as an empty list and 0 is returned.
// An error is returned when the value stored at key is not a list
func (p *Pipeline) LLen(ctx context.Context, key string) *redis.IntCmd {
	cmd := p.Pipe.LLen(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// LPop Removes and returns the front elements of the list stored at key
func (p *Pipeline) LPop(ctx context.Context, key string) *redis.StringCmd {
	cmd := p.Pipe.LPop(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// LPush Insert all the specified values at the head of the list stored at key.
// If key does not exist, it is created as empty list before performing the push operations.
// When key holds a value that is not a list, an error is returned
func (p *Pipeline) LPush(ctx context.Context, key string, value ...interface{}) *redis.IntCmd {
	cmd := p.Pipe.LPush(key, value...)
	p.setSpan(ctx, cmd)
	return cmd
}

// LPushX Inserts value at the head of the list stored at key, only if key already exists and holds a list.
// In contrary to LPUSH, no operation will be performed when key does not yet exist.
func (p *Pipeline) LPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	cmd := p.Pipe.LPushX(key, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// LRange Return the specified elements of the list stored at key.
// The offsets start and stop are zero-based indexes,
// with 0 being the first element of the list (the head of the list),
// 1 being the next element and so on.
func (p *Pipeline) LRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	cmd := p.Pipe.LRange(key, start, stop)
	p.setSpan(ctx, cmd)
	return cmd
}

// LRem Remove elements from a list
func (p *Pipeline) LRem(ctx context.Context, key string, count int64, value interface{}) *redis.IntCmd {
	cmd := p.Pipe.LRem(key, count, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// LSet Assign a value to a given list index
func (p *Pipeline) LSet(ctx context.Context, key string, index int64, value interface{}) *redis.StatusCmd {
	cmd := p.Pipe.LSet(key, index, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// LTrim Trim an existing list so that it will contain only the specified range of elements specified.
// Both start and stop are zero-based indexes, where 0 is the first element of the list (the head),
// 1 the next element and so on.
func (p *Pipeline) LTrim(ctx context.Context, key string, start, stop int64) *redis.StatusCmd {
	cmd := p.Pipe.LTrim(key, start, stop)
	p.setSpan(ctx, cmd)
	return cmd
}

// RPop Removes and returns the last elements of the list stored at key
func (p *Pipeline) RPop(ctx context.Context, key string) *redis.StringCmd {
	cmd := p.Pipe.RPop(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// RPopLPush RPOPLPUSH is equivalent to RPOP followed by LPUSH.
func (p *Pipeline) RPopLPush(ctx context.Context, source, destination string) *redis.StringCmd {
	cmd := p.Pipe.RPopLPush(source, destination)
	p.setSpan(ctx, cmd)
	return cmd
}

// RPush Insert all the specified values at the tail of the list stored at key.
// If key does not exist, it is created as empty list before performing the push operations.
// When key holds a value that is not a list, an error is returned
func (p *Pipeline) RPush(ctx context.Context, key string, value ...interface{}) *redis.IntCmd {
	cmd := p.Pipe.RPush(key, value...)
	p.setSpan(ctx, cmd)
	return cmd
}

// RPushX Inserts value at the tail of the list stored at key, only if key already exists and holds a list.
func (p *Pipeline) RPushX(ctx context.Context, key string, value interface{}) *redis.IntCmd {
	cmd := p.Pipe.RPushX(key, value)
	p.setSpan(ctx, cmd)
	return cmd
}
