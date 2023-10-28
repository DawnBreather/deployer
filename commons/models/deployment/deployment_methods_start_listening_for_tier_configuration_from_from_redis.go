package deployment

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"time"
)

func (d *Deployment) StartListeningForTierConfigurationFromRedisChannel() {
	d.InitializeTiersFields(d.tierName())
	go func() {
		subscribeForConfigurationFromRedisChannel(d.TierRedisChannel(d.tierName()), func(message string) {
			for PauseConfigurationListening {
				time.Sleep(1 * time.Second)
			}
			PauseConfigurationListening = true

			initialYaml := []byte(message)

			interfaceObject, err := UnfoldSecretsPlaceholdersInYaml(d.Secrets, initialYaml)
			if err != nil {
				logrus.Errorf("[E] Unfolding secrets for { Tier }: %v", err)
			} else {

				resultingYaml, err := yaml.Marshal(interfaceObject)

				decoder := yaml.NewDecoder(bytes.NewReader(resultingYaml))
				err = decoder.Decode(d.Tiers[d.tierName()])
				if err != nil {
					logrus.Errorf("[E] Unmarshalling tier configuration { %s } into { Tier } struct: %v", d.TierRedisChannel(d.tierName()), err)
				} else {
					d.Tiers[d.tierName()].initialYamlTemplate = initialYaml
				}

			}
			PauseConfigurationListening = false
		})
	}()
}

func (d *Deployment) InitializeTiersFields(tierName string) {
	if d.Tiers == nil {
		d.Tiers = map[string]*Tier{}
	}
	if d.Tiers[tierName] == nil {
		d.Tiers[tierName] = &Tier{}
	}
}
