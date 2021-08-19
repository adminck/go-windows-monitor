package systeminfo

var CmdMgr *Cmd
var WinServiceMgr *ServiceMgr
var WinPercentMgr *PercentMgr
var ProcExMgr *ProcessMgr
var HardwareInfoMgr HardwareInfo
var ProcessorCoreNum int

func init() {
	CmdMgr = NewCmdMgr()
	WinServiceMgr = NewServiceMgr()
	WinPercentMgr = NewPercentMgr()
	ProcExMgr = NewProcessMgr()
	HardwareInfoMgr.Run()
}
