package deployment

import (
	"github.com/DawnBreather/go-commons/aws"
	"github.com/sirupsen/logrus"
	"github.com/zabawaba99/firego"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
	"strings"
	"time"
)

func StartKeepingFirebaseAuthenticated() {
	go func() {
		for {
			f = NewFirebaseClientAuthenticated()

			if !IsFirebaseAuthenticated {
				IsFirebaseAuthenticated = true
			}

			time.Sleep(59 * time.Minute)

		}
	}()
}

func NewFirebaseClientAuthenticated() *firego.Firebase {
	d, ok := ExtractFirebaseAuthenticationData()
	if !ok {
		logrus.Fatalf("[F] extracting Firebase authentication data from { %s }", ABZ_DEPLOYER_AGENT_FIREBASE_AUTHENTICATION_JSON_URL)
	}

	conf, err := google.JWTConfigFromJSON([]byte(d), "https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/firebase.database")
	if err != nil {
		logrus.Fatalf("[F] error gettign JWT config from auth.json")
	}

	return firego.New(FIREBASE_RTDB_URL, conf.Client(oauth2.NoContext))
}

func ExtractFirebaseAuthenticationData() (value string, ok bool) {
	// i.e. aws::ssm://some/path/in/ssm
	url := os.Getenv(ABZ_DEPLOYER_AGENT_FIREBASE_AUTHENTICATION_JSON_URL)
	prefixAndKey := strings.SplitN(url, "://", 2)

	var prefix, key string

	if len(prefixAndKey) > 1 {
		prefix = prefixAndKey[0]
		key = prefixAndKey[1]

		switch prefix {
		case "aws::ssm":
			return aws.SecretsManager.GetSecret(key)
		}

	} else {
		logrus.Fatalf("[F] parsing { ABZ_DEPLOYER_AGENT_FIREBASE_AUTHENTICATION_JSON_URL } value. Please provide proper URL.")
	}

	return
}
