package rds

import (
	"fmt"
	"mon-dell-me4012/config"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/go-redis/redis"
)

var R *redis.Client = &redis.Client{}

func RedisInit(c *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host + ":" + strconv.Itoa(c.Redis.Port),
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Redis connection error: %s", err)
	}

	return client
}

func RedisSet(r *redis.Client, k, v string, e int) error {
	err := r.Set(k, v, time.Duration(e)*time.Second).Err()
	if err != nil {
		return fmt.Errorf("redis SET error: %s", err)
	}

	return nil
}

func RedisGet(r *redis.Client, k string) (string, error) {
	v, err := r.Get(k).Result()
	if err != nil {
		return "", nil
	}

	return v, nil
}
