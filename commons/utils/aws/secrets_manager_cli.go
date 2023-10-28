package aws

import (
  "fmt"
  "github.com/iancoleman/strcase"
  "github.com/urfave/cli/v2"
  "strings"
)

const (
  SOURCE_SECRETS = "source-secrets"
  SET_ENVIRONMENT_VARIABLES = "set-environment-variables"
)

// secrets aws-secrets-manager extract --source-secrets ${SOURCE_SECRETS} --set-environment-variabels

func GetSecretsManagerCliCommands()[]*cli.Command{
  return []*cli.Command{
    {
      Name: "aws-secrets-manager",
      Usage: "working with aws-secrets-manager",
      Subcommands: []*cli.Command{
        {
          Name: "extract",
          Usage: "extracts secrets",
          Action: extract,
          Flags: []cli.Flag{
            &cli.StringFlag{
              Name:    SOURCE_SECRETS,
              Aliases: []string{"c"},
              Usage:   "Provide name of secrets in format  `some/secret1,some/secret2,some/secretN`",
              EnvVars: []string{strcase.ToScreamingSnake(SOURCE_SECRETS)},
              Required: true,
            },
            &cli.BoolFlag{
              Name:    SET_ENVIRONMENT_VARIABLES,
              Aliases: []string{"s"},
              Usage:   "Will set variables as environment variables",
              EnvVars: []string{strcase.ToScreamingSnake(SET_ENVIRONMENT_VARIABLES)},
              Required: false,
            },
          },
        },
      },
    },
  }
}

func extract(c *cli.Context) error{
  result := ""
  secretsNames := c.String(SOURCE_SECRETS)
  for _, secretName := range strings.Split(secretsNames, ","){
    secretContent, ok := SecretsManager.GetSecret(secretName)
    if ! ok {
      secretContent = "# Error reading secret's content"
    }
    result += fmt.Sprintf("# From secret { %s }\n%s\n\n", secretName, secretContent)
  }
  fmt.Printf("%s", result)

  //if er := setEnvironmentVariablesIfEnabled(c.Bool(SET_ENVIRONMENT_VARIABLES), result); er != nil {
  //  return er
  //}

  return nil
}

//func setEnvironmentVariablesIfEnabled(shouldSetEnvironmentVariable bool, result string) error{
//  if shouldSetEnvironmentVariable {
//    for _, r := range result {
//
//    }
//    if er := os.Setenv(secretName, secretValue); er != nil  {
//      return er
//    }
//  }
//  return nil
//}