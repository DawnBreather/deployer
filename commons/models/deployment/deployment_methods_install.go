package deployment

import (
	"deployer/binaries"
	"encoding/base64"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func (d *Deployment) Install() {
	logrus.Infof("STARTING INSTALLATION SEQUENCE")
	defer logrus.Infof("INSTALLATION SEQUENCE ENDED")
	osFamily := runtime.GOOS
	if osFamily == "windows" && strings.ToLower(os.Getenv(ABZ_DEPLOYER_AGENT_INSTALLATION_ENABLED)) == "true" {
		logrus.Infof("[I] Identified { %s } OS: installing windows service { abz.deployer }", osFamily)
		file, err := os.Create("nssm.exe")
		if err != nil {
			logrus.Errorf("[E] creating file { nssm.exe }: %v", err)
		}

		nssmBinaryBytes, _ := base64.StdEncoding.DecodeString(binaries.NssmBinaryBase64)
		_, err = file.Write(nssmBinaryBytes)
		if err != nil {
			logrus.Errorf("[E] writing file { nssm.exe }: %v", err)
		}
		err = file.Close()
		if err != nil {
			logrus.Errorf("[E] closing file { nssm.exe }: %v", err)
		}

		//currentDirectory, err := os.Getwd()
		//if err != nil {
		//	logrus.Errorf("[E] getting working directory: %v", err)
		//}

		// https://stackoverflow.com/a/73409834/4265419
		//testSequence := []string{"/C", "dir"}

		executableFilePath, _ := filepath.Abs(os.Args[0])

		serviceStopCommandSequence := []string{"/C", "nssm.exe", "stop", "abz.deployer"}
		serviceRemoveCommandSequence := []string{"/C", "nssm.exe", "remove", "abz.deployer", "confirm"}
		serviceInstallCommandSequence := []string{"/C", "nssm.exe", "install", "abz.deployer", executableFilePath}
		serviceSetEnvironmentVariableCommandSequence := []string{"/C", "nssm.exe", "set", "abz.deployer", "AppEnvironmentExtra", fmt.Sprintf("%s=%s", ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID, os.Getenv(ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID)), fmt.Sprintf("%s=%s", ABZ_DEPLOYER_AGENT_TIER_ID, strings.ReplaceAll(os.Getenv(ABZ_DEPLOYER_AGENT_TIER_ID), "::", ".")), fmt.Sprintf("%s=%s", ABZ_DEPLOYER_AGENT_FIREBASE_AUTHENTICATION_JSON_URL, os.Getenv(ABZ_DEPLOYER_AGENT_FIREBASE_AUTHENTICATION_JSON_URL))}
		serviceStartCommandSequence := []string{"/C", "nssm.exe", "start", "abz.deployer"}

		//logrus.Infof("[I] Removing { abz.deployer } service if exists | OUTPUT")
		//cmd := exec.Command("cmd", serviceRemoveCommandSequence...)
		//
		//var stdBuffer bytes.Buffer
		//mw := io.MultiWriter(os.Stdout, &stdBuffer)
		//
		//cmd.Stdout = mw
		//cmd.Stderr = mw
		//
		//if err := cmd.Run(); err != nil {
		//	logrus.Errorf("OOPS")
		//}
		//
		//var dstB []byte
		//_, err = hex.Decode(dstB, stdBuffer.Bytes())
		//if err != nil {
		//	logrus.Errorf(err.Error())
		//}
		//logrus.Info("========", string(dstB), "=========")
		//
		//logrus.Infof("[I] Removing { abz.deployer } service if exists | GO-SHELL")
		//cm := execute.ExecTask{
		//	Command:     "./nssm.exe",
		//	Args:        []string{"remove", "abz.deployer", "confirm"},
		//	StreamStdio: false,
		//}
		//res, err := cm.Execute()
		//if err != nil {
		//	logrus.Info("========!!!!!", string(err.Error()), "!!!!!=========")
		//} else {
		//	logrus.Info("========!!!!!", strings.ReplaceAll(res.Stdout, "\x00", ""), "!!!!!=========")
		//	logrus.Info("========!!!!!", strings.ReplaceAll(res.Stderr, "\x00", ""), "!!!!!=========")
		//}
		//

		logrus.Infof("[I] Stopping { abz.deployer } service if exists")
		output, err := exec.Command("cmd", serviceStopCommandSequence...).CombinedOutput()
		output = []byte(strings.ReplaceAll(string(output), "\x00", ""))
		if err != nil {
			logrus.Errorf("[E] stopping { abz.deployer }: %v\n%s", err, string(output))
		} else {
			logrus.Infof("[I] %s", string(output))
		}

		logrus.Infof("[I] Removing { abz.deployer } service if exists")
		output, err = exec.Command("cmd", serviceRemoveCommandSequence...).CombinedOutput()
		output = []byte(strings.ReplaceAll(string(output), "\x00", ""))
		if err != nil {
			logrus.Errorf("[E] removing { abz.deployer }: %v\n%s", err, string(output))
		} else {
			logrus.Infof("[I] %s", string(output))
			time.Sleep(2 * time.Second)
		}

		logrus.Infof("[I] Installing { abz.deployer } service")
		c := exec.Command("cmd", serviceInstallCommandSequence...)
		output, err = c.CombinedOutput()
		output = []byte(strings.ReplaceAll(string(output), "\x00", ""))
		if err != nil {
			logrus.Errorf("[E] installing { abz.deployer } service over { nssm.exe }: %v\n%s", err, string(output))
		}
		output = []byte(strings.ReplaceAll(string(output), "\x00", ""))
		//logrus.Infof("[I] Installation  output: %s", string(output))
		logrus.Infof(string(output))
		//outputBase64Encoded := base64.StdEncoding.EncodeToString([]byte(output))
		//installationStatusRef, err := GetInstallationStatusRef(d.Name(), d.tierName(), agentId)
		//installationStatusRef, err := f.Ref(fmt.Sprintf("environments/%s/status/inventory/%s/installation", d.Name(), agentId))
		//if err != nil {
		//	logrus.Errorf("[E] referring to { installation status } in firebase: %v", err)
		//}

		// TODO: report installation status to Redis
		//err = d.installationStatusRef().Set(map[string]string{
		//	"state":  "accomplished",
		//	"output": string(output),
		//})
		//
		//
		//if err != nil {
		//	logrus.Errorf("[E] setting firebase ref { %s }: %v", strings.TrimPrefix(d.installationStatusRef().URL(), FIREBASE_RTDB_URL), err)
		//}

		logrus.Infof("[I] Setting { ABZ_DEPLOYER_ENVIRONMENT_ID } environment variable for { abz.deployer } service")
		output, err = exec.Command("cmd", serviceSetEnvironmentVariableCommandSequence...).CombinedOutput()
		output = []byte(strings.ReplaceAll(string(output), "\x00", ""))
		if err != nil {
			logrus.Errorf("[E] setting environment variables for { abz.deployer }: %v\n%s", err, string(output))
		}

		logrus.Infof("[I] Starting { abz.deployer } service")
		output, err = exec.Command("cmd", serviceStartCommandSequence...).CombinedOutput()
		output = []byte(strings.ReplaceAll(string(output), "\x00", ""))
		if err != nil {
			logrus.Errorf("[E] starting { abz.deployer } service: %v\n%s", err, output)
		}

		if strings.ToLower(os.Getenv(ABZ_DEPLOYER_AGENT_INSTALLATION_ENABLED)) == "true" {
			os.Exit(0)
		}

		//logrus.Infof("INSTALLATION SEQUENCE ENDED. EXITING")

	} else {
		logrus.Infof("[I] Skipping")
	}

}
