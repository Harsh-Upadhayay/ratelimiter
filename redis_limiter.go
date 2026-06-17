package ratelimiter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	redis redisAdapter
	algo  redisAlgorithm
}

func NewRedisLimiter(client *redis.Client, algo redisAlgorithm) (*RedisLimiter, error) {
	if client == nil {
		return nil, ErrNilRedisClient
	}

	if algo == nil {
		return nil, ErrNilAlgorithm
	}
	return &RedisLimiter{
		redis: &goRedisAdapter{client: client},
		algo:  algo,
	}, nil
}

func (l *RedisLimiter) Allow(ctx context.Context, key string) (Result, error) {
	if key == "" {
		return Result{}, ErrEmptyKey
	}

	script := l.algo.script()
	args := l.algo.args()

	raw, err := l.redis.eval(ctx, script, []string{key}, args)

	if err != nil {
		return Result{}, err
	}

	values, ok := raw.([]any)
	if !ok {
		return Result{}, ErrUnexpectedRedisResult
	}

	if len(values) != 3 {
		return Result{}, ErrUnexpectedRedisResult
	}

	allowedInt, ok := values[0].(int64)
	if !ok || !(allowedInt == 0 || allowedInt == 1) {
		return Result{}, ErrUnexpectedRedisResult
	}

	remaining, ok := values[1].(int64)
	if !ok {
		return Result{}, ErrUnexpectedRedisResult
	}

	retryAfterSeconds, ok := values[2].(int64)
	if !ok {
		return Result{}, ErrUnexpectedRedisResult
	}

	return Result{
		Allowed:    bool(allowedInt == 1),
		Remaining:  int(remaining),
		RetryAfter: time.Second * time.Duration(retryAfterSeconds),
	}, nil
}
