package gitwrapper_test

import (
  . "deployer/commons/utils/gitwrapper"
  "go.uber.org/zap"
  "go.uber.org/zap/zaptest/observer"
  "os"
  "path/filepath"
  "testing"
)

func TestNew(t *testing.T) {
  // Setup
  core, observedLogs := observer.New(zap.InfoLevel)
  logger := zap.New(core)
  repositoryURL := "https://example.com/repo.git"
  workingDirectory := "/tmp/workdir"
  failOnErrors := false

  // Exercise
  gitWrapper := New(repositoryURL, workingDirectory, failOnErrors, logger)

  // Verify
  if gitWrapper.RepositoryUrl != repositoryURL {
    t.Errorf("Expected repositoryUrl to be %s, but got %s", repositoryURL, gitWrapper.RepositoryUrl)
  }
  if gitWrapper.WorkingDirectory.GetPath() != workingDirectory {
    t.Errorf("Expected workingDirectory to be %s, but got %s", workingDirectory, gitWrapper.WorkingDirectory.GetPath())
  }
  if gitWrapper.FailOnErrorsEnabled != failOnErrors {
    t.Errorf("Expected failOnErrorsEnabled to be %v, but got %v", failOnErrors, gitWrapper.FailOnErrorsEnabled)
  }
  if gitWrapper.Logger != logger {
    t.Errorf("Expected logger to be %+v, but got %+v", logger, gitWrapper.Logger)
  }
  if len(observedLogs.All()) != 0 {
    t.Errorf("Expected no logs, but got %+v", observedLogs.All())
  }

  // Teardown (not needed for this test)
}

func TestCleanClone(t *testing.T) {
  // Setup
  core, observedLogs := observer.New(zap.InfoLevel)
  logger := zap.New(core)
  repositoryURL := "https://github.com/example/repo.git" // Use a real repository
  workingDirectory := "/tmp/test-repo"
  failOnErrors := false
  gitWrapper := New(repositoryURL, workingDirectory, failOnErrors, logger)

  // Exercise
  result := gitWrapper.CleanClone()

  // Verify
  if result == nil {
    t.Error("Expected CleanClone to return a GitWrapper instance, but got nil")
  }

  // Check if the working directory now contains a .git folder
  if _, err := os.Stat(filepath.Join(workingDirectory, ".git")); os.IsNotExist(err) {
    t.Error("Expected .git directory to exist after CleanClone, but it does not")
  }

  if len(observedLogs.All()) > 0 {
    t.Log("Logs were produced:")
    for _, log := range observedLogs.All() {
      t.Log(log)
    }
  }

  // Teardown
  os.RemoveAll(workingDirectory)
}

func TestCheckout(t *testing.T) {
  // Setup
  core, observedLogs := observer.New(zap.InfoLevel)
  logger := zap.New(core)
  repositoryURL := "https://example.com/repo.git"
  workingDirectory := "/tmp/workdir"
  failOnErrors := false
  gitWrapper := New(repositoryURL, workingDirectory, failOnErrors, logger)
  targetBranch := "main"

  // You may need to mock the git.Checkout function here

  // Exercise
  result := gitWrapper.Checkout(targetBranch)

  // Verify
  if result != gitWrapper {
    t.Error("Expected Checkout to return the original GitWrapper instance")
  }
  if len(observedLogs.All()) == 0 {
    t.Error("Expected logs to be produced, but none were found")
  }
  if gitWrapper.TargetRef != targetBranch {
    t.Errorf("Expected targetRef to be %s, but got %s", targetBranch, gitWrapper.TargetRef)
  }

  // Teardown
  os.RemoveAll(workingDirectory)
}

