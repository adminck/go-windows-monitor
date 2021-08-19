package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gopsutilnet "github.com/shirou/gopsutil/net"
	"go-windows-monitor/systeminfo"
	"net/http"
	"strconv"
)

type SystemInfo struct {
	Hostname        string   `json:"机器名"`
	Platform        string   `json:"系统"`
	PlatformVersion string   `json:"系统版本"`
	VideoChip       []string `json:"显卡"`
	ModelName       string   `json:"cpu型号"`
	Mac             string   `json:"MAC地址"`
	MemoryTotal     string   `json:"内存总大小"`
}

type NetIOCounters struct {
	BytesSent uint64 `json:"bytesSent"` // number of bytes sent
	BytesRecv uint64 `json:"bytesRecv"` // number of bytes received
}

func GetSystemInfo(c *gin.Context) {
	var systemInfo SystemInfo
	systemInfo.Hostname = systeminfo.HardwareInfoMgr.Hostname
	systemInfo.Platform = systeminfo.HardwareInfoMgr.Platform
	systemInfo.PlatformVersion = systeminfo.HardwareInfoMgr.PlatformVersion

	for _, v := range systeminfo.HardwareInfoMgr.VideoChip {
		systemInfo.VideoChip = append(systemInfo.VideoChip, fmt.Sprintf("%s(显存：%dGB)", v.Name, v.AdapterRAM/1024/1024/1024))
	}

	fmt.Println(systeminfo.HardwareInfoMgr.Processor)
	for _, v := range systeminfo.HardwareInfoMgr.Processor {
		systemInfo.ModelName = v.Name
		continue
	}

	for _, v := range systeminfo.HardwareInfoMgr.NetAdapter {
		if v.NetConnectionID == "本地连接" {
			systemInfo.Mac = v.MACAddress
			continue
		}
	}

	Total := 0
	for _, v := range systeminfo.HardwareInfoMgr.Memory {
		s, _ := strconv.Atoi(v.Capacity)
		Total += s
	}

	systemInfo.MemoryTotal = fmt.Sprintf("%dMB", Total/1024/1024)

	c.JSON(http.StatusOK, systemInfo)
}

func GetCpuMemory(c *gin.Context) {
	c.JSON(http.StatusOK, systeminfo.WinPercentMgr.Lookup())
}

func GetOsInfo(c *gin.Context) {
	c.JSON(http.StatusOK, systeminfo.HardwareInfoMgr)
}

func GetIOCounters(c *gin.Context) {
	var IOCounters NetIOCounters
	if IO, err := gopsutilnet.IOCounters(false); err != nil {
		IOCounters.BytesRecv = 0
		IOCounters.BytesSent = 0
	} else {
		IOCounters.BytesRecv = IO[0].BytesRecv
		IOCounters.BytesSent = IO[0].BytesSent
	}
	c.JSON(http.StatusOK, IOCounters)
}
