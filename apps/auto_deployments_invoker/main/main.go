package main

import (
	"deployer/apps/auto_deployments_invoker"
	"deployer/commons/utils/path"
	"os"
)

var deployments = auto_deployments_invoker.Deployments{}

func main() {

	defer auto_deployments_invoker.Logger.Sync()

	os.Exit(auto_deployments_invoker.NewEnvironmentsConfigurations().
		Initialize().
		SubmitAutodeploy())
}

type EnvironmentsConfigurationsGitAdapter struct {
	destinationPath path.Path
}

func (e EnvironmentsConfigurationsGitAdapter) cloneRepository() {
	e.destinationPath.RemoveIfExists()
}
