/**
 * go build -ldflags "-X main.version=0.1" app.go
 */

package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"
  "os/signal"
  "syscall"
)

type Config struct {
  Db          string `json:"app.db"`
  Port        string `json:"app.port"`
  LogDir      string `json:"log.dir"`
  LogFile     string `json:"log.file"`
  StaticDir   string `json:"static.dir"`
  TemplateDir string `json:"template.dir"`
}

var (
  version string
)

func init() {

}

func startup() {
  log.Println("done.")
}

func shutdown() {
  log.Println("done.")
}

func LoadConfig(file string) Config {
  var config Config
  configFile, err := os.Open(file)
  if err != nil {
    panic(err)
  }
  defer configFile.Close()
  jsonParser := json.NewDecoder(configFile)
  jsonParser.Decode(&config)
  return config
}

func main() {

  configFile := "etc/config.json" // relative to current working directory
  if len(os.Args) > 1 {
    configFile = os.Args[1]
  }
  config := LoadConfig(configFile)

  logfile := config.LogDir + string(os.PathSeparator) + config.LogFile
  logger, err := os.OpenFile(logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0644)
  if err != nil {
    log.Fatal("error opening file: %v", err)
  }
  defer logger.Close()
  log.SetOutput(logger)
  log.SetPrefix(fmt.Sprintf("[%d] ", os.Getpid()))

  log.Printf("Starting {{ app }} %s...", version)
  startup()

  fs := http.FileServer(http.Dir(config.StaticDir))
  http.Handle("/", fs)

  c := make(chan os.Signal)
  signal.Notify(c, os.Interrupt, syscall.SIGTERM)
  go func() {
    <-c
    log.Printf("Exiting {{ app }} %s...", version)
    shutdown()
    os.Exit(1)
  }()

  dir, err := os.Getwd()
  if err != nil {
    log.Fatal(err)
  }
  fmt.Printf("starting {{ app }} version %s, listening on port %s, see %s for details...\n", version, config.Port, dir + "/" + logfile)
  log.Fatal(http.ListenAndServe(config.Port, nil))

}
