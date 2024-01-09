package path_test

import (
  "deployer/commons/utils/path"
  "github.com/stretchr/testify/assert"
  "os"
  "testing"
)

func TestPath_SetPath(t *testing.T) {
  p := path.New("testpath")
  p.SetPath("newpath")
  assert.Equal(t, "newpath", p.GetPath())
}

func TestPath_GetFullPath(t *testing.T) {
  p := path.New("testpath")
  fullPath := p.GetFullPath()
  assert.NotEmpty(t, fullPath)
}

func TestPath_SetCompositePath(t *testing.T) {
  p := path.New("testpath")
  p.SetCompositePath("part1", "part2")
  assert.Equal(t, "part1/part2", p.GetPath())
}

func TestPath_Exists(t *testing.T) {
  p := path.New("testpath")
  assert.False(t, p.Exists())

  err := os.Mkdir("testpath", 0755)
  assert.NoError(t, err)
  defer os.Remove("testpath")

  assert.True(t, p.Exists())
}

func TestPath_MkdirAll(t *testing.T) {
  p := path.New("testdir/subdir")
  p.MkdirAll(0755)
  assert.True(t, p.Exists())
  defer os.RemoveAll("testdir")
}

func TestPath_IsFileOrDir(t *testing.T) {
  dirPath := path.New("testdir")
  dirPath.MkdirAll(0755)
  assert.Equal(t, path.DIRECTORY, dirPath.IsFileOrDir())
  defer os.RemoveAll("testdir")

  file, err := os.Create("testfile.txt")
  assert.NoError(t, err)
  file.Close()
  defer os.Remove("testfile.txt")

  filePath := path.New("testfile.txt")
  assert.Equal(t, path.FILE, filePath.IsFileOrDir())
}

func TestPath_IsFile(t *testing.T) {
  file, err := os.Create("testfile.txt")
  assert.NoError(t, err)
  file.Close()
  defer os.Remove("testfile.txt")

  p := path.New("testfile.txt")
  assert.True(t, p.IsFile())

  p.SetPath("not_a_file")
  assert.False(t, p.IsFile())
}

func TestPath_IsExistAndFile(t *testing.T) {
  file, err := os.Create("testfile.txt")
  assert.NoError(t, err)
  file.Close()
  defer os.Remove("testfile.txt")

  p := path.New("testfile.txt")
  assert.True(t, p.IsExistAndFile())

  p.SetPath("not_a_file")
  assert.False(t, p.IsExistAndFile())
}

func TestPath_IsExistAndDirectory(t *testing.T) {
  dirPath := path.New("testdir")
  dirPath.MkdirAll(0755)
  assert.True(t, dirPath.IsExistAndDirectory())
  defer os.RemoveAll("testdir")

  dirPath.SetPath("not_a_dir")
  assert.False(t, dirPath.IsExistAndDirectory())
}

func TestPath_IsDirectory(t *testing.T) {
  dirPath := path.New("testdir")
  dirPath.MkdirAll(0755)
  assert.True(t, dirPath.IsDirectory())
  defer os.RemoveAll("testdir")

  dirPath.SetPath("not_a_dir")
  assert.False(t, dirPath.IsDirectory())
}

func TestPath_Remove(t *testing.T) {
  file, err := os.Create("testfile.txt")
  assert.NoError(t, err)
  file.Close()

  p := path.New("testfile.txt")
  p.Remove()
  assert.False(t, p.Exists())
}

func TestPath_RemoveIfExists(t *testing.T) {
  file, err := os.Create("testfile.txt")
  assert.NoError(t, err)
  file.Close()

  p := path.New("testfile.txt")
  p.RemoveIfExists()
  assert.False(t, p.Exists())

  // Test with non-existing file
  p.SetPath("non_existent_file")
  p.RemoveIfExists()
}

func TestPath_ChdirIfDir(t *testing.T) {
  currentDir, err := os.Getwd()
  assert.NoError(t, err)

  dirPath := path.New("testdir")
  dirPath.MkdirAll(0755)
  defer os.RemoveAll("testdir")

  dirPath.ChdirIfDir()
  newDir, err := os.Getwd()
  assert.NoError(t, err)

  assert.NotEqual(t, currentDir, newDir)

  // Change back to original directory
  err = os.Chdir(currentDir)
  assert.NoError(t, err)
}

func TestPath_ChdirIfDir_NotADir(t *testing.T) {
  file, err := os.Create("testfile.txt")
  assert.NoError(t, err)
  file.Close()
  defer os.Remove("testfile.txt")

  filePath := path.New("testfile.txt")
  filePath.ChdirIfDir()

  // Ensure it doesn't change the current directory if the path is not a directory
  currentDir, err := os.Getwd()
  assert.NoError(t, err)
  newDir, err := os.Getwd()
  assert.NoError(t, err)
  assert.Equal(t, currentDir, newDir)
}

// TODO: Add tests for getFullPath function if it becomes exported or if you use an interface to mock it.
