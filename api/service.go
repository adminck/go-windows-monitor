package api

import (
	"github.com/gin-gonic/gin"
	"go-windows-monitor/systeminfo"
	"net/http"
)

type Win32Service struct {
	Name        string `json:"name"`
	Caption     string `json:"caption"`
	State       string `json:"state"`
	Description string `json:"description"`
	StartMode   string `json:"startMode"`
}

func GetServicelist(c *gin.Context) {
	var resp []Win32Service
	//size,_ := strconv.Atoi(c.DefaultQuery("size","10"))
	//page,_ := strconv.Atoi(c.DefaultQuery("page", "1"))
	//search := c.DefaultQuery("search", "")
	s := systeminfo.WinServiceMgr.Lookup()
	for _, v := range s {
		resp = append(resp, Win32Service{Name: v.Name,
			Caption:     v.Caption,
			State:       v.State,
			Description: v.Description,
			StartMode:   v.StartMode})
	}
	c.JSON(http.StatusOK, resp)
}
