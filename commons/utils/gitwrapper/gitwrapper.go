package gitwrapper

import (
  "deployer/commons/utils/path"
  "deployer/commons/utils/url"
  "fmt"
  "github.com/ldez/go-git-cmd-wrapper/v2/add"
  "github.com/ldez/go-git-cmd-wrapper/v2/branch"
  "github.com/ldez/go-git-cmd-wrapper/v2/checkout"
  "github.com/ldez/go-git-cmd-wrapper/v2/clone"
  "github.com/ldez/go-git-cmd-wrapper/v2/commit"
  "github.com/ldez/go-git-cmd-wrapper/v2/git"
  "github.com/ldez/go-git-cmd-wrapper/v2/pull"
  "github.com/ldez/go-git-cmd-wrapper/v2/push"
  "github.com/ldez/go-git-cmd-wrapper/v2/reset"
  "go.uber.org/zap"
  "os"
  "path/filepath"
  "strings"
)

type GitWrapper struct {
  logger              *zap.Logger
  failOnErrorsEnabled bool
  repositoryUrl       string
  workingDirectory    *path.Path

  targetRef string
}

//func (gs *GitWrapper) Initialize() *GitWrapper {
//	return gs
//}

func New(repositoryUrl, workingDirectory string, failOnErrorsEnabled bool, logger *zap.Logger) *GitWrapper {
  return &GitWrapper{
    repositoryUrl:       repositoryUrl,
    workingDirectory:    path.New(workingDirectory),
    failOnErrorsEnabled: failOnErrorsEnabled,
    logger:              logger,
  }
}

func (gs *GitWrapper) CleanClone() *GitWrapper {
  gs.workingDirectory.RemoveIfExists()
  output, err := git.Clone(clone.Repository(gs.repositoryUrl), clone.Directory(gs.workingDirectory.GetPath()))
  if err != nil {
    gs.logger.Error("git: failed cloning repository", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("repository", url.CleanUrlFromCredentials(gs.repositoryUrl)), zap.String("destination", gs.workingDirectory.GetFullPath()))
    if gs.failOnErrorsEnabled {
      os.Exit(1)
    }
    return gs
  }

  gs.logger.Info(output, zap.String("repository", url.CleanUrlFromCredentials(gs.repositoryUrl)), zap.String("destination", gs.workingDirectory.GetFullPath()))

  return gs
}

func (gs *GitWrapper) Checkout(target string) *GitWrapper {

  if gs.targetRef == target {
    return gs
  }

  gs.workingDirectory.ChdirIfDir()

  output, err := git.Checkout(checkout.NewBranch(target), checkout.Force)
  if err != nil {
    gs.logger.Error("git: error checking out", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("target", target))
    return gs
  }

  gs.logger.Info("git: checked out branch. "+output, zap.String("target", target))

  output, err = git.Status()
  if err != nil {
    gs.logger.Error("git: error getting status", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)))
    return gs
  }
  gs.logger.Info("git: status", zap.String("output", output))
  //
  //output, err = git.Fetch(fetch.RefSpec(target))
  //if err != nil {
  //	gs.logger.Error("git: error fetching", zap.NamedError("error", err), zap.String("target", fmt.Sprintf("origin %s", target)))
  //	return gs
  //}

  output, err = git.Reset(reset.Hard, reset.Path(target))
  if err != nil {
    gs.logger.Error("git: error resetting "+fmt.Sprintf("{ %s }", output), zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("target", fmt.Sprintf("origin %s", target)))
    return gs
  } else {
    gs.logger.Info("git: hard reset " + fmt.Sprintf("{ %s }", output))
  }

  output, err = git.Pull(pull.Repository("origin"), pull.Refspec(target), pull.Force, pull.NoRebase, pull.AllowUnrelatedHistories, pull.StrategyOption("theirs"))
  if err != nil {
    gs.logger.Error("git: error pulling "+fmt.Sprintf("{ %s }", output), zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("target", fmt.Sprintf("origin %s", target)))
    return gs
  }

  gs.logger.Info("git: pulled branch "+output, zap.String("target", fmt.Sprintf("origin %s", target)))

  gs.targetRef = target

  return gs
}

func (gs *GitWrapper) ChangeFileAndCommit(target, content string) (this *GitWrapper, ok bool) {
  gs.workingDirectory.ChdirIfDir()

  targetFilePath := path.New(target)
  if !targetFilePath.Exists() {
    targetFilePath.MkdirAll(0755)
    file, err := os.Create(target)
    if err != nil {
      gs.logger.Error("error creating file", zap.NamedError("error", err), zap.String("target", target))
      return gs, false
    }
    err = file.Close()
    if err != nil {
      gs.logger.Error("error closing file", zap.NamedError("error", err), zap.String("target", target))
      return gs, false
    }
  }

  err := os.WriteFile(target, []byte(content), 0644)
  if err != nil {
    gs.logger.Error("git: error writing to file", zap.NamedError("error", err), zap.String("target", target))
    return gs, false
  }

  gs.logger.Info("git: adjusted file", zap.String("target", target))

  output, err := git.Add(add.All)
  if err != nil {
    gs.logger.Error("git: error staging files", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)))
    return gs, false
  }

  gs.logger.Info("git: added file"+output, zap.String("target", target))

  output, err = git.Commit(commit.Message(fmt.Sprintf("Adjusted %s", target)))
  if err != nil {
    gs.logger.Error("git: error committing "+fmt.Sprintf("{ %s }", output), zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)))
    return gs, false
  }

  gs.logger.Info("git: committed " + output)

  return gs, true
}

func (gs *GitWrapper) PullAndPush() (ok bool) {
  output, err := git.Pull(pull.Repository("origin"), pull.Refspec(gs.targetRef))
  if err != nil {
    gs.logger.Error("git: error pulling", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("target", fmt.Sprintf("origin %s", gs.targetRef)))
    return false
  }

  gs.logger.Info("git: pulled branch "+output, zap.String("target", fmt.Sprintf("origin %s", gs.targetRef)))

  output, err = git.Push(push.Remote("origin"), push.RefSpec(gs.targetRef))
  if err != nil {
    gs.logger.Error("git: error pushing", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("target", fmt.Sprintf("origin %s", gs.targetRef)))
    return false
  }

  gs.logger.Info("git: pushed branch "+output, zap.String("target", fmt.Sprintf("origin %s", gs.targetRef)))

  return true
}

func (gs *GitWrapper) YamlFileExists(target string) (fullPath string, exists bool) {
  var extensions = []string{"yaml", "yml"}

  gs.workingDirectory.ChdirIfDir()

  for _, extension := range extensions {
    targetPathWithExtension := path.New(filepath.FromSlash(fmt.Sprintf("%s.%s", target, extension)))
    if targetPathWithExtension.IsExistAndFile() {
      return targetPathWithExtension.GetPath(), true
    }
  }

  return
}

func (gs *GitWrapper) RemoteBranchExists(targetBranch string) bool {

  gs.workingDirectory.ChdirIfDir()

  output, err := git.Branch(branch.All)
  if err != nil {
    gs.logger.Error("git: error listing branches", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)))
  }

  return strings.Contains(output, fmt.Sprintf("remotes/origin/%s", targetBranch))
}

func (gs *GitWrapper) Repository() string {
  return url.CleanUrlFromCredentials(gs.repositoryUrl)
}
