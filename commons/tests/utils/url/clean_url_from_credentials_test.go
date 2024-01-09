package url_test

import (
  "deployer/commons/utils/url"
  "testing"
)

func TestCleanUrlFromCredentials(t *testing.T) {
  tests := []struct {
    name     string
    input    string
    expected string
  }{
    {
      name:     "HTTP with credentials",
      input:    "http://user:pass@somehost.com/path",
      expected: "http://somehost.com/path",
    },
    {
      name:     "HTTPS with credentials",
      input:    "https://user:pass@somehost.com/path",
      expected: "https://somehost.com/path",
    },
    {
      name:     "Git protocol with credentials",
      input:    "git://user:pass@somehost.com/path",
      expected: "git://somehost.com/path",
    },
    {
      name:     "SSH with credentials",
      input:    "ssh://user:pass@somehost.com/path",
      expected: "ssh://somehost.com/path",
    },
    {
      name:     "URL without credentials",
      input:    "http://somehost.com/path",
      expected: "http://somehost.com/path",
    },
    {
      name:     "URL with special characters in credentials",
      input:    "http://user:pa@ss@somehost.com/path",
      expected: "http://somehost.com/path",
    },
    {
      name:     "URL with empty credentials",
      input:    "http://:@somehost.com/path",
      expected: "http://somehost.com/path",
    },
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      output := url.CleanUrlFromCredentials(tt.input)
      if output != tt.expected {
        t.Errorf("Expected %s, got %s", tt.expected, output)
      }
    })
  }
}
