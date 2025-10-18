package repository

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"strconv"
)

const luaScript = `
local current = redis.call('GET', KEYS[1])
if current == false then
    current = 0
else
    current = tonumber(current)
end

current = current + 1
redis.call('SET', KEYS[1], current, 'EX', ARGV[1])
return current
`

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo() RepoInterface {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		DB:       0,
		Password: "",
		Protocol: 2,
	})
	return &RedisRepo{
		client: client,
	}
}

func (r RedisRepo) Increment(ctx context.Context, key string, expiry int) (int, error) {
	cmd := r.client.Eval(ctx, luaScript, []string{key}, expiry)
	result, err := cmd.Int()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (r RedisRepo) GetByKey(ctx context.Context, key string) (int, error) {
	resultRaw, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	result, err := strconv.ParseInt(resultRaw, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(result), nil
}
