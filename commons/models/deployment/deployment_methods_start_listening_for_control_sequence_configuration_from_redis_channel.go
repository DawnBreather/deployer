package deployment

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"time"
)

var (
	IsLatestControlSequenceInitialized bool
)

func (d *Deployment) StartListeningForControlSequenceConfigurationFromRedisChannel() {
	go func() {
		subscribeForConfigurationFromRedisChannel(d.ControlSequenceRedisChannel(d.tierName()), func(message string) {
			for PauseConfigurationListening {
				time.Sleep(1 * time.Second)
			}
			for !IsSecretsDecryptionAccomplished {
				time.Sleep(1 * time.Second)
			}

			PauseConfigurationListening = true

			decoder := yaml.NewDecoder(bytes.NewReader([]byte(message)))
			err := decoder.Decode(&d.Tiers[d.tierName()].LatestControlSequence)
			if err != nil {
				logrus.Errorf("[E] Unmarshalling control sequence { %s } into { LatestControlSequence } struct: %v", d.ControlSequenceRedisChannel(d.tierName()), err)
			}
			IsLatestControlSequenceInitialized = true
			d.Tiers[d.tierName()].LatestControlSequence.Execute(d)

			PauseConfigurationListening = false
		})
	}()
}
