package main

import(
  "os"
  log "github.com/Sirupsen/logrus"
  "github.com/voxelbrain/goptions"
  "github.com/yosmudge/pagerbot/config"
  "github.com/yosmudge/pagerbot/updater"
)

type options struct {
  Verbose   bool            `goptions:"-v, --verbose, description='Log verbosely'"`
  Help      goptions.Help   `goptions:"-h, --help, description='Show help'"`
  Config    string          `goptions:"-c, --config, description='Config Yaml file to use'"`
}

func main() {
  parsedOptions := options{}

  parsedOptions.Config = "./config.yml"

  goptions.ParseAndFail(&parsedOptions)

  if parsedOptions.Verbose{
    log.SetLevel(log.DebugLevel)
  } else {
    log.SetLevel(log.InfoLevel)
  }

  log.SetFormatter(&log.TextFormatter{FullTimestamp:true})

  log.Debug("Logging verbosely!")

  err := config.Load(parsedOptions.Config)
  if err == nil {
    err = config.Config.Validate()
  }

  if err != nil{
    log.WithFields(log.Fields{
      "configFile": parsedOptions.Config,
      "error": err,
    }).Error("Could not load config file")
    os.Exit(1)
  }

  u, err := updater.New()
  if err != nil {
    log.WithFields(log.Fields{
      "error": err,
    }).Error("Could not start updater")
    os.Exit(1)
  }

  u.Start()
  u.Wg.Wait()
}
