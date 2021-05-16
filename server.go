package main

import (
	"encoding/gob"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yah01/CyDrive/config"
	. "github.com/yah01/CyDrive/consts"
	"github.com/yah01/CyDrive/model"
	"github.com/yah01/CyDrive/store"
)

var (
	userStore store.UserStore
	router    *gin.Engine
)

func InitServer(config config.Config) {

	if config.UserStoreType == "mem" {
		userStore = store.NewMemStore("user_data/user.json")
	}

	router = gin.Default()
	gob.Register(&model.User{})
	gob.Register(time.Time{})
}

func RunServer() {
	memStore := memstore.NewStore([]byte("ProjectMili"))

	router.Use(sessions.SessionsMany([]string{"user"}, memStore))
	router.Use(LoginAuth(router))
	// router.Use(SetFileInfo())

	router.POST("/login", LoginHandle)
	router.GET("/list/*path", ListHandle)

	router.GET("/file_info/*path", GetFileInfoHandle)
	router.PUT("/file_info/*path", PutFileInfoHandle)

	router.GET("/file/*path", DownloadHandle)
	router.PUT("/file/*path", UploadHandle)

	go ftm.Listen()
	router.Run(ListenPort)
}
