package pkg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	ctx    *context.Context
	log    *log.Logger
}

func NewClient(addr string, pass string, db int) (*Store, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	ctx := context.Background()

	store := Store{
		client: client,
		ctx:    &ctx,
		log:    log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return &store, nil
}

func (c *Store) Set(key, value string) {
	err := c.client.Set(*c.ctx, key, value, 0).Err()
	if err != nil {
		c.log.Fatalln(err)
	}
}

func (c *Store) Get(key string) string {
	val, err := c.client.Get(*c.ctx, key).Result()

	var message string

	if errors.Is(err, redis.Nil) {
		message = fmt.Sprintf("%s does not exist", key)
	} else {
		message = fmt.Sprintf("[GET]: %s has value %s", key, val)
	}

	time.Sleep(2 * time.Second)

	return message
}

func (c *Store) Del(key string) int64 {
	exists, err := c.client.Del(*c.ctx, key).Result()
	if err != nil {
		c.log.Fatalln(err)
	}

	return exists
}
