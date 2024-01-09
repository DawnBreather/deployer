package deployment

import (
  "bytes"
  "gopkg.in/yaml.v3"
  "io/ioutil"
)

func decodeYAMLFile(filepath string, out any) error {
  yamlData, err := ioutil.ReadFile(filepath)
  if err != nil {
    return err
  }

  decoder := yaml.NewDecoder(bytes.NewReader(yamlData))
  return decoder.Decode(out)
}
