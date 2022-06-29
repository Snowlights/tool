package vredis

import (
	"context"
	"fmt"
	"github.com/Snowlights/tool/cache/vredis/pipeline"
	"github.com/Snowlights/tool/vlog"
	"github.com/Snowlights/tool/vtrace"
	"github.com/go-redis/redis"
	"os"
	"strings"
)

type RedisConfig struct {
	Host string `json:"host" properties:"host"`
	Auth string `json:"auth" properties:"auth"`
	Db   int    `json:"db" properties:"db"`
}

type RedisClient struct {
	config *RedisConfig
	client *redis.Client
}

const (
	RedisHost = "REDIS_HOST"
	RedisAuth = "REDIS_AUTH"
)

func init() {
	host, ok := os.LookupEnv(RedisHost)
	if !ok {
		return
	}
	auth, ok := os.LookupEnv(RedisAuth)
	if !ok {
		return
	}
	config := &RedisConfig{
		Host: host,
		Auth: auth,
	}
	DefaultRedisClient, _ = NewRedisClient(context.Background(), config)
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

func (c *RedisClient) Select(ctx context.Context, db int) *redis.Cmd {
	cmd := c.client.Do("SELECT", db)
	c.setSpan(ctx, cmd)
	return cmd
}

func (c *RedisClient) Echo(ctx context.Context, message interface{}) *redis.Cmd {
	cmd := c.client.Do("ECHO", message)
	c.setSpan(ctx, cmd)
	return cmd
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}

func (c *RedisClient) Pipeline() *pipeline.Pipeline {
	return &pipeline.Pipeline{
		Pipe: c.client.Pipeline(),
	}
}