func TestChangeFileAndCommit(t *testing.T) {
  // Setup
  core, observedLogs := observer.New(zap.InfoLevel)
  logger := zap.New(core)
  repositoryURL := "https://example.com/repo.git"
  workingDirectory := "/tmp/workdir"
  failOnErrors := false
  gitWrapper := New(repositoryURL, workingDirectory, failOnErrors, logger)
  targetFile := "test.txt"
  fileContent := "Hello, world!"

  // You may need to mock the git.Add and git.Commit functions here

  // Exercise
  result, ok := gitWrapper.ChangeFileAndCommit(targetFile, fileContent)

  // Verify
  if !ok {
    t.Error("Expected ChangeFileAndCommit to succeed")
  }
  if result != gitWrapper {
    t.Error("Expected ChangeFileAndCommit to return the original GitWrapper instance")
  }
  if len(observedLogs.All()) == 0 {
    t.Error("Expected logs to be produced, but none were found")
  }

  // Check if the file was created and contains the correct content
  createdFileContent, err := os.ReadFile(filepath.Join(workingDirectory, targetFile))
  if err != nil {
    t.Errorf("Failed to read the created file: %v", err)
  }
  if string(createdFileContent) != fileContent {
    t.Errorf("Expected file content to be %s, but got %s", fileContent, string(createdFileContent))
  }

  // Teardown
  os.RemoveAll(workingDirectory)
}

func TestPullAndPush(t *testing.T) {
  // Setup
  core, observedLogs := observer.New(zap.InfoLevel)
  logger := zap.New(core)
  repositoryURL := "https://example.com/repo.git"
  workingDirectory := "/tmp/workdir"
  failOnErrors := false
  gitWrapper := New(repositoryURL, workingDirectory, failOnErrors, logger)
  gitWrapper.TargetRef = "main"

  // You may need to mock the git.Pull and git.Push functions here

  // Exercise
  ok := gitWrapper.PullAndPush()

  // Verify
  if !ok {
    t.Error("Expected PullAndPush to succeed")
  }
  if len(observedLogs.All()) == 0 {
    t.Error("Expected logs to be produced, but none were found")
  }

  // Teardown
  os.RemoveAll(workingDirectory)
}

func TestYamlFileExists(t *testing.T) {
  // Setup
  core, observedLogs := observer.New(zap.InfoLevel)
  logger := zap.New(core)
  repositoryURL := "https://example.com/repo.git"
  workingDirectory := "/tmp/workdir"
  failOnErrors := false
  gitWrapper := New(repositoryURL, workingDirectory, failOnErrors, logger)
  targetFile := "config.yaml"
  os.WriteFile(filepath.Join(workingDirectory, targetFile), []byte("content"), 0644)

  // Exercise
  fullPath, exists := gitWrapper.YamlFileExists("config")

  // Verify
  if !exists {
    t.Error("Expected YamlFileExists to return true for existing file")
  }
  if fullPath != filepath.Join(workingDirectory, targetFile) {
    t.Errorf("Expected fullPath to be %s, but got %s", filepath.Join(workingDirectory, targetFile), fullPath)
  }
  if len(observedLogs.All()) > 0 {
    t.Error("Expected no logs to be produced, but logs were found")
  }

  // Teardown
  os.RemoveAll(workingDirectory)
}

func TestRemoteBranchExists(t *testing.T) {
  // Setup
  core, observedLogs := observer.New(zap.InfoLevel)
  logger := zap.New(core)
  repositoryURL := "https://example.com/repo.git"
  workingDirectory := "/tmp/workdir"
  failOnErrors := false
  gitWrapper := New(repositoryURL, workingDirectory, failOnErrors, logger)
  targetBranch := "main"

  // You may need to mock the git.Branch function here

  // Exercise
  exists := gitWrapper.RemoteBranchExists(targetBranch)

  // Verify
  if exists {
    t.Error("Expected RemoteBranchExists to return false for non-existing branch")
  }
  if len(observedLogs.All()) == 0 {
    t.Error("Expected logs to be produced, but none were found")
  }

  // Teardown
  os.RemoveAll(workingDirectory)
}

func TestRepository(t *testing.T) {
  // Setup
  core, observedLogs := observer.New(zap.InfoLevel)
  logger := zap.New(core)
  repositoryURL := "https://user:pass@example.com/repo.git"
  expectedRepoURL := "https://example.com/repo.git"
  workingDirectory := "/tmp/workdir"
  failOnErrors := false
  gitWrapper := New(repositoryURL, workingDirectory, failOnErrors, logger)

  // Exercise
  repo := gitWrapper.Repository()

  // Verify
  if repo != expectedRepoURL {
    t.Errorf("Expected Repository to return %s, but got %s", expectedRepoURL, repo)
  }
  if len(observedLogs.All()) > 0 {
    t.Error("Expected no logs to be produced, but logs were found")
  }

  // Teardown
  // No files were created, so no need for teardown in this case
}
