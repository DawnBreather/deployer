package deployment

import (
	"strings"
)

const (
	KMS_ENCRYPTED_STRING_PREFIX = "enc::awskms:"
)

var (
	IsSecretsDecryptionAccomplished = false
)

func (s *Secrets) DecryptEncryptedValues() {
	decryptEncryptedValuesInSecrets((*map[string]any)(s))
	IsSecretsDecryptionAccomplished = true
}

func decryptEncryptedValuesInSecrets(obj *map[string]any) {
	for k, v := range *obj {
		switch v.(type) {
		case string:
			// Replace the string value with the new string value
			decryptedString := decryptSecretIfEncrypted(v.(string))
			(*obj)[k] = decryptedString
		case map[string]interface{}:
			// Recursively call the function on the nested map
			nestedObj := v.(map[string]interface{})
			decryptEncryptedValuesInSecrets(&nestedObj)
		case Secrets:
			// Recursively call the function on the nested map
			nestedObj := v.(Secrets)
			decryptEncryptedValuesInSecrets((*map[string]any)(&nestedObj))
		}
	}
}

func decryptSecretIfEncrypted(secretText string) string {
	if strings.HasPrefix(secretText, KMS_ENCRYPTED_STRING_PREFIX) {
		textToDecrypt := strings.TrimPrefix(secretText, KMS_ENCRYPTED_STRING_PREFIX)
		return decryptKmsEncryptedSecretString(textToDecrypt)
	}
	return secretText
}
