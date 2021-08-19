package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kbinani/win"
	"go-windows-monitor/cmd"
	"go-windows-monitor/systeminfo"
	"net/http"
	"strconv"
	"unsafe"
)

func SetExecute(c *gin.Context) {
	Execute := c.Query("exec")
	Guid := c.Query("Guid")

	if Execute == "1" {
		fmt.Println(getNetworkParams())
		Execute = "www.icafe8.com"
	}

	if Execute == "2" {
		Execute = "www.qq.com"
	}

	if info := systeminfo.CmdMgr.ExecCommand(Guid, Execute); info == "error" {
		c.JSON(http.StatusOK, "err")

	} else {
		c.JSON(http.StatusOK, info)
	}
}

func GetExecutePrint(c *gin.Context) {
	Count, _ := strconv.Atoi(c.Query("Count"))
	Guid := c.Query("Guid")

	s := systeminfo.CmdMgr.GetResp(Guid, Count)
	c.JSON(http.StatusOK, s)
}

func getNetworkParams() []string {
	dns := make([]string, 0)
	info := win.FIXED_INFO_W2KSP1{}
	size := uint32(unsafe.Sizeof(info))
	r := win.GetNetworkParams(&info, &size)
	if r == 0 {
		for ai := &info.DnsServerList; ai != nil; ai = ai.Next {
			d := fmt.Sprintf("%v.%v.%v.%v", ai.Context&0xFF, (ai.Context>>8)&0xFF, (ai.Context>>16)&0xFF, (ai.Context>>24)&0xFF)
			dns = append(dns, d)
		}
	} else if r == win.ValueOverflow {
		newBuffers := make([]byte, size)
		netParams := (win.PFIXED_INFO)(unsafe.Pointer(&newBuffers[0]))
		win.GetNetworkParams(netParams, &size)
		for ai := &netParams.DnsServerList; ai != nil; ai = ai.Next {
			d := fmt.Sprintf("%v.%v.%v.%v", ai.Context&0xFF, (ai.Context>>8)&0xFF, (ai.Context>>16)&0xFF, (ai.Context>>24)&0xFF)
			dns = append(dns, d)
		}
	}
	return dns
}

func Pty(c *gin.Context) {
	cmd.PtyHandler(c.Writer, c.Request)
}
