package router

import (
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"go-windows-monitor/api"
	"go-windows-monitor/assets"
)

// 初始化总路由

func Routers() *gin.Engine {
	assetInfo := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	var Router = gin.Default()
	//Router.Use(LoadTls())  // 打开就能玩https了
	// 跨域
	StaticFs := assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, AssetInfo: assetInfo, Prefix: "dist/static", Fallback: "index.html"}
	CmdFs := assetfs.AssetFS{Asset: assets.Asset, AssetDir: assets.AssetDir, AssetInfo: assetInfo, Prefix: "dist/cmd", Fallback: "index.html"}
	Router.Use(Cors())
	Router.StaticFS("/static", &StaticFs)
	Router.StaticFS("/cmd", &CmdFs)

	// 方便统一添加路由组前缀 多服务器上线使用
	ApiGroup := Router.Group("")
	InitHomeRouter(ApiGroup)
	return Router
}

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, x-token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

func InitHomeRouter(Router *gin.RouterGroup) {
	ApiRouter := Router.Group("")
	{
		ApiRouter.GET("", Home)                               // 创建Api
		ApiRouter.GET("GetProcesslist", api.GetProcesslist)   // 创建Api
		ApiRouter.GET("GetSystemInfo", api.GetSystemInfo)     // 创建Api
		ApiRouter.GET("GetCpuMemory", api.GetCpuMemory)       // 创建Api
		ApiRouter.GET("GetServicelist", api.GetServicelist)   // 创建Api
		ApiRouter.GET("GetOsInfo", api.GetOsInfo)             // 创建Api
		ApiRouter.GET("GetIOCounters", api.GetIOCounters)     // 创建Api
		ApiRouter.GET("GetExecutePrint", api.GetExecutePrint) // 创建Api
		ApiRouter.GET("SetExecute", api.SetExecute)           // 创建Api
		ApiRouter.POST("KillProcess", api.KillProcess)

		//cmd
		ApiRouter.GET("pty", api.Pty) // 创建Api

		//Flie
		ApiRouter.GET("CopyFile", api.BakFile)
		ApiRouter.GET("GetFileList", api.GetFileList)
		ApiRouter.GET("GetDiskinfo", api.GetDiskInfo)   // 创建Api
		ApiRouter.GET("DelFile", api.DelFile)           // 创建Api
		ApiRouter.GET("FileDownload", api.FileDownload) // 创建Api
	}
}

func Home(c *gin.Context) {
	c.Writer.WriteHeader(200)
	indexHtml, _ := assets.Asset("dist/index.html")
	_, _ = c.Writer.Write(indexHtml)
	c.Writer.Header().Add("Accept", "text/html")
	c.Writer.Flush()
}
