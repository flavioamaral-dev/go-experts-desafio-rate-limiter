package storage

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func InitRedis(addr string) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected to Redis")

	return &Redis{
		client: client,
	}, nil
}

func (a *Redis) AddAccess(ctx context.Context, keyType string, key string, maxAccesses int64) (bool, int64, error) {
	now := time.Now()
	clearBefore := now.Add(-time.Second)
	pipeline := a.client.Pipeline()

	redisKey := a.customRedisKey("access", keyType, key)
	pipeline.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(clearBefore.UnixMicro(), 10))
	count := pipeline.ZCard(ctx, redisKey)

	_, err := pipeline.Exec(ctx)
	if err != nil {
		fmt.Println("Error on pipeline exec", err)
		return false, 0, err
	}

	if count.Val() >= maxAccesses {
		return false, count.Val(), nil
	}

	pipeline = a.client.Pipeline()
	pipeline.ZAdd(ctx, redisKey, redis.Z{Score: float64(now.UnixMicro()), Member: now.Format(time.RFC3339Nano)})
	pipeline.Expire(ctx, redisKey, time.Second)

	_, err = pipeline.Exec(ctx)
	if err != nil {
		fmt.Println("Error on pipeline exec", err)
		return false, 0, err
	}

	return true, count.Val() + 1, nil
}

func (a *Redis) GetBlock(ctx context.Context, keyType string, key string) (*time.Time, error) {
	redisKey := a.customRedisKey("block", keyType, key)
	blockTime, err := a.client.Get(ctx, redisKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	blockTimeInt, err := strconv.ParseInt(blockTime, 10, 64)
	if err != nil {
		fmt.Println("Error parsing block time", err)
		return nil, err
	}

	blockTimeTime := time.Unix(0, blockTimeInt)
	return &blockTimeTime, nil
}

func (a *Redis) AddBlock(ctx context.Context, keyType string, key string, blockTimeMilliseconds int64) (*time.Time, error) {
	redisKey := a.customRedisKey("block", keyType, key)
	blockTime := time.Now().Add(time.Duration(blockTimeMilliseconds) * time.Millisecond)
	err := a.client.Set(ctx, redisKey, blockTime.UnixNano(), time.Duration(blockTimeMilliseconds)*time.Millisecond).Err()
	if err != nil {
		fmt.Println("Error setting block", err)
		return nil, err
	}

	return &blockTime, nil
}

func (s *Redis) customRedisKey(prefix string, keyType string, key string) string {
	return fmt.Sprintf(
		"%s-%s-%s",
		strings.ToLower(prefix),
		strings.ToLower(strings.ReplaceAll(keyType, "-", "_")),
		key,
	)
}
