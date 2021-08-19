package main

import (
	"github.com/pkg/errors"
	"go-windows-monitor/utils/windows"
	"strconv"
)

type Config struct {
	BSAddr string     `json:"bs"`
	Addr   string     `json:"addr"`
	BarID  string     `json:"bar_id"`
	Log    *LogConfig `json:"log"`
}

type LogConfig struct {
	LogLevel    string `json:"logLevel"`    // 日志级别，支持：off/trace/debug/info/warn/error/panic/fatal
	ReserveDays int    `json:"reserveDays"` // 日志文件保留天数
	MaxSize     int    `json:"maxSize"`     // 日志文件最大大小，单位：MB
	PrintScreen bool   `json:"printScreen"` // 是否打印至标准输出
}

var (
	defaultPrintScreen    = true
	defaultLogLevel       = "info"
	defaultLogReserveDays = 3
	defaultLogFileMaxSize = 100
)

func NewConfig() *Config {
	return &Config{
		Log: &LogConfig{
			LogLevel:    defaultLogLevel,
			ReserveDays: defaultLogReserveDays,
			MaxSize:     defaultLogFileMaxSize,
			PrintScreen: defaultPrintScreen,
		},
	}
}

func (c *Config) Load() error {
	/*content, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return errors.Wrap(err, "read config file error")
	}

	if err := json.Unmarshal(content, c); err != nil {
		return errors.Wrap(err, "unmarshal config error")
	}*/
	path := "./config.ini"
	c.Addr = windows.GetPrivateProfileString(path, "config", "addr")
	if c.Addr == "" {
		return errors.New("Addr is empty")
	}
	c.BarID = windows.GetPrivateProfileString(path, "config", "barid")
	c.logLoad(path)
	return nil
}

func (c *Config) logLoad(path string) {
	LogLevel := windows.GetPrivateProfileString(path, "log", "level")
	ReserveDays := windows.GetPrivateProfileString(path, "log", "reservedays")
	MaxSize := windows.GetPrivateProfileString(path, "log", "maxsize")
	PrintScreen := windows.GetPrivateProfileString(path, "log", "printscreen")
	if LogLevel != "" {
		c.Log.LogLevel = LogLevel
	}
	if ReserveDays != "" {
		c.Log.ReserveDays, _ = strconv.Atoi(ReserveDays)
	}
	if MaxSize != "" {
		c.Log.MaxSize, _ = strconv.Atoi(MaxSize)
	}
	if PrintScreen == "false" {
		c.Log.PrintScreen = false
	}
}
