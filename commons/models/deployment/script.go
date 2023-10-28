package deployment

import (
	"fmt"
	"github.com/DawnBreather/gosh"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type Script struct {
	Shell                string            `json:"shell" yaml:"shell"`
	EnvironmentVariables map[string]string `json:"environment_variables" yaml:"environment_variables"`
	Command              []string          `json:"command" yaml:"command"`
}

const (
	SHELL_TYPE_POWERSHELL = "powershell"
	SHELL_TYPE_BASH       = "bash"
	SHELL_TYPE_SHEBANG    = "sh"
)

func (s *Script) Execute(secrets Secrets) error {
	if len(s.Command) == 0 {
		return nil
	}
	logrus.Infof("SCRIPT -> EXECUTE STARTED")
	defer logrus.Infof("SCRIPT -> EXECUTE ENDED")

	compileSetOfEnvironmentVariables := func() (res []string) {
		for key, value := range s.EnvironmentVariables {
			res = append(res, fmt.Sprintf("%s=%s", key, TransformValuePlaceholderIntoValue(secrets, value)))
			//err := os.Setenv(key, transformValuePlaceholderIntoValue(secrets, value))
			//if err != nil {
			//  logrus.Warnf("[W] setting environment variable { %s }: %v", key, err)
			//}
		}
		return
	}

	executePowershellCommandSequence := func() error {
		//setEnvironmentVariable()
		for _, command := range s.Command {
			logrus.Infof("[I] executing { %s } command", command)
			err, _, stderr := gosh.PowershellOutput(command, compileSetOfEnvironmentVariables())
			if err != nil {
				errorMessage := fmt.Errorf("Error executing { powershell } script command { %s }: { err }=> %v ::: { stderr }=> %v", command, err, stderr)
				logrus.Warnf("[W] %s", errorMessage.Error())
				return errorMessage
			}
		}
		return nil
	}

	executeShellCommandSequence := func() error {
		//setEnvironmentVariable()
		for _, command := range s.Command {
			logrus.Infof("[I] executing { %s } command", command)
			err, _, stderr := gosh.ShellOutput(command, compileSetOfEnvironmentVariables())
			if err != nil {
				errorMessage := fmt.Errorf("Error executing { shell } script command { %s }: { err }=> %v ::: { stderr }=> %v", command, err, stderr)
				logrus.Warnf("[W] %s", errorMessage.Error())
				return errorMessage
			}
		}
		return nil
	}

	if s.Shell != "" {
		switch strings.ToLower(s.Shell) {
		case SHELL_TYPE_POWERSHELL:
			return executePowershellCommandSequence()
		case SHELL_TYPE_BASH, SHELL_TYPE_SHEBANG:
			return executeShellCommandSequence()
		}
	}

	if runtime.GOOS == "windows" {
		return executePowershellCommandSequence()
	} else {
		return executeShellCommandSequence()
	}

}
