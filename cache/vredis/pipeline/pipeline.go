package pipeline

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"github.com/Snowlights/tool/vtrace"
)

type Pipeline struct {
	Pipe redis.Pipeliner
	Host string
}

func (p *Pipeline) setSpan(ctx context.Context, cmd ...redis.Cmder) {
	if span := vtrace.SpanFromContent(ctx); span != nil {
		span.SetTag(vtrace.Component, vtrace.ComponentRedis)
		span.SetTag(vtrace.SpanKind, vtrace.SpanKindRedis)
		span.SetTag(vtrace.RedisCluster, p.Host)
		span.SetTag(vtrace.WithPipeline, true)
		span.SetTag(vtrace.RedisCmd, strings.Join(
			func() []string {
				var cmds []string
				for _, c := range cmd {
					for _, arg := range c.Args() {
						cmds = append(cmds, fmt.Sprint(arg))
					}
				}
				return cmds
			}(), " "))
	}
}

// Discard discards all the queued commands in a pipeline.
func (p *Pipeline) Discard(ctx context.Context) error {
	return p.Pipe.Discard()
}

// Exec executes all the queued commands in a pipeline.
func (p *Pipeline) Exec(ctx context.Context) ([]redis.Cmder, error) {
	cmd, err := p.Pipe.Exec()
	p.setSpan(ctx, cmd...)
	return cmd, err
}

// Close closes the pipeline.
func (p *Pipeline) Close(ctx context.Context) error {
	return p.Pipe.Close()
}
