package deployment

import "github.com/sirupsen/logrus"

// TODO: move to commons
func setFirebaseValueByPath(refPath string, value interface{}) error {
  ref, err := f.Ref(refPath)
  if err != nil {
    logrus.Errorf("[E] Referring to { %s } status in firebase: %v", refPath, err)
    return err
  }

  err = ref.Set(value)
  if err != nil {
    logrus.Errorf("[E] Setting { %s } in firebase: %v", refPath, err)
    return err
  }

  return nil
}

// TODO: move to commons
func getFirebaseValueByPath(refPath string) (value interface{}, err error) {
  ref, err := f.Ref(refPath)
  if err != nil {
    logrus.Errorf("[E] Referring to { %s } status in firebase: %v", refPath, err)
    return nil, err
  }

  err = ref.Value(&value)
  if err != nil {
    logrus.Errorf("[E] Setting { %s } in firebase: %v", refPath, err)
    return nil, err
  }

  return
}

// TODO: move to commons
func removeFirebaseValueByPath(refPath string) (err error) {
  ref, err := f.Ref(refPath)
  if err != nil {
    logrus.Errorf("[E] Referring to { %s } status in firebase: %v", refPath, err)
    return err
  }

  err = ref.Set(nil)
  if err != nil {
    logrus.Errorf("[E] Setting { %s } in firebase: %v", refPath, err)
    return err
  }

  return
}
