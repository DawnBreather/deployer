package aws

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func (sm *secretsManager) GetYamlAndUnmarshal(key string, out any) (err error) {
	value, ok := sm.GetSecret(key)
	if ok {
		err = yaml.Unmarshal([]byte(value), out)
	}
	if !ok && err != nil {
		logrus.Fatalf("[E] Unmarshalling YAML pulled from Secrets Manager { %s }: %v", key, err)
		return err
	}
	return nil
}
