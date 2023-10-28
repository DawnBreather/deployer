package deployment

import (
	"fmt"
	"github.com/DawnBreather/go-commons/go/pkg/mod/github.com/zabawaba99/firego@v0.0.0-20190331000051-3bcc4b6a4599"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

func (d *Deployment) StartListeningForTierConfiguration() {

	//var data map[string]interface{}
	//err := json.Unmarshal([]byte(sampleSchema), &data)
	//if err != nil {
	//	fmt.Println("error unmarshalling json", err)
	//}

	//err = f.Set(data)
	//if err != nil {
	//  fmt.Println("error sending data", err)
	//}

	go func() {
		notifications := make(chan firego.Event)

		//tierRefPath := fmt.Sprintf("environments/%v/tier", d.Name())

		//tierRef, _ := f.Ref(tierRefPath)
		//tierRef, _ := GetTierRef(d.Name(), d.tierName())
		//tierRef, _ := f.Ref("")

		var tierRef = d.tierRef()

		if err := tierRef.Watch(notifications); err != nil {
			log.Fatal(err)
		}

		defer tierRef.StopWatching()

		for /*event*/ event := range notifications {

			for PauseConfigurationListening {
				time.Sleep(1 * time.Second)
			}

			//if ! IsTierConfigurationInitialized {
			//	j, err := json.Marshal(event.Data)
			//	if err != nil {
			//		logrus.Errorf("[E] marshalling json: %v", err)
			//	}
			//	err = json.Unmarshal(j, d)
			//	if err != nil {
			//		logrus.Errorf("[E] unmarshalling json: %v", err)
			//	}
			//	fmt.Printf("Deployment: %v", d)
			//	IsTierConfigurationInitialized = true
			//}
			fmt.Printf("Event %#v\n", event)
			fmt.Printf("Type: %s\n", event.Type)
			fmt.Printf("Path: %s\n", event.Path)
			fmt.Printf("Data: %v\n", event.Data)
			if d.Tiers == nil {
				d.Tiers = map[string]*Tier{}
			}
			if d.Tiers[d.tierName()] == nil {
				d.Tiers[d.tierName()] = &Tier{}
			}
			err := tierRef.Value(d.Tiers[d.tierName()])
			if err != nil {
				logrus.Errorf("[E] unmarshalling value from { %s } into object of { Tier } struct: %v", strings.TrimPrefix(tierRef.URL(), FIREBASE_RTDB_URL), err)
			} else {
				if !IsTierConfigurationInitialized {
					IsTierConfigurationInitialized = true
				}
			}
		}
		fmt.Printf("Notifications have stopped")
		go d.StartListeningForTierConfiguration()
	}()
}

func (d *Deployment) StartListeningForSecretsConfiguration() {

	go func() {
		notifications := make(chan firego.Event)

		var secretsRef = d.environmentSecretsRef()

		if err := secretsRef.Watch(notifications); err != nil {
			log.Fatal(err)
		}

		defer secretsRef.StopWatching()

		for /*event*/ _ = range notifications {

			for PauseConfigurationListening {
				time.Sleep(1 * time.Second)
			}

			err := secretsRef.Value(&d.Secrets)
			if err != nil {
				logrus.Errorf("[E] unmarshalling value from { %s } into object of { Secrets } struct: %v", strings.TrimPrefix(secretsRef.URL(), FIREBASE_RTDB_URL), err)
			}
		}
		fmt.Printf("Notifications have stopped")
		go d.StartListeningForSecretsConfiguration()
	}()
}

func (d *Deployment) StartListeningForMetadataConfiguration() {

	go func() {
		notifications := make(chan firego.Event)

		var metadataRef = d.environmentMetadataRef()

		if err := metadataRef.Watch(notifications); err != nil {
			log.Fatal(err)
		}

		defer metadataRef.StopWatching()

		for /*event*/ _ = range notifications {

			for PauseConfigurationListening {
				time.Sleep(1 * time.Second)
			}

			err := metadataRef.Value(&d.Metadata)
			if err != nil {
				logrus.Errorf("[E] unmarshalling value from { %s } into object of { Metadata } struct: %v", strings.TrimPrefix(metadataRef.URL(), FIREBASE_RTDB_URL), err)
			}
		}
		fmt.Printf("Notifications have stopped")
		go d.StartListeningForMetadataConfiguration()
	}()
}

func (d *Deployment) StartListeningForLatestControlSequence() {

	go func() {
		notifications := make(chan firego.Event)

		//latestControlSequenceCreationDateRefPath := fmt.Sprintf("environments/%v/tier/latest_control_sequence/created_at", d.Name())

		//latestControlSequenceTimestampDateRef, _ := f.Ref(latestControlSequenceCreationDateRefPath)
		//latestControlSequenceTimestampDateRef, _ := GetLatestControlSequenceTimestampRef(d.Name(), d.tierName())

		var latestControlSequenceTimestampDateRef = d.latestControlSequenceTimestampRef()

		if err := latestControlSequenceTimestampDateRef.Watch(notifications); err != nil {
			log.Fatal(err)
		}

		defer latestControlSequenceTimestampDateRef.StopWatching()

		for event := range notifications {

			for PauseConfigurationListening {
				time.Sleep(1 * time.Second)
			}

			fmt.Printf("Event %#v\n", event)
			fmt.Printf("Type: %s\n", event.Type)
			fmt.Printf("Path: %s\n", event.Path)
			fmt.Printf("Data: %v\n", event.Data)

			//latestControlSequenceRefPath := fmt.Sprintf("environments/%v/tier/latest_control_sequence", d.Name())
			//latestControlSequenceRef, _ := f.Ref(latestControlSequenceRefPath)
			//latestControlSequenceRef, _ := GetLatestControlSequenceRef(d.Name(), d.tierName())
			err := d.latestControlSequenceRef().Value(&d.Tiers[d.tierName()].LatestControlSequence)

			if err != nil {
				logrus.Errorf("[E] unmarshalling value from { %s } into { d }: %v", strings.TrimPrefix(d.latestControlSequenceRef().URL(), FIREBASE_RTDB_URL), err)
			} else {
				d.Tiers[d.tierName()].LatestControlSequence.Execute(d)
			}
		}
		fmt.Printf("Notifications have stopped")
		go d.StartListeningForLatestControlSequence()
	}()
}
