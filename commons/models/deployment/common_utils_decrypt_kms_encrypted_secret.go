package deployment

import (
	"encoding/base64"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/sirupsen/logrus"
	"time"
)

func decryptKmsEncryptedSecretString(encryptedString string) string {
	for {

		time.Sleep(1 * time.Second)

		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
			Config: aws.Config{
				Region: aws.String(KMS_KEY_REGION),
			},
		}))

		kmsSvc := kms.New(sess)

		ciphertextBlob, err := base64.StdEncoding.DecodeString(encryptedString)
		if err != nil {
			logrus.Errorf("[E] decoding base64 encoded string { %s }: %v", encryptedString, err)
			continue
			//return "", err
		}

		decryptInput := &kms.DecryptInput{
			CiphertextBlob: ciphertextBlob,
		}
		decryptOutput, err := kmsSvc.Decrypt(decryptInput)
		if err != nil {
			logrus.Errorf("[E] decrypting secret string { %s } by KMS key: %v", encryptedString, err)
			continue
		}

		return string(decryptOutput.Plaintext)
	}
}
