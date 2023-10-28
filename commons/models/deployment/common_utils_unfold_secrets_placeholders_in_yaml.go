package deployment

import (
	"gopkg.in/yaml.v3"
	"strings"
)

func UnfoldSecretsPlaceholdersInYaml(s Secrets, yamlBytes []byte) (interface{}, error) {
	var data interface{}
	err := yaml.Unmarshal(yamlBytes, &data)
	if err != nil {
		return nil, err
	}

	return traverseAndModify(s, data, updateValue), err
}

func traverseAndModify(s Secrets, data interface{}, modifier func(Secrets, interface{}) interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		for key, val := range v {
			v[key] = traverseAndModify(s, val, modifier)
		}
	case []interface{}:
		for i, val := range v {
			v[i] = traverseAndModify(s, val, modifier)
		}
	default:
		return modifier(s, v)
	}

	return data
}

func updateValue(s Secrets, value interface{}) interface{} {
	if strVal, ok := value.(string); ok {
		if strings.HasPrefix(strVal, "${") && strings.HasSuffix(strVal, "}") {
			return TransformValuePlaceholderIntoValue(s, strVal)
		}
	}
	return value
}
