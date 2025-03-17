package cache

import (
	"strings"
	"sync"
	"time"
)

type Cacheitem struct {
	Value      string    `json:"value"`      // Etiquetas JSON para los campos
	Expiration time.Time `json:"expiration"` // Etiquetas JSON para los campos
}

type Cache struct {
	data map[string]Cacheitem
	mu   sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]Cacheitem),
	}
}

// ttl significa Time To Live
func (c *Cache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Now().Add(ttl)

	c.data[key] = Cacheitem{
		Value:      value,
		Expiration: expiration,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	// Si el tiempo actual es mayor que el tiempo de expiración, eliminar la clave
	if time.Now().After(item.Expiration) {
		delete(c.data, key)
		return nil, false
	}

	return item.Value, true
}

func (c *Cache) CleanUpExpiredKeys(interval time.Duration) {
	go func() {
		for {
			time.Sleep(interval)

			c.mu.Lock()
			now := time.Now()
			for key, item := range c.data {
				if now.After(item.Expiration) {
					delete(c.data, key)
				}
			}
			c.mu.Unlock()
		}
	}()
}

func (c *Cache) ListKeys(pattern string) []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	var keys []string
	for key := range c.data {
		if strings.HasPrefix(key, strings.TrimSuffix(pattern, "*")) {
			keys = append(keys, key)
		}
	}

	return keys
}
