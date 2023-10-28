package deployment

import (
	git2 "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func (d *Deployment) StartPullingConfigurationFromGitRepository() {
	repository, err := cleanCloneConfigurationFromGit()

	if err == nil {

		if !IsConfigurationFromGitRepositoryRetrieved {
			IsConfigurationFromGitRepositoryRetrieved = true
		}

		startPullingLatestConfigurationFromGit(repository)

		//logrus.Errorf("[E] Extracting working tree from the Git repository: %v", err)
	} else {
		logrus.Errorf("[E] Getting working tree from Git repository: %v", err)
	}
}

func cleanCloneConfigurationFromGit() (*git2.Repository, error) {
	if _, err := os.Stat(CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH); os.IsNotExist(err) {
		// Directory doesn't exist, no need to remove
		return cloneConfigurationFromGit()
	}
	err := os.RemoveAll(CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH)
	if err != nil {
		logrus.Errorf("[E] Removing configuraiton folder { %s }: %v", CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH, err)
		return nil, err
	}

	return cloneConfigurationFromGit()
}

func cloneConfigurationFromGit() (*git2.Repository, error) {
	return git2.PlainClone(CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH, false, &git2.CloneOptions{
		URL:               GIT_CREDENTIALS.RepositoryHttpsUrl,
		Auth:              GIT_CREDENTIALS.GetBasicAuth(),
		RemoteName:        "origin",
		ReferenceName:     plumbing.NewBranchReferenceName(CONFIGURATION_STORAGE_GIT_BRANCH_NAME),
		SingleBranch:      true,
		NoCheckout:        false,
		Depth:             0,
		RecurseSubmodules: 0,
		Progress:          nil,
		Tags:              0,
		InsecureSkipTLS:   false,
		CABundle:          nil,
	})
}

func startPullingLatestConfigurationFromGit(repository *git2.Repository) {
	worktree, err := repository.Worktree()
	if err == nil {

		//var tiersYamlBytes, oldTiersYamlBytes []byte
		//var tiers map[string]*Tier
		var err error

		go func() {
			for {
				time.Sleep(time.Duration(CONFIGURATION_STORAGE_GIT_PULL_PERIOD) * time.Second)

				//repository, err = git2.PlainOpen(CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH + "/" + ".git")
				//if err != nil {
				//	logrus.Errorf("[E] opening Git repository { %s }: %v", CONFIGURATION_STORAGE_GIT_CLONE_DESTINATION_PATH, err)
				//	continue
				//}
				err = repository.Fetch(&git2.FetchOptions{
					Auth:       GIT_CREDENTIALS.GetBasicAuth(),
					Progress:   os.Stdout,
					RemoteName: "origin",
				})
				if err != nil {
					if err.Error() != "already up-to-date" {
						logrus.Errorf("[E] Fetching updates for the Git repository: %v", err)
						continue
					}
				}

				err = pullLatestConfigurationFromGit(worktree)
				if err != nil {
					if err.Error() != "already up-to-date" {
						logrus.Errorf("[E] Pulling updates for the Git repository: %v", err)
						continue
					}
				}

				logrus.Infof("[I] Successfully pulled")
			}
		}()
	} else {
		if err.Error() != "already up-to-date" {
			logrus.Errorf("[E] Opening worktree for the Git repository: {%v}", err)
		}
	}
}

func pullLatestConfigurationFromGit(w *git2.Worktree) error {
	return w.Pull(&git2.PullOptions{
		Auth:          GIT_CREDENTIALS.GetBasicAuth(),
		Progress:      os.Stdout,
		RemoteName:    "origin",
		ReferenceName: plumbing.NewBranchReferenceName(CONFIGURATION_STORAGE_GIT_BRANCH_NAME),
	})
}

//func parseLatestTiersConfiguration(tiersYamlBytes, oldTiersYamlBytes []byte) ([]byte, []byte, map[string]*Tier, error) {
//	var tiers map[string]*Tier
//	var err error
//
//	tiersYamlBytes, err = ioutil.ReadFile(CONFIGURATION_STORAGE_TIERS_PATH)
//	if err != nil {
//		logrus.Errorf("[E] Reading { %s } file: %v", CONFIGURATION_STORAGE_TIERS_PATH, err)
//		return tiersYamlBytes, oldTiersYamlBytes, nil, err
//	}
//
//	if !bytes.Equal(tiersYamlBytes, oldTiersYamlBytes) {
//		oldTiersYamlBytes = tiersYamlBytes
//
//		err = yaml.Unmarshal(tiersYamlBytes, &tiers)
//		if err != nil {
//			logrus.Errorf("[E] Unmarshalling tiers configuration: %v", err)
//			return tiersYamlBytes, oldTiersYamlBytes, nil, err
//		}
//	}
//
//	return tiersYamlBytes, oldTiersYamlBytes, tiers, nil
//}
