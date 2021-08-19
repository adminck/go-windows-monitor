package main

import (
	"fmt"
	"go-windows-monitor/router"
	"go-windows-monitor/utils/comm"
	"go-windows-monitor/utils/log"
	"go-windows-monitor/utils/windows"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"runtime/debug"
	"syscall"
	"time"
)

func init() {
	//flag.StringVar(&configFile, "c", "./config.json", "config file path")
	//flag.Parse()
}

var config = NewConfig()

func main() {
	//change working directory
	os.Chdir(filepath.Dir(comm.GetCurrentProcessPath()))
	//go func() {
	//	http.ListenAndServe(":6789", nil)
	//}()

	go func() {
		for {
			debug.FreeOSMemory()
			time.Sleep(time.Minute)
		}
	}()

	initLog(config.Log)
	if err := config.Load(); err != nil {
		log.Error("config.json read failed")
		return
	}

	log.Infoln("###################################################")
	log.Infof("config.bs: %s", config.BSAddr)
	log.Infof("config.Log.LogLevel: %s", config.Log.LogLevel)
	log.Infof("config.Log.ReserveDays: %d", config.Log.ReserveDays)
	log.Infof("config.Log.MaxSize: %d", config.Log.MaxSize)
	log.Infoln("###################################################")

	//go ReportHardwareInfo()

	keepAlive()

	if err := RunHttpServer(config.Addr); err != nil {
		log.Error("HttpServer start failed err: %s", err.Error())
		return
	}

	//进程退出
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Infoln("recv signal interrupt, exit...")
}

func RunHttpServer(HttpAddress string) error {
	Router := router.Routers()
	hs := &http.Server{
		Addr:           ":" + HttpAddress,
		Handler:        Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	err := hs.ListenAndServe()
	return err
}

func initLog(cfg *LogConfig) {
	level := log.StringToLevel(cfg.LogLevel)
	log.SetLevel(level)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var (
		writer io.Writer
		err    error
	)
	writer, err = log.NewFileWriter(
		path.Join(filepath.Dir(comm.GetCurrentProcessPath()), "log", "barmonitor.log"),
		log.ReserveDays(cfg.ReserveDays),
		log.RotateByDaily(true),
		log.LogFileMaxSize(cfg.MaxSize),
	)
	if err != nil {
		log.Errorln("create file writer error:", err)
		return
	}

	if cfg.PrintScreen {
		writer = io.MultiWriter(writer, os.Stdout)
	}

	log.SetOutput(writer)
}

func keepAlive() {
	s := fmt.Sprintf("Global\\%s.tsa", comm.GetCurrentProcessName())
	//fmt.Println(s)
	log.Info("create event=", s)
	handler, err := windows.CreateEvent(nil, 0, 0, s)
	if nil != err {
		log.Error("CreateEvent failed")
		return
	}
	go func(syscall.Handle) {
		for {
			time.Sleep(1 * time.Second)
			windows.SetEvent(handler)
		}
	}(handler)
}
