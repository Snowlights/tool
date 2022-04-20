package pipeline

import (
	"context"
	"github.com/go-redis/redis"
)

func (p *Pipeline) ZAdd(ctx context.Context, key string, zList ...redis.Z) *redis.IntCmd {
	cmd := p.Pipe.ZAdd(key, zList...)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZCard gets the number of elements in the sorted set at key.
func (p *Pipeline) ZCard(ctx context.Context, key string) *redis.IntCmd {
	cmd := p.Pipe.ZCard(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZCount gets the number of elements in the sorted set at key with a score between min and max.
func (p *Pipeline) ZCount(ctx context.Context, key string, min, max string) *redis.IntCmd {
	cmd := p.Pipe.ZCount(key, min, max)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZIncrBy increments the score of member in the sorted set stored at key by increment.
func (p *Pipeline) ZIncrBy(ctx context.Context, key string, increment float64, member string) *redis.FloatCmd {
	cmd := p.Pipe.ZIncrBy(key, increment, member)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZRange gets the specified range of elements in the sorted set stored at key.
func (p *Pipeline) ZRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	cmd := p.Pipe.ZRange(key, start, stop)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZRangeWithScores gets the specified range of elements in the sorted set stored at key.
func (p *Pipeline) ZRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	cmd := p.Pipe.ZRangeWithScores(key, start, stop)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZRank gets the rank of member in the sorted set stored at key.
func (p *Pipeline) ZRank(ctx context.Context, key, member string) *redis.IntCmd {
	cmd := p.Pipe.ZRank(key, member)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZRem removes the specified members from the sorted set stored at key.
func (p *Pipeline) ZRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	cmd := p.Pipe.ZRem(key, members...)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZRemRangeByRank removes all elements in the sorted set stored at key with rank between start and stop.
func (p *Pipeline) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *redis.IntCmd {
	cmd := p.Pipe.ZRemRangeByRank(key, start, stop)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZRemRangeByScore removes all elements in the sorted set stored at key with a score between min and max.
func (p *Pipeline) ZRemRangeByScore(ctx context.Context, key, min, max string) *redis.IntCmd {
	cmd := p.Pipe.ZRemRangeByScore(key, min, max)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZRevRange gets the specified range of elements in the sorted set stored at key.
func (p *Pipeline) ZRevRange(ctx context.Context, key string, start, stop int64) *redis.StringSliceCmd {
	cmd := p.Pipe.ZRevRange(key, start, stop)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZRevRangeWithScores gets the specified range of elements in the sorted set stored at key.
func (p *Pipeline) ZRevRangeWithScores(ctx context.Context, key string, start, stop int64) *redis.ZSliceCmd {
	cmd := p.Pipe.ZRevRangeWithScores(key, start, stop)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZRevRank gets the rank of member in the sorted set stored at key.
func (p *Pipeline) ZRevRank(ctx context.Context, key, member string) *redis.IntCmd {
	cmd := p.Pipe.ZRevRank(key, member)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZScore gets the score of member in the sorted set at key.
func (p *Pipeline) ZScore(ctx context.Context, key, member string) *redis.FloatCmd {
	cmd := p.Pipe.ZScore(key, member)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZUnionStore is a sorted set union.
func (p *Pipeline) ZUnionStore(ctx context.Context, dest string, store redis.ZStore, keys ...string) *redis.IntCmd {
	cmd := p.Pipe.ZUnionStore(dest, store, keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZInterStore is a sorted set intersection.
func (p *Pipeline) ZInterStore(ctx context.Context, dest string, store redis.ZStore, keys ...string) *redis.IntCmd {
	cmd := p.Pipe.ZInterStore(dest, store, keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// ZScan is a sorted set scan.
func (p *Pipeline) ZScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	cmd := p.Pipe.ZScan(key, cursor, match, count)
	p.setSpan(ctx, cmd)
	return cmd
}
