package deployment

import (
	"github.com/sirupsen/logrus"
	"time"
)

var (
	IsLatestControlSequenceInitialized bool
)

func (d *Deployment) StartListeningForControlSequenceConfigurationFromGit() {
	go func() {
		for {

			time.Sleep(time.Duration(CONFIGURATION_STORAGE_GIT_PULL_PERIOD) * time.Second)

			for PauseConfigurationListening {
				time.Sleep(1 * time.Second)
			}

			if !IsLatestControlSequenceInitialized && IsConfigurationFromGitRepositoryRetrieved {
				err := unmarshalYamlToStruct(CONFIGURATION_STORAGE_CONTROL_SEQUENCES_PATH, &d.Tiers[d.tierName()].LatestControlSequence)
				if err != nil {
					logrus.Errorf("[E] Unmarshalling control sequence configuration { %s } into { LatestControlSequence } struct: %v", CONFIGURATION_STORAGE_CONTROL_SEQUENCES_PATH, err)
					continue
				}
				IsLatestControlSequenceInitialized = true
			}

			if hasFileChanged(CONFIGURATION_STORAGE_CONTROL_SEQUENCES_PATH) {
				err := unmarshalYamlToStruct(CONFIGURATION_STORAGE_CONTROL_SEQUENCES_PATH, &d.Tiers[d.tierName()].LatestControlSequence)
				if err != nil {
					logrus.Errorf("[E] Unmarshalling control sequence configuration { %s } into { LatestControlSequence } struct: %v", CONFIGURATION_STORAGE_CONTROL_SEQUENCES_PATH, err)
					continue
				}
				d.Tiers[d.tierName()].LatestControlSequence.Execute(d)
			}
		}
	}()
}
