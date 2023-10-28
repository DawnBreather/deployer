package deployment

import (
  "github.com/google/uuid"
  "github.com/sirupsen/logrus"
  "io"
  "net/http"
  "os"
  "strings"
)

func AgentName() string {
  hostname, err := os.Hostname()
  if err != nil {
    logrus.Errorf("[E] getting { hostname }: %v", err)
    hostname = uuid.New().String()
  }
  resp, _ := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
  if resp != nil {
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusOK {
      bodyBytes, err := io.ReadAll(resp.Body)
      if err != nil {
        logrus.Errorf("[E] reading EC2 instance name from response body: %v", err)
      }
      return string(bodyBytes)
    }
  }

  //remove unsupported symbols (for firebase nodes)
  hostname = strings.ReplaceAll(hostname, ".", "-")

  return hostname
}
