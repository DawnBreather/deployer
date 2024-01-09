package deployment

import (
  "bytes"
  "context"
  aws2 "deployer/commons/utils/aws"
  "github.com/goccy/go-yaml"
  cp "github.com/otiai10/copy"
  "github.com/sirupsen/logrus"
  "io/ioutil"
  "os"
  "strings"
  "sync"
)

var filesInfos = map[string]os.FileInfo{}
var filesInfosMutex = &sync.Mutex{}

func TransformValuePlaceholderIntoValue(data map[string]any, path string) string {
  return replacePlaceholdersWithValues(data, path)
}

func RemoveContentsOfDirectory(dir string) {
  cleanDirectory(dir)
}

func Unzip(source, destination string) error {
  return extractZipFile(context.TODO(), source, destination)
}

func hasFileChanged(filename string) bool {
  changed, _ := checkFileChanges(&filesInfos, filesInfosMutex, filename)
  return changed
}

func UnmarshalYAMLToStruct(filepath string, out any) error {
  return decodeYAMLFile(filepath, out)
}

//func TransformValuePlaceholderIntoValue(data map[string]any, path string) string {
//	re := regexp.MustCompile(`\${.*?}`)
//	result := path
//	for _, placeholder := range re.FindAllString(path, -1) {
//		valuePath := strings.TrimSuffix(strings.TrimPrefix(placeholder, `${`), `}`)
//		if value, isFound := getValueByPathFromMap(data, valuePath, ""); isFound {
//			result = strings.ReplaceAll(result, placeholder, value.(string))
//		}
//	}
//	return result
//
//	//if strings.HasPrefix(path, `${`) && strings.HasSuffix(path, "}") {
//	//	path = strings.TrimSuffix(strings.TrimPrefix(path, `${`), `}`)
//	//	if result, isFound := getValueByPathFromMap(data, path, ""); isFound {
//	//		return result.(string)
//	//	}
//	//}
//	//return path
//}
//
//// TODO: move to common utils
//func getValueByPathFromMap(data map[string]any, key string, passedKey string) (result any, found bool) {
//
//	keyAndPath := strings.SplitN(key, ".", 2)
//	currentKey := keyAndPath[0]
//	if passedKey != "" {
//		passedKey = passedKey + "." + currentKey
//	} else {
//		passedKey = currentKey
//	}
//
//	if _, isKeyExistInData := data[currentKey]; !isKeyExistInData {
//		logrus.Warnf("[W] key path { %s } not found", passedKey)
//		return
//	} else {
//
//		if len(keyAndPath) > 1 {
//			remainingPath := keyAndPath[1]
//			switch data[currentKey].(type) {
//			case map[string]any:
//				if result, found = getValueByPathFromMap(data[currentKey].(map[string]any), remainingPath, passedKey); found {
//					return
//				}
//			case Secrets:
//				if result, found = getValueByPathFromMap(data[currentKey].(Secrets), remainingPath, passedKey); found {
//					return
//				}
//			}
//		} else {
//			return data[currentKey], true
//		}
//	}
//
//	return nil, false
//}
//
//// TODO: move to common utilities
//// https://stackoverflow.com/a/33451503/4265419
//func RemoveContentsOfDirectory(dir string) {
//	d, err := os.Open(dir)
//	if err != nil {
//		logrus.Errorf("[E] Opening directory { %s }: %v", dir, err)
//	}
//	defer d.Close()
//	names, err := d.Readdirnames(-1)
//	if err != nil {
//		logrus.Errorf("[E] Opening directory { %s }: %v", dir, err)
//	}
//	for _, name := range names {
//		err = os.RemoveAll(filepath.Join(dir, name))
//		if err != nil {
//			logrus.Errorf("[E] Removing all at { %s }: %v", filepath.Join(dir, name), err)
//		}
//	}
//}
//
//// TODO: move to common utilities
//func Unzip(source, destination string) (err error) {
//	data, _ := os.ReadFile(source)
//	buffer := bytes.NewBuffer(data)
//	err = extract.Zip(context.TODO(), buffer, destination, nil)
//	if err != nil {
//		logrus.Errorf("[E] extracting zip file { %s } to { %s }: %v", source, destination, err)
//		return err
//	}
//
//	return nil
//}

// TODO: move to common utils
func removeContentsOfDestinationLocation(to string) {
  switch getLocationSchema(to) {
  case S3_SCHEMA:
    region, bucket, _, _ := parseS3BucketLocation(to)
    aws2.S3EmptyBucket(region, bucket)
  case FILESYSTEM_SCHEMA:
    RemoveContentsOfDirectory(to)
  default:
    return
  }
}

// TODO: move to common utils
func copy(from, to string) {
  switch getLocationSchema(to) {
  case S3_SCHEMA:
    region, _, _, regionFlag := parseS3BucketLocation(to)
    aws2.S3Upload(region, from, strings.TrimSuffix(to, regionFlag))
  case FILESYSTEM_SCHEMA:
    err := cp.Copy(from, to)
    if err != nil {
      logrus.Errorf("[E] copying from { %s } to { %s }", from, to)
    }
  default:
    return
  }
}

//var filesInfos = map[string]os.FileInfo{}
//var filesInfosMutex = &sync.Mutex{}
//
//func hasFileChanged(filename string) (res bool) {
//	info, err := os.Stat(filename)
//	if err != nil {
//		logrus.Errorf("[E] collecting fileInfo for { %s }: %v", filename, err)
//	} else {
//
//		filesInfosMutex.Lock()
//		if _, ok := filesInfos[filename]; ok {
//			res = info.ModTime() == filesInfos[filename].ModTime() && info.Size() == filesInfos[filename].Size()
//		}
//		filesInfos[filename] = info
//		filesInfosMutex.Unlock()
//	}
//	return
//}

func unmarshalYamlToStruct(filepath string, out any) error {
  // Read the YAML file into a byte slice
  yamlData, err := ioutil.ReadFile(filepath)
  if err != nil {
    return err
  }

  //return yaml.Unmarshal(yamlData, out)

  //if filepath == CONFIGURATION_STORAGE_SECRETS_PATH {
  //	logrus.Infof("%v", string(yamlData))
  //}

  //// Create a decoder and unmarshal the YAML data into a struct
  decoder := yaml.NewDecoder(bytes.NewReader(yamlData))
  return decoder.Decode(out)
}
