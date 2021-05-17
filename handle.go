package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	. "github.com/yah01/CyDrive/consts"
	"github.com/yah01/CyDrive/env"
	"github.com/yah01/CyDrive/model"
	"github.com/yah01/CyDrive/utils"
)

var (
	localEnv   = env.NewLocalEnv()
	currentEnv env.Env
)

func LoginHandle(c *gin.Context) {
	username, ok := c.GetPostForm("username")
	if !ok {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "no user name",
			Data:    nil,
		})
		return
	}

	password, ok := c.GetPostForm("password")
	if !ok {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "no password",
			Data:    nil,
		})
		return
	}

	user := userStore.GetUserByName(username)
	if user == nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "no such user",
			Data:    nil,
		})
		return
	}
	if utils.PasswordHash(user.Password) != password {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusAuthError,
			Message: "user name or password not correct",
			Data:    nil,
		})
		return
	}

	userSession := sessions.DefaultMany(c, "user")

	userSession.Set("userStruct", &user)
	userSession.Set("expire", time.Now().Add(time.Hour*12))
	err := userSession.Save()
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	userJson, err := json.Marshal(user.SafeUser)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: "Welcome to CyDrive!",
		Data:    string(userJson),
	})
}

func ListHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	path := c.Param("path")

	path = strings.Trim(path, "/")
	absPath := strings.Join([]string{user.RootDir, path}, "/")

	fileList, err := currentEnv.ReadDir(absPath)
	for i := range fileList {
		fileList[i].FilePath = strings.ReplaceAll(fileList[i].FilePath, "\\", "/")
		fileList[i].FilePath = strings.TrimPrefix(fileList[i].FilePath, user.RootDir)
	}
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	fileListJson, err := json.Marshal(fileList)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: "list done",
		Data:    string(fileListJson),
	})
}

func GetFileInfoHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	filePath := c.Param("path")
	filePath = strings.Trim(filePath, "/")
	absFilePath := filepath.Join(user.RootDir, filePath)

	fileInfo, err := currentEnv.Stat(absFilePath)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: "get file info done",
		Data:    fileInfo,
	})
}

func PutFileInfoHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	filePath := c.Param("path")
	filePath = strings.Trim(filePath, "/")
	absFilePath := filepath.Join(user.RootDir, filePath)

	_, err := currentEnv.Stat(absFilePath)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	len := c.Request.ContentLength
	fileInfoJson := make([]byte, len)
	c.Request.Body.Read(fileInfoJson)

	fileInfo := model.FileInfo{}
	if err := json.Unmarshal(fileInfoJson, &fileInfo); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: "error when parsing file info",
			Data:    nil,
		})
		return
	}

	openFile, err := currentEnv.OpenFile(absFilePath, os.O_RDWR, os.FileMode(fileInfo.FileMode))
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}
	defer openFile.Close()

	if err = openFile.Chmod(os.FileMode(fileInfo.FileMode)); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if err = currentEnv.Chtimes(absFilePath, time.Now(), time.Unix(fileInfo.ModifyTime, 0)); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
			Data:    nil,
		})

		return
	}

	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: "put file info done",
		Data:    nil,
	})
}

func DownloadHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	// relative path
	filePath := c.Param("path")

	// absolute filepath
	filePath = strings.Join([]string{user.RootDir, filePath}, "/")
	fileInfo, _ := currentEnv.Stat(filePath)

	if fileInfo.IsDir() {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: "not a file",
		})
		return
	}

	// range
	var begin, _ int64 = 0, fileInfo.Size() - 1
	bytesRange := c.GetHeader("Range")
	if len(bytesRange) > 0 {
		begin, _ = utils.UnpackRange(bytesRange)
	}

	taskId := ftm.AddTask(&fileInfo, user, DownloadTaskType, begin)

	uFileInfo := fileInfo
	uFileInfo.FilePath, _ = filepath.Rel(user.RootDir, uFileInfo.FilePath)
	uFileInfo.FilePath = strings.ReplaceAll(uFileInfo.FilePath, "\\", "/")
	jsonBytes, err := json.Marshal(uFileInfo)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: fmt.Sprintf("serialize file info failed: %+v", err),
		})
		return
	}

	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: fmt.Sprint(taskId),
		Data:    string(jsonBytes),
	})
}

func UploadHandle(c *gin.Context) {
	userI, _ := c.Get("user")
	user := userI.(*model.User)

	filePath := c.Param("path")

	filePath = strings.Join([]string{user.RootDir, filePath}, "/")
	fileDir := filepath.Dir(filePath)
	if err := currentEnv.MkdirAll(fileDir, 0666); err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: err.Error(),
		})
		return
	}

	var fileInfo model.FileInfo

	jsonBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusIoError,
			Message: fmt.Sprintf("read request body error: %+v", err),
		})
		return
	}
	if len(jsonBytes) == 0 || json.Unmarshal(jsonBytes, &fileInfo) != nil {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusInternalError,
			Message: fmt.Sprintf("need file info"),
		})
		return
	}

	// Check file size
	if fileInfo.Size_ > FileSizeLimit {
		c.JSON(http.StatusOK, model.Resp{
			Status:  StatusFileTooLargeError,
			Message: "file is too large",
		})
		return
	}

	// Check user storage capability
	if fileInfo.Size_ > user.Cap {
		c.JSON(http.StatusOK, model.Resp{
			Status: StatusFileTooLargeError,
			Message: fmt.Sprintf("no enough capability, free storage: %vMiB, and size of the file: %vMiB",
				(user.Cap-user.Usage)>>20, fileInfo.Size_>>20), // Convert Byte to MB
		})
		return
	}

	taskId := ftm.AddTask(&fileInfo, user, UploadTaskType, fileInfo.Size_)

	c.JSON(http.StatusOK, model.Resp{
		Status:  StatusOk,
		Message: fmt.Sprint(taskId),
	})

	// if err = saveFile.Chmod(os.FileMode(fileInfo.FileMode)); err != nil {
	// 	c.JSON(http.StatusOK, model.Resp{
	// 		Status:  StatusInternalError,
	// 		Message: err.Error(),
	// 		Data:    nil,
	// 	})
	// 	return
	// }
	//
	// saveFile.Close()
	//
	// if err = currentEnv.Chtimes(filePath, time.Now(), time.Unix(fileInfo.ModifyTime, 0)); err != nil {
	// 	c.JSON(http.StatusOK, model.Resp{
	// 		Status:  StatusInternalError,
	// 		Message: err.Error(),
	// 		Data:    nil,
	// 	})
	// 	return
	// }
}
