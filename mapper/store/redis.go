package store

import (
	"context"
	"fmt"
	"mapper/config"
	"net"
	"net/url"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

var updateStreamScript = redis.NewScript(`
	local key = KEYS[1]
	local value = ARGV[1]
	local ttl = tonumber(ARGV[2])

	local current = redis.call('GET', key)
	if current == false or current ~= value then
		redis.call('SET', key, value, 'EX', ttl)
	else
		redis.call('EXPIRE', key, ttl)
	end

	return "OK"
`)

func NewRedisStore(cfg *config.Config) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.RedisAddr,
		DialTimeout:  1 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStore{client: client}, nil
}

func (s *RedisStore) UpdateStream(ctx context.Context, stream, host string, ttl int) error {
	key := fmt.Sprintf("streams:%s", stream)
	host, err := hostFromURL(host)
	if err != nil {
		return fmt.Errorf("failed to get host from URL: %w", err)
	}

	if updateStreamScript.Run(ctx, s.client, []string{key}, host, ttl).Err(); err != nil {
		return fmt.Errorf("failed to execute update stream script for %s: %w", stream, err)
	}

	return nil
}

func (s *RedisStore) UpdateStreams(ctx context.Context, streams map[string][]string, ttl int) error {
	for host, streamList := range streams {
		for _, stream := range streamList {
			if err := s.UpdateStream(ctx, stream, host, ttl); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *RedisStore) GetStream(ctx context.Context, term string) (string, error) {
	key := fmt.Sprintf("streams:%s", term)
	return s.client.Get(ctx, key).Result()
}

func hostFromURL(u string) (string, error) {
	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	host, _, err := net.SplitHostPort(parsedURL.Host)

	return host, err
}
