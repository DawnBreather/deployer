package deployment

import (
  "bytes"
  "context"
  "fmt"
  "github.com/sirupsen/logrus"
  "gopkg.in/yaml.v3"
  "time"
)

func getYamlFromRedis(key string) ([]byte, error) {
  val, err := getConnectionToRedis().Get(context.TODO(), key).Result()
  if err != nil {
    logrus.Errorf("[E] Getting value from Redis for { %s } key: %v", key, err)
    return nil, err
  }
  return []byte(val), nil
}

func decodeYamlFromBytes(yamlBytes []byte, out any) error {
  decoder := yaml.NewDecoder(bytes.NewReader(yamlBytes))
  if err := decoder.Decode(out); err != nil {
    return err
  }
  return nil
}

func handleTierConfigurationDecryption(s Secrets, key string, val []byte, out *Tier) error {
  var resultingTierObject Tier
  initialVal := val

  if !IsSecretsDecryptionAccomplished {
    time.Sleep(1 * time.Second) // Consider using a more robust synchronization mechanism
  }

  result, err := UnfoldSecretsPlaceholdersInYaml(s, val)
  if err != nil {
    return err
  }

  val, err = yaml.Marshal(result)
  if err != nil {
    return err
  }

  if err := decodeYamlFromBytes(val, &resultingTierObject); err != nil {
    logrus.Errorf("[E] Unmarhsalling value from Redis key { %s }: %v", key, err)
    return err
  }

  resultingTierObject.initialYamlTemplate = initialVal
  *out = resultingTierObject
  return nil
}

func GetYamlEntryFromRedisWithDecryption(s Secrets, key string, out any, d *Deployment) error {
  val, err := getYamlFromRedis(key)
  if err != nil {
    return err
  }

  if key == d.TierRedisConfigurationEntry(d.tierName()) {
    if tierOut, ok := out.(*Tier); ok {
      return handleTierConfigurationDecryption(s, key, val, tierOut)
    } else {
      return fmt.Errorf("expected output to be of type *Tier, but got %T", out)
    }
  }

  if err := decodeYamlFromBytes(val, out); err != nil {
    logrus.Errorf("[E] Unmarhsalling value from Redis key { %s }: %v", key, err)
    return err
  }

  return nil
}

func GetYamlEntryFromRedisWithoutDecryption(key string, out any) error {
  val, err := getYamlFromRedis(key)
  if err != nil {
    return err
  }

  if err := decodeYamlFromBytes(val, out); err != nil {
    logrus.Errorf("[E] Unmarhsalling value from Redis key { %s }: %v", key, err)
    return err
  }

  return nil
}
