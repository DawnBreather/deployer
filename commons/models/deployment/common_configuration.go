package deployment

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (d *Deployment) ControlSequenceRedisChannel(tierName string) string {
	return fmt.Sprintf("%s/control_sequences/%s", d.Name(), strings.ReplaceAll(tierName, "::", "."))
}

func (d *Deployment) TierRedisChannel(tierName string) string {
	return fmt.Sprintf("%s/tiers/%s", d.Name(), strings.ReplaceAll(tierName, "::", "."))
}

func (d *Deployment) SecretsRedisChannel() string {
	return fmt.Sprintf("%s/secrets", d.Name())
}

func (d *Deployment) SecretsRedisConfigurationEntryKey() string {
	return fmt.Sprintf("configuration/%s/secrets", d.Name())
}

func (d *Deployment) TierRedisConfigurationEntry(tierName string) string {
	return fmt.Sprintf("configuration/%s/tiers/%s", d.Name(), strings.ReplaceAll(tierName, "::", "."))
}

func (d *Deployment) ControlSequenceRedisConfigurationEntry(tierName string) string {
	return fmt.Sprintf("configuration/%s/control_sequences/%s", d.Name(), strings.ReplaceAll(tierName, "::", "."))
}

func subscribeForConfigurationFromRedisChannel(channel string, actionItem func(string)) {

	for {
		time.Sleep(1 * time.Second)

		// Create a new PubSub channel
		pubsub := getConnectionToRedis().Subscribe(context.Background(), channel)

		// Wait for confirmation that subscription is created
		_, err := pubsub.Receive(context.Background())
		if err != nil {
			logrus.Errorf("[E] error subscribing for Redis channel { %s }: %v", channel, err)
			continue
		}

		// Create a channel to receive messages
		ch := pubsub.Channel()

		// Read messages from the channel
		for msg := range ch {
			actionItem(msg.Payload)
			//fmt.Println("Received message:\n", msg.Payload)
		}
	}
}

func getConfigurationFromRedis(entryKey string, actionItem func(string)) {

}
