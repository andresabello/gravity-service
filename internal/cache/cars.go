package cache

import (
	"sync"
)

type MakeModelYear struct {
	Make   string   `json:"make"`
	Models []string `json:"models"`
	Years  []string `json:"years"`
}

type ModelMakeYear struct {
	Model string   `json:"model"`
	Makes []string `json:"makes"`
	Years []string `json:"years"`
}

type Cache struct {
	memo map[string]interface{}
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		memo: make(map[string]interface{}),
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result, exists := c.memo[key]
	return result, exists
}

func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.memo[key] = value
}
