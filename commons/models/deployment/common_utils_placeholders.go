package deployment

import (
  "github.com/sirupsen/logrus"
  "regexp"
  "strings"
)

func replacePlaceholdersWithValues(data map[string]any, path string) string {
  regex := regexp.MustCompile(`\${.*?}`)
  result := path
  for _, placeholder := range regex.FindAllString(path, -1) {
    valuePath := strings.TrimSuffix(strings.TrimPrefix(placeholder, `${`), `}`)
    if value, isFound := getValueByPathFromMap(data, valuePath, ""); isFound {
      result = strings.ReplaceAll(result, placeholder, value.(string))
    }
  }
  return result
}

func getValueByPathFromMap(data map[string]any, key string, passedKey string) (result any, found bool) {

  keyAndPath := strings.SplitN(key, ".", 2)
  currentKey := keyAndPath[0]
  if passedKey != "" {
    passedKey = passedKey + "." + currentKey
  } else {
    passedKey = currentKey
  }

  if _, isKeyExistInData := data[currentKey]; !isKeyExistInData {
    logrus.Warnf("[W] key path { %s } not found", passedKey)
    return
  } else {

    if len(keyAndPath) > 1 {
      remainingPath := keyAndPath[1]
      switch data[currentKey].(type) {
      case map[string]any:
        if result, found = getValueByPathFromMap(data[currentKey].(map[string]any), remainingPath, passedKey); found {
          return
        }
      case Secrets:
        if result, found = getValueByPathFromMap(data[currentKey].(Secrets), remainingPath, passedKey); found {
          return
        }
      }
    } else {
      return data[currentKey], true
    }
  }

  return nil, false
}
