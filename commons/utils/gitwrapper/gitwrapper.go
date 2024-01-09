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
  Logger              *zap.Logger
  FailOnErrorsEnabled bool
  RepositoryUrl       string
  WorkingDirectory    *path.Path

  TargetRef string
}

//func (gs *GitWrapper) Initialize() *GitWrapper {
//	return gs
//}

func New(repositoryUrl, workingDirectory string, failOnErrorsEnabled bool, logger *zap.Logger) *GitWrapper {
  return &GitWrapper{
    RepositoryUrl:       repositoryUrl,
    WorkingDirectory:    path.New(workingDirectory),
    FailOnErrorsEnabled: failOnErrorsEnabled,
    Logger:              logger,
  }
}

func (gs *GitWrapper) CleanClone() *GitWrapper {
  gs.WorkingDirectory.RemoveIfExists()
  output, err := git.Clone(clone.Repository(gs.RepositoryUrl), clone.Directory(gs.WorkingDirectory.GetPath()))
  if err != nil {
    gs.Logger.Error("git: failed cloning repository", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("repository", url.CleanUrlFromCredentials(gs.RepositoryUrl)), zap.String("destination", gs.WorkingDirectory.GetFullPath()))
    if gs.FailOnErrorsEnabled {
      os.Exit(1)
    }
    return gs
  }

  gs.Logger.Info(output, zap.String("repository", url.CleanUrlFromCredentials(gs.RepositoryUrl)), zap.String("destination", gs.WorkingDirectory.GetFullPath()))

  return gs
}

func (gs *GitWrapper) Checkout(target string) *GitWrapper {

  if gs.TargetRef == target {
    return gs
  }

  gs.WorkingDirectory.ChdirIfDir()

  output, err := git.Checkout(checkout.NewBranch(target), checkout.Force)
  if err != nil {
    gs.Logger.Error("git: error checking out", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("target", target))
    return gs
  }

  gs.Logger.Info("git: checked out branch. "+output, zap.String("target", target))

  output, err = git.Status()
  if err != nil {
    gs.Logger.Error("git: error getting status", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)))
    return gs
  }
  gs.Logger.Info("git: status", zap.String("output", output))
  //
  //output, err = git.Fetch(fetch.RefSpec(target))
  //if err != nil {
  //	gs.logger.Error("git: error fetching", zap.NamedError("error", err), zap.String("target", fmt.Sprintf("origin %s", target)))
  //	return gs
  //}

  output, err = git.Reset(reset.Hard, reset.Path(target))
  if err != nil {
    gs.Logger.Error("git: error resetting "+fmt.Sprintf("{ %s }", output), zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("target", fmt.Sprintf("origin %s", target)))
    return gs
  } else {
    gs.Logger.Info("git: hard reset " + fmt.Sprintf("{ %s }", output))
  }

  output, err = git.Pull(pull.Repository("origin"), pull.Refspec(target), pull.Force, pull.NoRebase, pull.AllowUnrelatedHistories, pull.StrategyOption("theirs"))
  if err != nil {
    gs.Logger.Error("git: error pulling "+fmt.Sprintf("{ %s }", output), zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("target", fmt.Sprintf("origin %s", target)))
    return gs
  }

  gs.Logger.Info("git: pulled branch "+output, zap.String("target", fmt.Sprintf("origin %s", target)))

  gs.TargetRef = target

  return gs
}

func (gs *GitWrapper) ChangeFileAndCommit(target, content string) (this *GitWrapper, ok bool) {
  gs.WorkingDirectory.ChdirIfDir()

  targetFilePath := path.New(target)
  if !targetFilePath.Exists() {
    targetFilePath.MkdirAll(0755)
    file, err := os.Create(target)
    if err != nil {
      gs.Logger.Error("error creating file", zap.NamedError("error", err), zap.String("target", target))
      return gs, false
    }
    err = file.Close()
    if err != nil {
      gs.Logger.Error("error closing file", zap.NamedError("error", err), zap.String("target", target))
      return gs, false
    }
  }

  err := os.WriteFile(target, []byte(content), 0644)
  if err != nil {
    gs.Logger.Error("git: error writing to file", zap.NamedError("error", err), zap.String("target", target))
    return gs, false
  }

  gs.Logger.Info("git: adjusted file", zap.String("target", target))

  output, err := git.Add(add.All)
  if err != nil {
    gs.Logger.Error("git: error staging files", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)))
    return gs, false
  }

  gs.Logger.Info("git: added file"+output, zap.String("target", target))

  output, err = git.Commit(commit.Message(fmt.Sprintf("Adjusted %s", target)))
  if err != nil {
    gs.Logger.Error("git: error committing "+fmt.Sprintf("{ %s }", output), zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)))
    return gs, false
  }

  gs.Logger.Info("git: committed " + output)

  return gs, true
}

func (gs *GitWrapper) PullAndPush() (ok bool) {
  output, err := git.Pull(pull.Repository("origin"), pull.Refspec(gs.TargetRef))
  if err != nil {
    gs.Logger.Error("git: error pulling", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("target", fmt.Sprintf("origin %s", gs.TargetRef)))
    return false
  }

  gs.Logger.Info("git: pulled branch "+output, zap.String("target", fmt.Sprintf("origin %s", gs.TargetRef)))

  output, err = git.Push(push.Remote("origin"), push.RefSpec(gs.TargetRef))
  if err != nil {
    gs.Logger.Error("git: error pushing", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)), zap.String("target", fmt.Sprintf("origin %s", gs.TargetRef)))
    return false
  }

  gs.Logger.Info("git: pushed branch "+output, zap.String("target", fmt.Sprintf("origin %s", gs.TargetRef)))

  return true
}

func (gs *GitWrapper) YamlFileExists(target string) (fullPath string, exists bool) {
  var extensions = []string{"yaml", "yml"}

  gs.WorkingDirectory.ChdirIfDir()

  for _, extension := range extensions {
    targetPathWithExtension := path.New(filepath.FromSlash(fmt.Sprintf("%s.%s", target, extension)))
    if targetPathWithExtension.IsExistAndFile() {
      return targetPathWithExtension.GetPath(), true
    }
  }

  return
}

func (gs *GitWrapper) RemoteBranchExists(targetBranch string) bool {

  gs.WorkingDirectory.ChdirIfDir()

  output, err := git.Branch(branch.All)
  if err != nil {
    gs.Logger.Error("git: error listing branches", zap.NamedError("error", fmt.Errorf("%s, %s", err.Error(), output)))
  }

  return strings.Contains(output, fmt.Sprintf("remotes/origin/%s", targetBranch))
}

func (gs *GitWrapper) Repository() string {
  return url.CleanUrlFromCredentials(gs.RepositoryUrl)
}
