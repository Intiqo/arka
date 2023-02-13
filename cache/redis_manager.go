package cache

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"

	"github.com/adwitiyaio/arka/constants"
	"github.com/adwitiyaio/arka/logger"
	"github.com/adwitiyaio/arka/secrets"
)

const hostKey = "REDIS_HOST"
const portKey = "REDIS_PORT"
const passwordKey = "REDIS_PASSWORD"
const databaseKey = "REDIS_DATABASE"

type redisCacheManager struct {
	sm     secrets.Manager
	client *redis.Client
}

func (r *redisCacheManager) initialize() {
	host := r.sm.GetValueForKey(hostKey)
	port := r.sm.GetValueForKey(portKey)
	pass := r.sm.GetValueForKey(passwordKey)
	db, err := strconv.Atoi(r.sm.GetValueForKey(databaseKey))
	if err != nil {
		db = 1
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	opts := &redis.Options{
		Addr: addr,
		DB:   db,
	}
	if os.Getenv("CI") != "true" {
		opts.Password = pass
	}
	r.client = redis.NewClient(opts)

	_, err = r.client.Ping().Result()
	if err != nil {
		logger.Log.Fatal().Err(err).Stack().Msg("unable to connect to redis")
	}
}

func (r *redisCacheManager) GetStatus() string {
	status, err := r.client.Ping().Result()
	if err != nil || !strings.EqualFold(status, "pong") {
		logger.Log.Error().Err(err).Msg("cache connection may be down")
		return constants.SystemStatusDown
	}
	return constants.SystemStatusUp
}

func (r *redisCacheManager) Set(key string, val interface{}) error {
	return r.client.Set(key, val, 0).Err()
}

func (r *redisCacheManager) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}
