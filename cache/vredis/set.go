package vredis

import (
	"context"
	"github.com/go-redis/redis"
)

// SAdd adds the specified members to the set stored at key.
// Specified members that are already a member of this set are ignored.
// If key does not exist, a new set is created before adding the specified members.
// An error is returned when the value stored at key is not a set.
func (c *RedisClient) SAdd(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	cmd := c.client.SAdd(key, members...)
	c.setSpan(ctx, cmd)
	return cmd
}

// SCard returns the set cardinality (number of elements) of the set stored at key.
func (c *RedisClient) SCard(ctx context.Context, key string) *redis.IntCmd {
	cmd := c.client.SCard(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// SDiff returns the members of the set resulting from the difference between the first set
// and all the successive sets.
func (c *RedisClient) SDiff(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	cmd := c.client.SDiff(keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// SDiffStore is like SDiff, but instead of returning the resulting set, it is stored in destination.
// If destination already exists, it is overwritten.
func (c *RedisClient) SDiffStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	cmd := c.client.SDiffStore(destination, keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// SInter returns the members of the set resulting from the intersection of all the given sets.
func (c *RedisClient) SInter(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	cmd := c.client.SInter(keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// SInterStore is like SInter, but instead of returning the resulting set, it is stored in destination.
// If destination already exists, it is overwritten.
func (c *RedisClient) SInterStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	cmd := c.client.SInterStore(destination, keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// SIsMember returns if member is a member of the set stored at key.
func (c *RedisClient) SIsMember(ctx context.Context, key string, member interface{}) *redis.BoolCmd {
	cmd := c.client.SIsMember(key, member)
	c.setSpan(ctx, cmd)
	return cmd
}

// SMembers returns all the members of the set value stored at key.
func (c *RedisClient) SMembers(ctx context.Context, key string) *redis.StringSliceCmd {
	cmd := c.client.SMembers(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// SMove moves member from the set at source to the set at destination.
// This operation is atomic.
func (c *RedisClient) SMove(ctx context.Context, source, destination, member string) *redis.BoolCmd {
	cmd := c.client.SMove(source, destination, member)
	c.setSpan(ctx, cmd)
	return cmd
}

// SPop removes and returns a random element from the set value stored at key.
func (c *RedisClient) SPop(ctx context.Context, key string) *redis.StringCmd {
	cmd := c.client.SPop(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// SRandMember returns a random element from the set value stored at key.
func (c *RedisClient) SRandMember(ctx context.Context, key string) *redis.StringCmd {
	cmd := c.client.SRandMember(key)
	c.setSpan(ctx, cmd)
	return cmd
}

// SRandMemberN selects random members from a set.
func (c *RedisClient) SRandMemberN(ctx context.Context, key string, count int64) *redis.StringSliceCmd {
	cmd := c.client.SRandMemberN(key, count)
	c.setSpan(ctx, cmd)
	return cmd
}

// SRem removes the specified members from the set stored at key.
// Specified members that are not a member of this set are ignored.
// If key does not exist, it is treated as an empty set and this command returns 0.
func (c *RedisClient) SRem(ctx context.Context, key string, members ...interface{}) *redis.IntCmd {
	cmd := c.client.SRem(key, members...)
	c.setSpan(ctx, cmd)
	return cmd
}

// SUnion returns the members of the set resulting from the union of all the given sets.
func (c *RedisClient) SUnion(ctx context.Context, keys ...string) *redis.StringSliceCmd {
	cmd := c.client.SUnion(keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// SUnionStore is like SUnion, but instead of returning the resulting set, it is stored in destination.
// If destination already exists, it is overwritten.
func (c *RedisClient) SUnionStore(ctx context.Context, destination string, keys ...string) *redis.IntCmd {
	cmd := c.client.SUnionStore(destination, keys...)
	c.setSpan(ctx, cmd)
	return cmd
}

// SScan is like Scan, but with the additional ability to specify the collection
// of keys to scan. SScan is atomic and isolated.
// The Redis documentation says: "Returns all elements of the collection
// between start and end."
// In this case, the "collection" is a set.
func (c *RedisClient) SScan(ctx context.Context, key string, cursor uint64, match string, count int64) *redis.ScanCmd {
	cmd := c.client.SScan(key, cursor, match, count)
	c.setSpan(ctx, cmd)
	return cmd
}
