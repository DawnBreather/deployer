package main

import (
  "github.com/pterm/pterm"
  "github.com/sirupsen/logrus"
  "time"
)

//deployerctl set /environments/%s --value file://demo.adhereit.cohero-health.com.json
//deployerctl status
//deployerctl set /environments/%s/tier/credentials --value file://credentials.json
//deployerctl set /environments/%s/tier/credentials --value file://credentials.yaml
//deployerctl set /environments/%s/tier/credentials --value file://credentials.yml
//deployerctl set /environments/%s/tier/credentials/aws --value file://aws.credentials.yaml
//deployerctl set /environments/%s/tier/artifacts --value file://artifacts.yaml
// TODO: rename health-checks for health_checks
//deployerctl set /environments/%s/tier/health-checks --value file://artifacts.yaml
//deployerctl get /environments/%s/tier/artifacts -o yaml
//deployerctl get /environments/%s/tier/artifacts -o json
//deployerctl remove /environments/%s/tier/artifacts
//deployerctl inventory prune
//deployerctl deploy
//deployerctl extract-logs all | DESKTOP-JGBQ13D --stdout --stderr
//deployerctl check locations

func main() {
  logrus.SetFormatter(&logrus.TextFormatter{
    DisableQuote: true,
    ForceColors:  true,
  })

  //var d = deployment.Deployment{
  //	Name: os.Getenv(deployment.ABZ_DEPLOYER_AGENT_ENVIRONMENT_ID),
  //}

  //d.StartListeningForTierConfiguration()

  //for !d.ConfigurationInitializedFromFirebaseSource {
  //	time.Sleep(100 * time.Millisecond)
  //}

  //deployment.StartListeningForFirebaseValueByPath(fmt.Sprintf("environments/%s/status", d.Name()), func(value interface{}) {
  //	jsonBytes, _ := json.Marshal(value)
  //	var status deployment.Status
  //	json.Unmarshal(jsonBytes, &status)
  //	//jsonPrettied := pretty.Color(pretty.Pretty(jsonBytes), pretty.TerminalStyle)
  //	jsonPrettied := pretty.Pretty(jsonBytes)
  //
  //	//fmt.Printf("%#v\n", status)
  //	fmt.Println(string(jsonPrettied))
  //	os.Exit(0)
  //})
  //
  //for {
  //	time.Sleep(1 * time.Second)
  //}

  pterm.Info.Println("The previous text will stay in place, while the area updates.")
  pterm.Print("\n\n") // Add two new lines as spacer.

  area, _ := pterm.DefaultArea.WithCenter().Start() // Start the Area printer, with the Center option.
  for i := 0; i < 10; i++ {
    str, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString(time.Now().Format("15:04:05"))).Srender() // Save current time in str.
    area.Update(str)                                                                                                // Update Area contents.
    time.Sleep(time.Second)
  }
  area.Stop()

}

func inventoryPrune() {

}
