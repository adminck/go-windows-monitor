package api

import (
	"github.com/gin-gonic/gin"
	"go-windows-monitor/systeminfo"
	"net/http"
	"strconv"
)

type ProcessInfo struct {
	Pid         int32  `json:"pid"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Memory      uint64 `json:"memory"`
	Cpu         int    `json:"cpu"`
	Username    string `json:"username"`
	Uid         uint32 `json:"sessionid"`
	ThreadCount uint32 `json:"thread_count"`
	HandleCount uint32 `json:"handle_count"`
}

func GetProcesslist(c *gin.Context) {
	var resp []ProcessInfo
	//size,_ := strconv.Atoi(c.DefaultQuery("size","10"))
	//page,_ := strconv.Atoi(c.DefaultQuery("page", "1"))
	processes, _ := systeminfo.ProcExMgr.Lookup(1, 999)
	for _, v := range processes {
		resp = append(resp, ProcessInfo{
			Pid:         v.Pid,
			Name:        v.Name,
			Path:        v.Path,
			Memory:      v.Memory,
			Cpu:         v.Cpu,
			Username:    v.Username,
			Uid:         v.SessionId,
			ThreadCount: v.ThreadCount,
			HandleCount: v.HandleCount})
	}
	c.JSON(http.StatusOK, resp)
}

func KillProcess(c *gin.Context) {
	pid, _ := strconv.Atoi(c.PostForm("pid"))
	c.JSON(http.StatusOK, systeminfo.ProcExMgr.KillProcess(int32(pid)))
}
