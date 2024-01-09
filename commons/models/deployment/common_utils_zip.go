package deployment

import (
  "bytes"
  "context"
  "github.com/codeclysm/extract/v3"
  "github.com/sirupsen/logrus"
  "os"
)

func extractZipFile(ctx context.Context, source, destination string) error {
  data, err := os.ReadFile(source)
  if err != nil {
    logrus.Errorf("Error reading zip file %s: %v", source, err)
    return err
  }

  buffer := bytes.NewBuffer(data)
  err = extract.Zip(ctx, buffer, destination, nil)
  if err != nil {
    logrus.Errorf("Error extracting zip file %s to %s: %v", source, destination, err)
    return err
  }

  return nil
}
