package ratelimiter

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type goRedisAdapter struct {
	client *redis.Client
}

func (a *goRedisAdapter) eval(ctx context.Context, script string, keys []string, args []string) (any, error) {
	redisArgs := make([]any, 0, len(args)) // Current size is 0, but we preallocate len(args) space to it.
	for _, arg := range args {
		redisArgs = append(redisArgs, arg)
	}

	return a.client.Eval(ctx, script, keys, redisArgs...).Result()
}
