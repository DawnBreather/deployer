package deployment

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"time"
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

func PublishMessageToRedisChannel(channel, message string) error {
	return getConnectionToRedis().Publish(context.Background(), channel, message).Err()
}

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

func GetYamlEntryFromRedisWithDecryption(s Secrets, key string, out any, d *Deployment) error {

	var val []byte

	valStr, err := getConnectionToRedis().Get(context.TODO(), key).Result()
	if err != nil {
		logrus.Errorf("[E] Getting value from Redis for { %s } key: %v", key, err)
		return err
	}
	val = []byte(valStr)

	//decoder := yaml.NewDecoder(bytes.NewReader([]byte(val)))
	//err = decoder.Decode(out)

	//fmt.Println(val)

	if key == d.TierRedisConfigurationEntry(d.tierName()) {

		var resultingTierObject Tier
		initialVal := val

		if !IsSecretsDecryptionAccomplished {
			time.Sleep(1 * time.Second)
		}

		result, _ := UnfoldSecretsPlaceholdersInYaml(s, val)
		val, err = yaml.Marshal(result)
		if err != nil {
			return err
		}

		decoder := yaml.NewDecoder(bytes.NewReader(val))
		err = decoder.Decode(&resultingTierObject)
		if err != nil {
			logrus.Errorf("[E] Unmarhsalling value from Redis key { %s }: %v", key, err)
			return err
		}
		resultingTierObject.initialYamlTemplate = initialVal
		*(out.(*Tier)) = resultingTierObject
		return nil

	}

	decoder := yaml.NewDecoder(bytes.NewReader(val))
	err = decoder.Decode(out)
	if err != nil {
		logrus.Errorf("[E] Unmarhsalling value from Redis key { %s }: %v", key, err)
		return err
	}

	return nil
}

func GetYamlEntryFromRedisWithoutDecryption(key string, out any) error {

	var val []byte

	valStr, err := getConnectionToRedis().Get(context.TODO(), key).Result()
	if err != nil {
		logrus.Errorf("[E] Getting value from Redis for { %s } key: %v", key, err)
		return err
	}
	val = []byte(valStr)

	decoder := yaml.NewDecoder(bytes.NewReader(val))
	err = decoder.Decode(out)
	if err != nil {
		logrus.Errorf("[E] Unmarhsalling value from Redis key { %s }: %v", key, err)
		return err
	}

	return nil
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
