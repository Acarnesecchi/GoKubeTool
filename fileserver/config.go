package fileserver

import (
	"errors"
	"time"
)

type Config struct {
	storage      string
	storageAddr  string
	listenAddr   string
	queue        string
	queueAddr    string
	volumeType   string
	databaseAddr string
	timeout      time.Duration
}

func NewConfig() Config {
	return Config{
		storage:     "local",
		storageAddr: "/opt/devtools/storage/",
		listenAddr:  ":50058",
		queue:       "rabbitmq",
		queueAddr:   ":50059",
		volumeType:  "host-storage",
		timeout:     10 * time.Second,
	}
}

func (c Config) WithStorage(s string, addr ...string) (Config, error) {
	switch s {
	case "local":
		c.storage = s
		return c, nil
	case "s3":
		c.storage = s
		if addr[0] != "" {
			c.storageAddr = addr[0]
			return c, nil
		}
		return c, errors.New("empty storage address")
	case "bucket":
		c.storage = s
		if addr[0] != "" {
			c.storageAddr = addr[0]
			return c, nil
		}
		return c, errors.New("empty storage address")
	default:
		return c, errors.New("unknown storage option")
	}
}

func (c Config) WithListenAddr(s string) Config {
	c.listenAddr = s
	return c
}

func (c Config) WithQueueAddr(s string) Config {
	c.queueAddr = s
	return c
}

func (c Config) WithVolume(s string) (Config, error) {
	switch s {
	case "host-storage":
		c.volumeType = s
	case "persistentVolume":
		c.volumeType = s
	default:
		return c, errors.New("unknown volume option")
	}
	return c, nil
}

func (c Config) WithDatabaseAddr(s string) Config {
	c.databaseAddr = s
	return c
}
