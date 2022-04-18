package vredis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"vtool/vlog"
	"vtool/vtrace"
)

type RedisConfig struct {
	Host string
	Auth string
	Db   int
}

type RedisClient struct {
	config *RedisConfig
	client *redis.Client
}

var DefaultRedisClient *RedisClient

func NewRedisClient(ctx context.Context, cfg *RedisConfig) (*RedisClient, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Password: cfg.Auth,
		DB:       cfg.Db,
	})

	pong, err := c.Ping().Result()
	if err != nil {
		vlog.FatalF(ctx, "init redis failed, error is %s ", err.Error())
		return nil, err
	} else {
		vlog.InfoF(ctx, "init redis success, pong is %s", pong)
	}

	cli := &RedisClient{
		config: cfg,
		client: c,
	}
	DefaultRedisClient = cli
	return DefaultRedisClient, nil
}

func (c *RedisClient) setSpan(ctx context.Context, cmd redis.Cmder) {
	if span := vtrace.SpanFromContent(ctx); span != nil {
		span.SetTag(vtrace.Component, vtrace.ComponentRedis)
		span.SetTag(vtrace.SpanKind, vtrace.SpanKindRedis)
		span.SetTag(vtrace.RedisCluster, c.config.Host)
		span.SetTag(vtrace.RedisCmd, strings.Join(
			func() []string {
				var cmds []string
				for _, arg := range cmd.Args() {
					cmds = append(cmds, fmt.Sprint(arg))
				}
				return cmds
			}(), " "))
	}
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}
