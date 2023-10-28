package deployment

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"time"
)

func (d *Deployment) StartListeningForSecretsConfigurationFromRedisChannel() {
	go func() {
		subscribeForConfigurationFromRedisChannel(d.SecretsRedisChannel(), func(message string) {
			for PauseConfigurationListening {
				time.Sleep(1 * time.Second)
			}

			decoder := yaml.NewDecoder(bytes.NewReader([]byte(message)))
			err := decoder.Decode(&d.Secrets)
			if err != nil {
				logrus.Errorf("[E] Unmarshalling secrets configuration { %s } into { Secrets } struct: %v", CONFIGURATION_STORAGE_SECRETS_PATH, err)
			}
			IsSecretsDecryptionAccomplished = false
			d.Secrets.DecryptEncryptedValues()
		})
	}()
}

func (d *Deployment) StartListeningForSecretsConfigurationFromGit() {
	go func() {
		for {

			for PauseConfigurationListening {
				time.Sleep(1 * time.Second)
			}

			if hasFileChanged(CONFIGURATION_STORAGE_SECRETS_PATH) {
				err := unmarshalYamlToStruct(CONFIGURATION_STORAGE_SECRETS_PATH, &d.Secrets)
				if err != nil {
					logrus.Errorf("[E] Unmarshalling secrets configuration { %s } into { Secrets } struct: %v", CONFIGURATION_STORAGE_SECRETS_PATH, err)
				}
				IsSecretsDecryptionAccomplished = false
				d.Secrets.DecryptEncryptedValues()
			}
			time.Sleep(time.Duration(CONFIGURATION_STORAGE_GIT_PULL_PERIOD) * time.Second)
		}
	}()
}
