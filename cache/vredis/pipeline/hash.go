package pipeline

import (
	"context"
	"github.com/go-redis/redis"
)

// HDel deletes the specified fields from the hash stored at key.
// returns the number of fields that were removed from the hash stored at key.
func (p *Pipeline) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd {
	cmd := p.Pipe.HDel(key, fields...)
	p.setSpan(ctx, cmd)
	return cmd
}

// HExists returns if field is an existing field in the hash stored at key.
func (p *Pipeline) HExists(ctx context.Context, key, field string) *redis.BoolCmd {
	cmd := p.Pipe.HExists(key, field)
	p.setSpan(ctx, cmd)
	return cmd
}

// HGet returns the value associated with field in the hash stored at key.
func (p *Pipeline) HGet(ctx context.Context, key, field string) *redis.StringCmd {
	cmd := p.Pipe.HGet(key, field)
	p.setSpan(ctx, cmd)
	return cmd
}

// HGetAll returns all fields and values of the hash stored at key.
func (p *Pipeline) HGetAll(ctx context.Context, key string) *redis.StringStringMapCmd {
	cmd := p.Pipe.HGetAll(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// HIncrBy increments the number stored at field in the hash stored at key by increment.
func (p *Pipeline) HIncrBy(ctx context.Context, key, field string, incr int64) *redis.IntCmd {
	cmd := p.Pipe.HIncrBy(key, field, incr)
	p.setSpan(ctx, cmd)
	return cmd
}

// HIncrByFloat increments the float64 value stored at field in the hash stored at key by increment.
func (p *Pipeline) HIncrByFloat(ctx context.Context, key, field string, incr float64) *redis.FloatCmd {
	cmd := p.Pipe.HIncrByFloat(key, field, incr)
	p.setSpan(ctx, cmd)
	return cmd
}

// HKeys returns all field names in the hash stored at key.
func (p *Pipeline) HKeys(ctx context.Context, key string) *redis.StringSliceCmd {
	cmd := p.Pipe.HKeys(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// HLen returns the number of fields contained in the hash stored at key.
func (p *Pipeline) HLen(ctx context.Context, key string) *redis.IntCmd {
	cmd := p.Pipe.HLen(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// HMGet returns the values associated with the specified fields in the hash stored at key.
func (p *Pipeline) HMGet(ctx context.Context, key string, fields ...string) *redis.SliceCmd {
	cmd := p.Pipe.HMGet(key, fields...)
	p.setSpan(ctx, cmd)
	return cmd
}

// HMSet sets the specified fields to their respective values in the hash stored at key.
func (p *Pipeline) HMSet(ctx context.Context, key string, fields map[string]interface{}) *redis.StatusCmd {
	cmd := p.Pipe.HMSet(key, fields)
	p.setSpan(ctx, cmd)
	return cmd
}

// HSet sets field in the hash stored at key to value.
func (p *Pipeline) HSet(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	cmd := p.Pipe.HSet(key, field, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// HSetNX sets field in the hash stored at key to value, only if field does not yet exist.
func (p *Pipeline) HSetNX(ctx context.Context, key, field string, value interface{}) *redis.BoolCmd {
	cmd := p.Pipe.HSetNX(key, field, value)
	p.setSpan(ctx, cmd)
	return cmd
}

// HVals returns all values in the hash stored at key.
func (p *Pipeline) HVals(ctx context.Context, key string) *redis.StringSliceCmd {
	cmd := p.Pipe.HVals(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// HScan returns all fields and values of the hash stored at key.
func (p *Pipeline) HScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	cmd := p.Pipe.HScan(key, cursor, match, count)
	p.setSpan(ctx, cmd)
	return cmd
}
