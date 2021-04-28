package log_writter

import (
	"context"
	"github.com/go-redis/redis/v8"
)

type RedisListWriter struct {
	rdb *redis.Client
	key string
}

// NewRedisListWriter make a redis writer to push every traffic into a redis list
func NewRedisListWriter(rdbLogKey string, rdb *redis.Client) *RedisListWriter {
	return &RedisListWriter{rdb: rdb, key: rdbLogKey}
}

func (r *RedisListWriter) Write(p []byte) (n int, err error) {
	_, err = r.rdb.LPush(context.Background(), r.key, p).Result()
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
