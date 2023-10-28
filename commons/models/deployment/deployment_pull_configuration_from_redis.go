package deployment

import "time"

func (d *Deployment) PullConfigurationFromRedisWithDecryptionAndDeploy() {

	d.InitializeTiersFields(d.tierName())

	for {
		time.Sleep(1 * time.Second)

		err := GetYamlEntryFromRedisWithDecryption(d.Secrets, d.SecretsRedisConfigurationEntryKey(), &d.Secrets, d)
		if err != nil {
			continue
		}
		d.Secrets.DecryptEncryptedValues()

		err = GetYamlEntryFromRedisWithDecryption(d.Secrets, d.TierRedisConfigurationEntry(d.tierName()), d.Tiers[d.tierName()], d)
		if err != nil {
			continue
		}
		IsTierConfigurationInitialized = true

		err = GetYamlEntryFromRedisWithDecryption(d.Secrets, d.ControlSequenceRedisConfigurationEntry(d.tierName()), &d.Tiers[d.tierName()].LatestControlSequence, d)
		if err != nil {
			continue
		}
		d.Tiers[d.tierName()].LatestControlSequence.Execute(d)

		return
	}
}

func (d *Deployment) PullConfigurationFromRedisWithoutDecryptionForMultipleTiers(tierNames []string) {

	for {
		time.Sleep(1 * time.Second)

		err := GetYamlEntryFromRedisWithoutDecryption(d.SecretsRedisConfigurationEntryKey(), &d.Secrets)
		if err != nil {
			continue
		}

		for _, tierName := range tierNames {

			d.InitializeTiersFields(tierName)

			err = GetYamlEntryFromRedisWithoutDecryption(d.TierRedisConfigurationEntry(tierName), d.Tiers[tierName])
			if err != nil {
				continue
			}

			err = GetYamlEntryFromRedisWithoutDecryption(d.ControlSequenceRedisConfigurationEntry(tierName), &d.Tiers[tierName].LatestControlSequence)
			if err != nil {
				continue
			}

		}

		return
	}
}
