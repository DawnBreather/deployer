package deployment

import (
  "github.com/google/uuid"
  "github.com/sirupsen/logrus"
  "io"
  "net/http"
  "os"
  "strings"
)

const ec2InstanceMetadataURL = "http://169.254.169.254/latest/meta-data/instance-id"

func AgentName() string {
  hostname, err := getHostName()
  if err != nil {
    hostname = generateUUID()
  }

  instanceID, err := getEC2InstanceID()
  if err == nil {
    return instanceID
  }

  return sanitizeHostName(hostname)
}

func getHostName() (string, error) {
  hostname, err := os.Hostname()
  if err != nil {
    logrus.Errorf("[E] getting { hostname }: %v", err)
    return "", err
  }
  return hostname, nil
}

func generateUUID() string {
  return uuid.New().String()
}

func getEC2InstanceID() (string, error) {
  resp, err := http.Get(ec2InstanceMetadataURL)
  if err != nil {
    logrus.Errorf("[E] making request to EC2 instance metadata: %v", err)
    return "", err
  }
  defer resp.Body.Close()

  if resp.StatusCode == http.StatusOK {
    bodyBytes, err := io.ReadAll(resp.Body)
    if err != nil {
      logrus.Errorf("[E] reading EC2 instance name from response body: %v", err)
      return "", err
    }
    return string(bodyBytes), nil
  }
  return "", nil
}

func sanitizeHostName(hostname string) string {
  return strings.ReplaceAll(hostname, ".", "-")
}
