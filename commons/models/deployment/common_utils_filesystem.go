package deployment

import (
  "github.com/sirupsen/logrus"
  "os"
  "path/filepath"
  "sync"
)

// cleanDirectory removes all contents of the specified directory.
func cleanDirectory(dir string) error {
  d, err := os.Open(dir)
  if err != nil {
    logrus.Errorf("Error opening directory %s: %v", dir, err)
    return err
  }
  defer d.Close()

  names, err := d.Readdirnames(-1)
  if err != nil {
    logrus.Errorf("Error reading directory names in %s: %v", dir, err)
    return err
  }

  for _, name := range names {
    path := filepath.Join(dir, name)
    if err := os.RemoveAll(path); err != nil {
      logrus.Errorf("Error removing all content in %s: %v", path, err)
    }
  }
  return nil
}

// CheckFileChanges checks if the specified file has changed since the last check.
func checkFileChanges(filesInfos *map[string]os.FileInfo, filesInfosMutex *sync.Mutex, filename string) (bool, error) {
  info, err := os.Stat(filename)
  if err != nil {
    logrus.Errorf("Error collecting fileInfo for %s: %v", filename, err)
    return false, err
  }

  filesInfosMutex.Lock()
  defer filesInfosMutex.Unlock()

  prevInfo, ok := (*filesInfos)[filename]
  (*filesInfos)[filename] = info
  if !ok {
    return false, nil
  }

  return info.ModTime() == prevInfo.ModTime() && info.Size() == prevInfo.Size(), nil
}
