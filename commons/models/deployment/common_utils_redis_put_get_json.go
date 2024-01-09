package deployment

import (
  "context"
  "encoding/json"
  "github.com/sirupsen/logrus"
  "time"
)

func PutIntoRedis(key string, object any, ttl time.Duration) error {
  data, err := json.Marshal(object)
  if err != nil {
    logrus.Errorf("[E] Marhalling object into JSON for Redis: %v", err)
    return err
  }

  err = getConnectionToRedis().Set(context.Background(), key, data, ttl).Err()
  if err != nil {
    logrus.Errorf("[E] Saving object into Redis under { %s } key: %v", key, err)
    return err
  }
  return nil
}

func GetJsonEntryFromRedis(key string, out any) error {
  val, err := getConnectionToRedis().Get(context.TODO(), key).Result()
  if err != nil {
    logrus.Errorf("[E] Getting value from Redis for { %s } key: %v", key, err)
    return err
  }

  err = json.Unmarshal([]byte(val), out)
  if err != nil {
    logrus.Errorf("[E] Unmarhsalling value from Redis key { %s }: %v", key, err)
    return err
  }
  return nil
}
