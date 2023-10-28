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
			PauseConfigurationListening = true
			IsSecretsDecryptionAccomplished = false

			decoder := yaml.NewDecoder(bytes.NewReader([]byte(message)))
			err := decoder.Decode(&d.Secrets)
			if err != nil {
				logrus.Errorf("[E] Unmarshalling secrets configuration { %s } into { Secrets } struct: %v", CONFIGURATION_STORAGE_SECRETS_PATH, err)
			}
			d.Secrets.DecryptEncryptedValues()

			tierInterfaceObject, err := UnfoldSecretsPlaceholdersInYaml(d.Secrets, d.Tiers[d.tierName()].initialYamlTemplate)
			if err != nil {
				logrus.Errorf("[E] Unfolding secrets for { Tier }: %v", err)
			} else {

				tierResultingYaml, err := yaml.Marshal(tierInterfaceObject)

				decoder = yaml.NewDecoder(bytes.NewReader(tierResultingYaml))
				err = decoder.Decode(d.Tiers[d.tierName()])
				if err != nil {
					logrus.Errorf("[E] Unmarshalling tier configuration { %s } into { Tier } struct: %v", d.TierRedisChannel(d.tierName()), err)
				}

			}

			PauseConfigurationListening = false
		})
	}()
}
