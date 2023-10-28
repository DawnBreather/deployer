package deployment

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"time"
)

func (d *Deployment) StartListeningForTierConfigurationFromRedisChannel() {
	initializeTiersFields(d)
	go func() {
		subscribeForConfigurationFromRedisChannel(d.TierRedisChannel(d.tierName()), func(message string) {
			for PauseConfigurationListening {
				time.Sleep(1 * time.Second)
			}

			decoder := yaml.NewDecoder(bytes.NewReader([]byte(message)))
			err := decoder.Decode(d.Tiers[d.tierName()])
			if err != nil {
				logrus.Errorf("[E] Unmarshalling tier configuration { %s } into { Tier } struct: %v", d.TierRedisChannel(d.tierName()), err)
			}
		})
	}()
}

func initializeTiersFields(d *Deployment) {
	if d.Tiers == nil {
		d.Tiers = map[string]*Tier{}
	}
	if d.Tiers[d.tierName()] == nil {
		d.Tiers[d.tierName()] = &Tier{}
	}
}
