package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService provides Redis-based caching functionality
type CacheService struct {
	client     *redis.Client
	defaultTTL time.Duration
	keyPrefix  string
}

// CacheConfig contains Redis cache configuration
type CacheConfig struct {
	Host       string
	Port       int
	Password   string
	DB         int
	DefaultTTL time.Duration
	KeyPrefix  string
}

// NewCacheService creates a new cache service instance
func NewCacheService(config CacheConfig) *CacheService {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	// Test the connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	return &CacheService{
		client:     client,
		defaultTTL: config.DefaultTTL,
		keyPrefix:  config.KeyPrefix,
	}
}

// buildKey creates a prefixed cache key
func (s *CacheService) buildKey(key string) string {
	if s.keyPrefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", s.keyPrefix, key)
}

// Set stores a value in cache with optional TTL
func (s *CacheService) Set(ctx context.Context, key string, value interface{}, ttl ...time.Duration) error {
	cacheKey := s.buildKey(key)

	// Determine TTL
	cacheTTL := s.defaultTTL
	if len(ttl) > 0 {
		cacheTTL = ttl[0]
	}

	// Serialize value to JSON
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return s.client.Set(ctx, cacheKey, data, cacheTTL).Err()
}

// Get retrieves a value from cache
func (s *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	cacheKey := s.buildKey(key)

	data, err := s.client.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get cache value: %w", err)
	}

	// Deserialize from JSON
	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return nil
}

// Delete removes a value from cache
func (s *CacheService) Delete(ctx context.Context, key string) error {
	cacheKey := s.buildKey(key)
	return s.client.Del(ctx, cacheKey).Err()
}

// Exists checks if a key exists in cache
func (s *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	cacheKey := s.buildKey(key)
	count, err := s.client.Exists(ctx, cacheKey).Result()
	return count > 0, err
}

// SetTTL updates the TTL of an existing key
func (s *CacheService) SetTTL(ctx context.Context, key string, ttl time.Duration) error {
	cacheKey := s.buildKey(key)
	return s.client.Expire(ctx, cacheKey, ttl).Err()
}

// GetTTL returns the remaining TTL of a key
func (s *CacheService) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	cacheKey := s.buildKey(key)
	return s.client.TTL(ctx, cacheKey).Result()
}

// SetMulti stores multiple key-value pairs
func (s *CacheService) SetMulti(ctx context.Context, items map[string]interface{}, ttl ...time.Duration) error {
	pipe := s.client.Pipeline()

	// Determine TTL
	cacheTTL := s.defaultTTL
	if len(ttl) > 0 {
		cacheTTL = ttl[0]
	}

	for key, value := range items {
		cacheKey := s.buildKey(key)
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
		}
		pipe.Set(ctx, cacheKey, data, cacheTTL)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetMulti retrieves multiple values from cache
func (s *CacheService) GetMulti(ctx context.Context, keys []string) (map[string]interface{}, error) {
	if len(keys) == 0 {
		return make(map[string]interface{}), nil
	}

	// Build cache keys
	cacheKeys := make([]string, len(keys))
	for i, key := range keys {
		cacheKeys[i] = s.buildKey(key)
	}

	// Get all values
	values, err := s.client.MGet(ctx, cacheKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get multiple cache values: %w", err)
	}

	// Build result map
	result := make(map[string]interface{})
	for i, value := range values {
		if value != nil {
			var data interface{}
			if err := json.Unmarshal([]byte(value.(string)), &data); err != nil {
				continue // Skip invalid JSON
			}
			result[keys[i]] = data
		}
	}

	return result, nil
}

// DeletePattern deletes all keys matching a pattern
func (s *CacheService) DeletePattern(ctx context.Context, pattern string) (int64, error) {
	cachePattern := s.buildKey(pattern)

	// Get all keys matching the pattern
	keys, err := s.client.Keys(ctx, cachePattern).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get keys for pattern: %w", err)
	}

	if len(keys) == 0 {
		return 0, nil
	}

	// Delete all matching keys
	deleted, err := s.client.Del(ctx, keys...).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to delete keys: %w", err)
	}

	return deleted, nil
}

// Increment atomically increments a counter
func (s *CacheService) Increment(ctx context.Context, key string, delta int64) (int64, error) {
	cacheKey := s.buildKey(key)
	return s.client.IncrBy(ctx, cacheKey, delta).Result()
}

// Decrement atomically decrements a counter
func (s *CacheService) Decrement(ctx context.Context, key string, delta int64) (int64, error) {
	cacheKey := s.buildKey(key)
	return s.client.DecrBy(ctx, cacheKey, delta).Result()
}

// SetNX sets a key only if it doesn't exist (atomic lock)
func (s *CacheService) SetNX(ctx context.Context, key string, value interface{}, ttl time.Duration) (bool, error) {
	cacheKey := s.buildKey(key)

	data, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	return s.client.SetNX(ctx, cacheKey, data, ttl).Result()
}

// GetSet atomically sets a new value and returns the old value
func (s *CacheService) GetSet(ctx context.Context, key string, value interface{}) (string, error) {
	cacheKey := s.buildKey(key)

	data, err := json.Marshal(value)
	if err != nil {
		return "", fmt.Errorf("failed to marshal value: %w", err)
	}

	return s.client.GetSet(ctx, cacheKey, data).Result()
}

