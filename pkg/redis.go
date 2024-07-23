package pkg

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	ctx    *context.Context
	elog   *log.Logger
	log    *log.Logger
}

func NewClient(addr string, pass string, db int) (*Store, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       db,
	})

	ctx := context.Background()
	_, err := client.Conn().Ping(ctx).Result()
	if err != nil {
		log.Println("Failde to connect to redis")
		log.Fatalln(err)
	}

	store := Store{
		client: client,
		ctx:    &ctx,
		elog:   log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		log:    log.New(os.Stdout, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}

	return &store, nil
}

func (c *Store) Set(key, value string) error {
	err := c.client.Set(*c.ctx, key, value, 0).Err()
	if errors.Is(err, redis.ErrClosed) {
		c.elog.Println(err)
	} else if err != nil {
		c.elog.Println(err)
		return err
	} else {
		return nil
	}
	return errors.New("error setting key")
}

func (c *Store) Get(key string) (string, error) {
	val, err := c.client.Get(*c.ctx, key).Result()

	var message string
	if errors.Is(err, redis.Nil) {
		message = fmt.Sprintf("%s does not exist", key)
	} else if err != nil {
		c.elog.Println(err)
		return "", err
	} else {
		message = fmt.Sprintf("[GET]: %s has value %s", key, val)
	}
	return message, nil
}

func (c *Store) Del(key string) (int64, error) {
	exists, err := c.client.Del(*c.ctx, key).Result()
	if errors.Is(err, redis.TxFailedErr) {
		c.log.Println("Reddis connection closed")
		c.log.Fatalln(err)
	} else if err != nil {
		c.elog.Fatalln(err)
	} else {
		return exists, nil
	}
	return -1, errors.New("error deleting key")
}
