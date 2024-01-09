package deployment

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func getConnectionToRedis() *redis.Client {
  if redisClient == nil {
    redisClient = redis.NewClient(&redis.Options{
      Addr:     fmt.Sprintf("%s:%s", REDIS_CREDENTIALS.Host, REDIS_CREDENTIALS.Port),
      Password: REDIS_CREDENTIALS.Password,
      TLSConfig: &tls.Config{
        InsecureSkipVerify: true,
      },
    })
  }

  return redisClient
}

func GetRedisKeysByMask(mask string) ([]string, error) {
  var cursor uint64
  var keys []string
  for {
    var (
      keysInCurrentIteration []string
      err                    error
    )
    keysInCurrentIteration, cursor, err = getConnectionToRedis().Scan(context.Background(), cursor, mask, 1000).Result()
    if err != nil {
      return nil, err
    }
    keys = append(keys, keysInCurrentIteration...)
    if cursor == 0 {
      break
    }
  }
  return keys, nil
}
