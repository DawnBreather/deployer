package deployment

import "context"

func PublishMessageToRedisChannel(channel, message string) error {
  return getConnectionToRedis().Publish(context.Background(), channel, message).Err()
}
