package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/disk"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strings"
)

type Diskinfo struct {
	Device      string `json:"device"`
	Mountpoint  string `json:"mountpoint"`
	Fstype      string `json:"fstype"`
	Total       string `json:"total"`
	Free        string `json:"free"`
	Used        string `json:"used"`
	UsedPercent string `json:"usedPercent"`
}

type FileInfo struct {
	Name  string `json:"name"`
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
}

func GetDiskInfo(c *gin.Context) {
	var diskInfos []FileInfo
	parts, err := disk.Partitions(true)
	if err != nil {
		c.JSON(http.StatusOK, diskInfos)
		return
	}
	for _, part := range parts {
		var diskInfo FileInfo
		diskInfo.Name = part.Device
		diskInfo.Path = part.Device
		diskInfo.IsDir = false

		diskInfos = append(diskInfos, diskInfo)
	}
	c.JSON(http.StatusOK, diskInfos)
}

func GetFileList(c *gin.Context) {
	var FileList []FileInfo
	path := c.Query("path")
	dpath := path
	if !strings.Contains(path, "\\") {
		dpath = path + "\\"
	}
	rd, _ := ioutil.ReadDir(dpath)
	for _, fi := range rd {
		FileList = append(FileList, FileInfo{Name: fi.Name(), Path: path + "\\" + fi.Name(), IsDir: fi.IsDir()})
	}

	sort.Slice(FileList, func(i, j int) bool {
		if FileList[i].IsDir == FileList[j].IsDir {
			return FileList[i].Name < FileList[j].Name
		}

		if FileList[i].IsDir && !FileList[j].IsDir {
			return true
		} else {
			return false
		}
	})

	c.JSON(http.StatusOK, FileList)
}

func DelFile(c *gin.Context) {
	path := c.Query("path")
	if err := os.Remove(path); err != nil {
		c.String(http.StatusOK, err.Error())
	} else {
		c.String(http.StatusOK, "OK")
	}
}

func BakFile(c *gin.Context) {
	path := c.Query("path")
	var newpath string
	if index := strings.LastIndexAny(path, "."); index != -1 {
		newpath = path[:index] + "_bak" + path[index:]
	} else {
		newpath = path + "_bak"
	}
	if err := CopyFile(path, newpath); err != nil {
		c.String(http.StatusOK, err.Error())
	} else {
		c.String(http.StatusOK, "OK")
	}
}

func FileDownload(c *gin.Context) {
	path := c.Query("path")
	name := c.Query("name")
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", name))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File(path)
}

func CopyFile(source, dest string) error {
	if source == "" || dest == "" {
		return errors.New("path is null")
	}
	//??????????????????
	source_open, err := os.Open(source)
	//???????????????????????????????????????????????? defer ????????????????????????
	if err != nil {
		return err
	}
	defer source_open.Close()
	//???????????????????????? ????????????????????????????????? ????????? 644????????????????????????linux ????????????
	dest_open, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 644)
	if err != nil {
		return err
	}
	//???????????????????????????????????????????????? defer ????????????????????????
	defer dest_open.Close()
	//??????????????????
	_, err = io.Copy(dest_open, source_open)
	if err != nil {
		return err
	}

	return nil
}
