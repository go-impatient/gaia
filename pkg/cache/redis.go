package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

var _ Client = &redisClient{}
var ctx = context.Background()

func NewRedisClient(c *redis.Client, encoding Encoding) *redisClient {
	return &redisClient{
		client:   c,
		encoding: encoding,
	}
}

type redisClient struct {
	expireTime time.Duration
	client     *redis.Client
	encoding   Encoding
}

func (c *redisClient) Clear(ctx context.Context) error {
	if c == nil {
		return errors.New("redis client is disable")
	}
	client := c.client.WithContext(ctx)
	err := client.FlushDB(ctx).Err()
	return err
}

func (c *redisClient) Get(ctx context.Context, key string, data interface{}) error {
	if c == nil {
		return errors.New("redis cache is disabled")
	}
	client := c.client.WithContext(ctx)
	cmd := client.Get(ctx, key)
	b, err := cmd.Bytes()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return err
	}

	return c.encoding.Decode(b, data)
}

func (c *redisClient) Set(ctx context.Context, key string, data interface{}, expiration time.Time) error {
	if c == nil {
		return errors.New("redis cache is disabled")
	}

	data, err := c.encoding.Encode(data)
	if err != nil {
		return err
	}

	client := c.client.WithContext(ctx)
	cmd := client.Set(ctx, key, data, TtlForExpiration(expiration))
	return cmd.Err()
}

func (c *redisClient) Add(ctx context.Context, key string, data interface{}, expiration time.Time) error {
	if c == nil {
		return errors.New("redis cache is disabled")
	}

	b, err := c.encoding.Encode(data)
	if err != nil {
		return err
	}
	client := c.client.WithContext(ctx)
	cmd := client.SetNX(ctx, key, b, TtlForExpiration(expiration))
	if !cmd.Val() {
		return ErrNotStored
	}
	return cmd.Err()
}

func (c *redisClient) Delete(ctx context.Context, key string) error {
	if c == nil {
		return errors.New("redis cache is disable")
	}

	client := c.client.WithContext(ctx)
	err := client.Del(ctx, key).Err()
	return err
}

func (c *redisClient) Increment(ctx context.Context, key string, delta uint64) (uint64, error) {
	if c == nil {
		return 0, errors.New("redis cache is disable")
	}

	client := c.client.WithContext(ctx)
	cmd := client.IncrBy(ctx, key, int64(delta))
	val, err := cmd.Result()
	return uint64(val), err
}

func (c *redisClient) Decrement(ctx context.Context, key string, delta uint64) (uint64, error) {
	if c == nil {
		return 0, errors.New("redis cache is disable")
	}

	client := c.client.WithContext(ctx)
	cmd := client.DecrBy(ctx, key, int64(delta))
	val, err := cmd.Result()
	return uint64(val), err
}
