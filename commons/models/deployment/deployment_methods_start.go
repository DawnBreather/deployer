package deployment

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func (d *Deployment) WaitForInitialization() {
	for !IsTierConfigurationInitialized {
		logrus.Infof("Waiting for configuration initailization")
		time.Sleep(3000 * time.Millisecond)
	}

	for !IsSecretsDecryptionAccomplished {
		logrus.Infof("Waiting for configuration initailization")
		time.Sleep(1000 * time.Millisecond)
	}
}

func (t *Tier) StartRunningEntrypoint() {
	//serviceStartCommandSequence := []string{"/C", "nssm.exe", "start", "abz.deployer"}

	//logrus.Infof("[I] Removing { abz.deployer } service if exists | OUTPUT")

	go func() {

		//entrypointCtx, EntrypointCtxCancel = context.WithCancel(ctx)
		//entrypointCtxDone = entrypointCtx.Done
		var err error

		if len(t.Entrypoint.Command) > 1 {
			entrypointCommand = exec.Command(t.Entrypoint.Command[0], t.Entrypoint.Command[1:]...)
			//cmd = exec.CommandContext(entrypointCtx, t.Entrypoint.Command[0], t.Entrypoint.Command[1:]...)
		} else if len(t.Entrypoint.Command) == 1 {
			entrypointCommand = exec.Command(t.Entrypoint.Command[0], "")
			//cmd = exec.CommandContext(entrypointCtx, t.Entrypoint.Command[0], "")
		} else {
			return
		}

		var stdBuffer bytes.Buffer
		mw := io.MultiWriter(os.Stdout, &stdBuffer)

		entrypointCommand.Stdout = mw
		entrypointCommand.Stderr = mw

		entrypointStdin, err = entrypointCommand.StdinPipe()
		if err != nil {
			logrus.Errorf("[E] taking over stdin of entrypoint process: %v", err)
		}

		if err := entrypointCommand.Run(); err != nil {
			if !strings.Contains(err.Error(), "signal: killed") {
				logrus.Errorf("[E] Executing { %v %v }: %v", entrypointCommand.Path, entrypointCommand.Args, err)
			}
		}

	}()

}
