package pipeline

import (
	"context"
	"github.com/go-redis/redis"
)

// SAdd adds the specified members to the set stored at key.
// Specified members that are already a member of this set are ignored.
// If key does not exist, a new set is created before adding the specified members.
// An error is returned when the value stored at key is not a set.
func (p *Pipeline) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	cmd := p.Pipe.SAdd(key, members...)
	p.setSpan(ctx, cmd)
	return cmd
}

// SCard returns the set cardinality (number of elements) of the set stored at key.
func (p *Pipeline) SCard(ctx context.Context, key string) *redis.IntCmd {
	cmd := p.Pipe.SCard(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// SDiff returns the members of the set resulting from the difference between the first set
// and all the successive sets.
func (p *Pipeline) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	cmd := p.Pipe.SDiff(keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// SDiffStore is like SDiff, but instead of returning the resulting set, it is stored in destination.
// If destination already exists, it is overwritten.
func (p *Pipeline) SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	cmd := p.Pipe.SDiffStore(destination, keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// SInter returns the members of the set resulting from the intersection of all the given sets.
func (p *Pipeline) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	cmd := p.Pipe.SInter(keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// SInterStore is like SInter, but instead of returning the resulting set, it is stored in destination.
// If destination already exists, it is overwritten.
func (p *Pipeline) SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	cmd := p.Pipe.SInterStore(destination, keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// SIsMember returns if member is a member of the set stored at key.
func (p *Pipeline) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	cmd := p.Pipe.SIsMember(key, member)
	p.setSpan(ctx, cmd)
	return cmd
}

// SMembers returns all the members of the set value stored at key.
func (p *Pipeline) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	cmd := p.Pipe.SMembers(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// SMove moves member from the set at source to the set at destination.
// This operation is atomip.
func (p *Pipeline) SMove(ctx context.Context, source, destination, member string) *redis.BoolCmd {
	cmd := p.Pipe.SMove(source, destination, member)
	p.setSpan(ctx, cmd)
	return cmd
}

// SPop removes and returns a random element from the set value stored at key.
func (p *Pipeline) SPop(ctx context.Context, key string) *redis.StringCmd {
	cmd := p.Pipe.SPop(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// SRandMember returns a random element from the set value stored at key.
func (p *Pipeline) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	cmd := p.Pipe.SRandMember(key)
	p.setSpan(ctx, cmd)
	return cmd
}

// SRandMemberN selects random members from a set.
func (p *Pipeline) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	cmd := p.Pipe.SRandMemberN(key, count)
	p.setSpan(ctx, cmd)
	return cmd
}

// SRem removes the specified members from the set stored at key.
// Specified members that are not a member of this set are ignored.
// If key does not exist, it is treated as an empty set and this command returns 0.
func (p *Pipeline) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	cmd := p.Pipe.SRem(key, members...)
	p.setSpan(ctx, cmd)
	return cmd
}

// SUnion returns the members of the set resulting from the union of all the given sets.
func (p *Pipeline) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	cmd := p.Pipe.SUnion(keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// SUnionStore is like SUnion, but instead of returning the resulting set, it is stored in destination.
// If destination already exists, it is overwritten.
func (p *Pipeline) SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	cmd := p.Pipe.SUnionStore(destination, keys...)
	p.setSpan(ctx, cmd)
	return cmd
}

// SScan is like Scan, but with the additional ability to specify the collection
// of keys to scan. SScan is atomic and isolated.
// The Redis documentation says: "Returns all elements of the collection
// between start and end."
// In this case, the "collection" is a set.
func (p *Pipeline) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	cmd := p.Pipe.SScan(key, cursor, match, count)
	p.setSpan(ctx, cmd)
	return cmd
}
