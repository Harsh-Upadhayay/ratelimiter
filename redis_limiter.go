package ratelimiter

import "context"

type RedisLimiter struct {
	redis redisAdapter
	algo  redisAlgorithm
}

func (l *RedisLimiter) Allow(ctx context.Context, key string) (Result, error) {

	return Result{}, nil
}
