package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/yah01/CyDrive/config"
)

var (
	dbConfig config.Config

	isOnline      bool
	serverAddress string
)

func init() {
	// Parse args
	flag.BoolVar(&isOnline, "--online", false, "whether is online")
	flag.StringVar(&serverAddress, "-h", "localhost", "set the CyDrive Server address")
	flag.Parse()

	// Read DB config
	// dbConfigFile, err := ioutil.ReadFile("db_config.yaml")
	// if err != nil {
	//	panic(err)
	// }
	// if err = yaml.Unmarshal(dbConfigFile, &dbConfig); err != nil {
	//	panic(err)
	// }

	logFile, err := os.OpenFile("log", os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

	currentEnv = localEnv
}

func main() {
	dbConfig.UserStoreType = "mem"
	InitServer(dbConfig)
	RunServer()
}
