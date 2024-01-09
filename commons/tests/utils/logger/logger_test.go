package logger_test

import (
  . "deployer/commons/utils/logger"
  "github.com/sirupsen/logrus"
  "github.com/stretchr/testify/assert"
  "testing"
)

func TestLogger_NewLogger(t *testing.T) {
  logger := New()
  assert.NotNil(t, logger, "Expected logger to be non-nil")
  assert.IsType(t, &logrus.Logger{}, logger, "Expected logger to be of type *logrus.Logger")

  textFormatter, ok := logger.Formatter.(*logrus.TextFormatter)
  assert.True(t, ok, "Expected logger formatter to be of type *logrus.TextFormatter")
  assert.True(t, textFormatter.ForceColors, "Expected ForceColors to be true")
  assert.False(t, textFormatter.DisableColors, "Expected DisableColors to be false")
  assert.False(t, textFormatter.ForceQuote, "Expected ForceQuote to be false")
  assert.False(t, textFormatter.DisableQuote, "Expected DisableQuote to be false")
  assert.False(t, textFormatter.EnvironmentOverrideColors, "Expected EnvironmentOverrideColors to be false")
  assert.False(t, textFormatter.DisableTimestamp, "Expected DisableTimestamp to be false")
  assert.True(t, textFormatter.FullTimestamp, "Expected FullTimestamp to be true")
  assert.Equal(t, "20060102150405", textFormatter.TimestampFormat, "Expected TimestampFormat to be '20060102150405'")
  assert.False(t, textFormatter.DisableSorting, "Expected DisableSorting to be false")
  assert.Nil(t, textFormatter.SortingFunc, "Expected SortingFunc to be nil")
  assert.False(t, textFormatter.DisableLevelTruncation, "Expected DisableLevelTruncation to be false")
  assert.False(t, textFormatter.PadLevelText, "Expected PadLevelText to be false")
  assert.False(t, textFormatter.QuoteEmptyFields, "Expected QuoteEmptyFields to be false")
  assert.Nil(t, textFormatter.FieldMap, "Expected FieldMap to be nil")
  assert.Nil(t, textFormatter.CallerPrettyfier, "Expected CallerPrettyfier to be nil")
}
