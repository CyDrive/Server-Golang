package main

import (
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/yah01/CyDrive/config"
)

var (
	dbConfig config.Config
	log      *logrus.Logger

	isOnline      bool
	serverAddress string
)

func init() {
	// Parse args
	flag.BoolVar(&isOnline, "--online", false, "whether is online")
	flag.StringVar(&serverAddress, "-h", "localhost", "set the CyDrive Server address")
	flag.Parse()

	// Read DB config
	//dbConfigFile, err := ioutil.ReadFile("db_config.yaml")
	//if err != nil {
	//	panic(err)
	//}
	//if err = yaml.Unmarshal(dbConfigFile, &dbConfig); err != nil {
	//	panic(err)
	//}

	log = logrus.New()
	logFile,err := os.OpenFile("log", os.O_CREATE|os.O_APPEND, 0777)
	if err!=nil {
		panic(err)
	}
	log.Out = logFile
	log.SetNoLock()
	log.SetReportCaller(true)
	log.SetFormatter(&logrus.JSONFormatter{})
}

func main() {
	dbConfig.UserStoreType = "mem"
	InitServer(dbConfig)
	RunServer()
}
