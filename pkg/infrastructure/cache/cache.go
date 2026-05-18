package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrCacheKeyExists = errors.New("cache key already exists")
	ErrCacheKeyEmpty  = errors.New("cache key is empty")
	ErrCacheMiss      = errors.New("cache miss")
)

var (
	instance *Cache
	once     sync.Once
)

type Data struct {
	Key      string
	Value    any
	Duration time.Duration
	Override bool
}

type Cache struct {
	items map[string]item
	mutex sync.RWMutex
	now   func() time.Time
}

type item struct {
	value     any
	expiresAt time.Time
}

func Connect() *Cache {
	once.Do(func() {
		instance = New()
	})
	return instance
}

func New() *Cache {
	return &Cache{
		items: map[string]item{},
		now:   time.Now,
	}
}

func (c *Cache) Store(ctx context.Context, data Data) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if data.Key == "" {
		return ErrCacheKeyEmpty
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !data.Override {
		cached, ok := c.items[data.Key]
		if ok && !cached.expired(c.now()) {
			return ErrCacheKeyExists
		}
	}

	c.items[data.Key] = item{
		value:     data.Value,
		expiresAt: expiresAt(c.now(), data.Duration),
	}
	return nil
}

func (c *Cache) Get(ctx context.Context, key string) (any, error) {
	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	if key == "" {
		return nil, ErrCacheKeyEmpty
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	cached, ok := c.items[key]
	if !ok {
		return nil, ErrCacheMiss
	}
	if cached.expired(c.now()) {
		delete(c.items, key)
		return nil, ErrCacheMiss
	}
	return cached.value, nil
}

func (c *Cache) Delete(ctx context.Context, key string) error {
	if err := checkContext(ctx); err != nil {
		return err
	}
	if key == "" {
		return ErrCacheKeyEmpty
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	delete(c.items, key)
	return nil
}

func (c *Cache) Flush(ctx context.Context) error {
	if err := checkContext(ctx); err != nil {
		return err
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items = map[string]item{}
	return nil
}

func (c *Cache) Len() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.pruneExpired()
	return len(c.items)
}

func (c *Cache) pruneExpired() {
	now := c.now()
	for key, cached := range c.items {
		if cached.expired(now) {
			delete(c.items, key)
		}
	}
}

func (i item) expired(now time.Time) bool {
	return !i.expiresAt.IsZero() && !now.Before(i.expiresAt)
}

func expiresAt(now time.Time, duration time.Duration) time.Time {
	if duration <= 0 {
		return time.Time{}
	}
	return now.Add(duration)
}

func checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
