package ratelimiter

import "context"

type RedisLimter struct {
	redis redisAdapter
	algo  redisAlgorithm
}

func (l *RedisLimter) Allow(ctx context.Context, key string) (Result, error) {
	return Result{}, nil
}
