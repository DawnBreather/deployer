package gitwrapper

import "go.uber.org/zap"

type Logger interface {
  Info(msg string, fields ...zap.Field)
  Error(msg string, err error, fields ...zap.Field)
}

type FileSystem interface {
  RemoveIfExists(path string) error
  ChdirIfDir(path string) error
  CreateFile(target string) error
  WriteFile(target string, content []byte) error
  FileExists(target string) bool
}

type GitClient interface {
  Clone(repoURL, directory string) (string, error)
  Checkout(branch string, options ...string) (string, error)
  Status() (string, error)
  Fetch(refSpec string) (string, error)
  Reset(options ...string) (string, error)
  Pull(repository, refSpec string, options ...string) (string, error)
  Push(remote, refSpec string) (string, error)
  Branch(options ...string) (string, error)
  Add(options ...string) (string, error)
  Commit(message string) (string, error)
}
