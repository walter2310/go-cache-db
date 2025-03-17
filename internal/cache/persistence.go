package cache

import (
	"encoding/json"
	"os"
)

func (c *Cache) SaveToDisk(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	return encoder.Encode(c.data)
}

func (c *Cache) LoadFromDisk(filename string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		c.data = make(map[string]Cacheitem)
		return nil
	}

	return json.Unmarshal(data, &c.data)
}
