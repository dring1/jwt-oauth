package Cache

import "sync"

type Cache struct {
	data map[string]string
	sync.Mutex
}

func (c *Cache) Insert(key string, value string) {
	c.Lock()
	c.data[key] = value
	c.Unlock()
}

func (c *Cache) Remove(key string) {
	c.Lock()
	delete(c.data, key)
	c.Unlock()
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}
