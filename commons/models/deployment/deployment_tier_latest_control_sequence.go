package deployment

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"syscall"
)

// TODO: move filename into global constant
func (l *LatestControlSequence) Execute(d *Deployment) {
	logrus.Infof("CONTROL SEQUENCE -> EXECUTION STARTED")
	defer logrus.Infof("CONTROL SEQUENCE -> EXECUTION ENDED")
	PauseConfigurationListening = true

	logrus.Infof("[I] Identified { configuration } command { %#v } created at { %s }", l.CommandSequence, l.CreatedAt)

	var err error

	var inStatusLatestCotrolSequence LatestControlSequence

	GetJsonEntryFromRedis(fmt.Sprintf(string(RDB_LATEST_CONTROL_SEQUENCE_STATUS_PATH), d.Name(), d.tierName(), agentId), &inStatusLatestCotrolSequence)

	if inStatusLatestCotrolSequence.CreatedAt != l.CreatedAt {

		for _, command := range l.CommandSequence {

			logrus.Infof("[I] Executing { %s } command", command)

			switch command {
			case COMMAND_SEQUENCE_DEPLOY:
				logrus.Infof("[I] { %s } control sequence identified: executing...", command)
				d.Tiers[d.tierName()].DeployArtifacts(d)
			case COMMAND_SEQUENCE_RESTART_ENTRYPOINT:
				logrus.Infof("[I] { %s } control sequence identified: executing...", command)
				if entrypointCommand != nil {
					err = entrypointCommand.Process.Signal(syscall.SIGINT)
				}

				d.Tiers[d.tierName()].StartRunningEntrypoint()
			case COMMAND_SEQUENCE_STOP_ENTRYPOINT:
				logrus.Infof("[I] { %s } control sequence identified: executing...", command)
				if entrypointCommand != nil {
					err = entrypointCommand.Process.Signal(syscall.SIGTERM)

					if err != nil {
						logrus.Errorf("[E] sending SIGTERM signal to entrypoint process: %v", err)
					}
				}

			case COMMAND_SEQUENCE_START_ENTRYPOINT:
				logrus.Infof("[I] { %s } control sequence identified: executing...", command)
				d.Tiers[d.tierName()].StartRunningEntrypoint()
			default:
				logrus.Warnf("[W] Command sequence { %s } not recognized", command)
				return
			}
		}

	} else {
		logrus.Infof("[I] Skipping: sequence had been executed { %s == %s }", inStatusLatestCotrolSequence.CreatedAt, l.CreatedAt)
	}

	//PutIntoRedis(fmt.Sprintf(string(RDB_LATEST_CONTROL_SEQUENCE_STATUS_PATH), d.Name(), d.tierName(), agentId), l, 6000*time.Second)
	GetNodeStatus().LatestControlSequence = *l
	PauseConfigurationListening = false

}