// ListPush adds elements to the beginning of a list
func (s *CacheService) ListPush(ctx context.Context, key string, values ...interface{}) error {
	cacheKey := s.buildKey(key)

	// Serialize all values
	serializedValues := make([]interface{}, len(values))
	for i, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value at index %d: %w", i, err)
		}
		serializedValues[i] = data
	}

	return s.client.LPush(ctx, cacheKey, serializedValues...).Err()
}

// ListPop removes and returns the first element of a list
func (s *CacheService) ListPop(ctx context.Context, key string, dest interface{}) error {
	cacheKey := s.buildKey(key)

	data, err := s.client.LPop(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to pop from list: %w", err)
	}

	// Deserialize from JSON
	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal list value: %w", err)
	}

	return nil
}

// ListLength returns the length of a list
func (s *CacheService) ListLength(ctx context.Context, key string) (int64, error) {
	cacheKey := s.buildKey(key)
	return s.client.LLen(ctx, cacheKey).Result()
}

// SetAdd adds elements to a set
func (s *CacheService) SetAdd(ctx context.Context, key string, values ...interface{}) error {
	cacheKey := s.buildKey(key)

	// Serialize all values
	serializedValues := make([]interface{}, len(values))
	for i, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value at index %d: %w", i, err)
		}
		serializedValues[i] = data
	}

	return s.client.SAdd(ctx, cacheKey, serializedValues...).Err()
}

// SetMembers returns all members of a set
func (s *CacheService) SetMembers(ctx context.Context, key string) ([]string, error) {
	cacheKey := s.buildKey(key)
	return s.client.SMembers(ctx, cacheKey).Result()
}

// SetRemove removes elements from a set
func (s *CacheService) SetRemove(ctx context.Context, key string, values ...interface{}) error {
	cacheKey := s.buildKey(key)

	// Serialize all values
	serializedValues := make([]interface{}, len(values))
	for i, value := range values {
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value at index %d: %w", i, err)
		}
		serializedValues[i] = data
	}

	return s.client.SRem(ctx, cacheKey, serializedValues...).Err()
}

// HashSet sets a field in a hash
func (s *CacheService) HashSet(ctx context.Context, key, field string, value interface{}) error {
	cacheKey := s.buildKey(key)

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return s.client.HSet(ctx, cacheKey, field, data).Err()
}

// HashGet gets a field from a hash
func (s *CacheService) HashGet(ctx context.Context, key, field string, dest interface{}) error {
	cacheKey := s.buildKey(key)

	data, err := s.client.HGet(ctx, cacheKey, field).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get hash field: %w", err)
	}

	// Deserialize from JSON
	if err := json.Unmarshal([]byte(data), dest); err != nil {
		return fmt.Errorf("failed to unmarshal hash value: %w", err)
	}

	return nil
}

// HashGetAll gets all fields from a hash
func (s *CacheService) HashGetAll(ctx context.Context, key string) (map[string]string, error) {
	cacheKey := s.buildKey(key)
	return s.client.HGetAll(ctx, cacheKey).Result()
}

// HashDelete deletes fields from a hash
func (s *CacheService) HashDelete(ctx context.Context, key string, fields ...string) error {
	cacheKey := s.buildKey(key)
	return s.client.HDel(ctx, cacheKey, fields...).Err()
}

// FlushAll clears all cache entries (use with caution)
func (s *CacheService) FlushAll(ctx context.Context) error {
	return s.client.FlushAll(ctx).Err()
}

// GetStats returns cache statistics
func (s *CacheService) GetStats(ctx context.Context) (map[string]interface{}, error) {
	info, err := s.client.Info(ctx, "memory", "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get Redis info: %w", err)
	}

	// Get database size
	dbSize, err := s.client.DBSize(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get database size: %w", err)
	}

	stats := map[string]interface{}{
		"redis_info": info,
		"db_size":    dbSize,
	}

	return stats, nil
}

// Close closes the Redis connection
func (s *CacheService) Close() error {
	return s.client.Close()
}

// CacheError represents cache-related errors
type CacheError struct {
	Operation string
	Key       string
	Err       error
}

func (e CacheError) Error() string {
	return fmt.Sprintf("cache %s failed for key '%s': %v", e.Operation, e.Key, e.Err)
}

// ErrCacheMiss indicates a cache miss
var ErrCacheMiss = fmt.Errorf("cache miss")

// WithCache is a helper function to implement cache-aside pattern
func (s *CacheService) WithCache(ctx context.Context, key string, ttl time.Duration, fn func() (interface{}, error), dest interface{}) error {
	// Try to get from cache first
	err := s.Get(ctx, key, dest)
	if err == nil {
		return nil // Cache hit
	}

	if err != ErrCacheMiss {
		// Log cache error but continue with function execution
		// In production, you might want to use a proper logger here
	}

	// Cache miss or error, execute the function
	result, err := fn()
	if err != nil {
		return err
	}

	// Store result in cache (fire and forget)
	go func() {
		bgCtx := context.Background()
		s.Set(bgCtx, key, result, ttl)
	}()

	// Copy result to destination
	data, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("failed to marshal result: %w", err)
	}

	return json.Unmarshal(data, dest)
}
