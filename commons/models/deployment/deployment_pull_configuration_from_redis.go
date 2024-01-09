package deployment

import (
  "time"
)

const sleepDuration = 1 * time.Second

func (d *Deployment) PullConfigurationFromRedisWithDecryptionAndDeploy() {
  d.InitializeTiersFields(d.tierName())

  for {
    time.Sleep(sleepDuration)

    if err := d.pullAndDecryptSecrets(); err != nil {
      continue
    }

    if err := d.pullAndDecryptTierConfiguration(); err != nil {
      continue
    }
    IsTierConfigurationInitialized = true

    if err := d.pullAndExecuteControlSequence(); err != nil {
      continue
    }

    return
  }
}

func (d *Deployment) pullAndDecryptSecrets() error {
  err := GetYamlEntryFromRedisWithDecryption(d.Secrets, d.SecretsRedisConfigurationEntryKey(), &d.Secrets, d)
  if err != nil {
    return err
  }
  d.Secrets.DecryptEncryptedValues()
  return nil
}

func (d *Deployment) pullAndDecryptTierConfiguration() error {
  return GetYamlEntryFromRedisWithDecryption(d.Secrets, d.TierRedisConfigurationEntry(d.tierName()), d.Tiers[d.tierName()], d)
}

func (d *Deployment) pullAndExecuteControlSequence() error {
  err := GetYamlEntryFromRedisWithDecryption(d.Secrets, d.ControlSequenceRedisConfigurationEntry(d.tierName()), &d.Tiers[d.tierName()].LatestControlSequence, d)
  if err != nil {
    return err
  }
  d.Tiers[d.tierName()].LatestControlSequence.Execute(d)
  return nil
}

func (d *Deployment) PullConfigurationFromRedisWithoutDecryptionForMultipleTiers(tierNames []string) {
  for {
    time.Sleep(sleepDuration)

    if err := d.pullSecretsWithoutDecryption(); err != nil {
      continue
    }

    for _, tierName := range tierNames {
      if err := d.pullTierConfigurationWithoutDecryption(tierName); err != nil {
        continue
      }
    }

    return
  }
}

func (d *Deployment) pullSecretsWithoutDecryption() error {
  return GetYamlEntryFromRedisWithoutDecryption(d.SecretsRedisConfigurationEntryKey(), &d.Secrets)
}

func (d *Deployment) pullTierConfigurationWithoutDecryption(tierName string) error {
  d.InitializeTiersFields(tierName)

  if err := GetYamlEntryFromRedisWithoutDecryption(d.TierRedisConfigurationEntry(tierName), d.Tiers[tierName]); err != nil {
    return err
  }

  return GetYamlEntryFromRedisWithoutDecryption(d.ControlSequenceRedisConfigurationEntry(tierName), &d.Tiers[tierName].LatestControlSequence)
}
