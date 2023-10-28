package aws

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func (sm *secretsManager) GetJsonAndUnmarshal(key string, out any) (err error) {
	value, ok := sm.GetSecret(key)
	if ok {
		err = json.Unmarshal([]byte(value), out)
	}
	if !ok && err != nil {
		logrus.Fatalf("[E] Unmarshalling JSON pulled from Secrets Manager { %s }: %v", key, err)
		return err
	}
	return nil
}
