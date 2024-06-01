package config

import (
	"log"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	env map[string]string
	mu  sync.RWMutex
}

var instance *Config
var once sync.Once

func Instance() *Config {
	once.Do(func() {
		instance = &Config{}
		instance.loadEnvFile(".env")
	})
	return instance
}

func (c *Config) GetEnvVar(key, def string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if val, ok := c.env[key]; ok {
		return val
	}
	return def
}

func (c *Config) loadEnvFile(envFile string) {
	var err error
	c.env, err = godotenv.Read(envFile)

	if err != nil {
		log.Println("Failed to load .env file. Fallback to default values")
		c.env = make(map[string]string)
	}
}
