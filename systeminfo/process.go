package systeminfo

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/process"
	"go-windows-monitor/utils/log"
	"golang.org/x/sys/windows"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	modkernel32                    = windows.NewLazySystemDLL("kernel32.dll")
	procQueryFullProcessImageNameW = modkernel32.NewProc("QueryFullProcessImageNameW")
	procGetProcessHandleCount      = modkernel32.NewProc("GetProcessHandleCount")
	procProcessIdToSessionId       = modkernel32.NewProc("ProcessIdToSessionId")

	modpsapi                     = windows.NewLazySystemDLL("psapi.dll")
	procGetProcessMemoryInfo     = modpsapi.NewProc("GetProcessMemoryInfo")
	procGetProcessImageFileNameW = modpsapi.NewProc("GetProcessImageFileNameW")
)

type ProcessInfo struct {
	Pid              int32
	Name             string
	Path             string
	Memory           uint64
	Cpu              int
	Username         string
	SessionId        uint32
	ThreadCount      uint32
	HandleCount      uint32
	lastTime         uint64
	dwElepsedTime    time.Duration
	ProcessorCoreNum int
	ParentProcessID  uint32
}

func getProcessHandleCount(h windows.Handle) (uint32, error) {
	var HandleCount uint32
	r1, _, e1 := syscall.Syscall(procGetProcessHandleCount.Addr(), 2, uintptr(h), uintptr(unsafe.Pointer(&HandleCount)), 0)
	if r1 == 0 {
		if e1 != 0 {
			return 0, error(e1)
		} else {
			return 0, syscall.EINVAL
		}
	}

	//ret.RSS = uint64(mem.WorkingSetSize)
	//ret.VMS = uint64(mem.PagefileUsage)

	return HandleCount, nil
}

func getModuleFileName(module windows.Handle) (n string, err error) {
	buf := make([]uint16, syscall.MAX_LONG_PATH)
	size := uint32(syscall.MAX_LONG_PATH)
	if err := procQueryFullProcessImageNameW.Find(); err == nil { // Vista+
		ret, _, err := procQueryFullProcessImageNameW.Call(
			uintptr(module),
			uintptr(0),
			uintptr(unsafe.Pointer(&buf[0])),
			uintptr(unsafe.Pointer(&size)))
		if ret == 0 {
			return "", err
		}
		return windows.UTF16ToString(buf[:]), nil
	}
	// XP fallback
	ret, _, err := procGetProcessImageFileNameW.Call(uintptr(module), uintptr(unsafe.Pointer(&buf[0])), uintptr(size))
	if ret == 0 {
		return "", err
	}
	return windows.UTF16ToString(buf[:]), nil
}

func getProcessMemoryInfo(h windows.Handle) (uint64, error) {
	var mem process.PROCESS_MEMORY_COUNTERS
	//var ret process.MemoryInfoStat
	r1, _, e1 := syscall.Syscall(procGetProcessMemoryInfo.Addr(), 3, uintptr(h), uintptr(unsafe.Pointer(&mem)), uintptr(unsafe.Sizeof(mem)))
	if r1 == 0 {
		if e1 != 0 {
			return 0, error(e1)
		} else {
			return 0, syscall.EINVAL
		}
	}

	//ret.RSS = uint64(mem.WorkingSetSize)
	//ret.VMS = uint64(mem.PagefileUsage)

	return uint64(mem.PagefileUsage), nil
}

func (p *ProcessInfo) getCpuPercent(c windows.Handle) {
	var CPU windows.Rusage
	if err := windows.GetProcessTimes(c, &CPU.CreationTime, &CPU.ExitTime, &CPU.KernelTime, &CPU.UserTime); err != nil {
		return
	}

	KernelTime := uint64(CPU.KernelTime.LowDateTime) + uint64(CPU.KernelTime.HighDateTime)<<32
	UserTime := uint64(CPU.UserTime.LowDateTime) + uint64(CPU.UserTime.HighDateTime)<<32
	ProcessTime := (KernelTime + UserTime) / 10000

	if p.lastTime == 0 {
		p.lastTime = ProcessTime
		return
	}

	p.Cpu = int((ProcessTime - p.lastTime) * 100 / uint64(p.dwElepsedTime) / uint64(p.ProcessorCoreNum))

	if p.Cpu > 100 || p.Cpu < 0 {
		p.Cpu = 0
	}
	p.lastTime = ProcessTime
}

func (p *ProcessInfo) getUserName(c windows.Handle) {
	var token syscall.Token
	err := syscall.OpenProcessToken(syscall.Handle(c), syscall.TOKEN_QUERY, &token)
	if err != nil {
		return
	}
	defer token.Close()
	tokenUser, err := token.GetTokenUser()
	if err != nil {
		return
	}

	if user, _, _, err := tokenUser.User.Sid.LookupAccount(""); err != nil {
		return
	} else {
		p.Username = user
	}

}

func (p *ProcessInfo) ProcessIdToSessionId() {
	var SessionId uint32
	r1, _, e1 := syscall.Syscall(procProcessIdToSessionId.Addr(), 2, uintptr(p.Pid), uintptr(unsafe.Pointer(&SessionId)), 0)
	if r1 == 0 {
		if e1 != 0 {
			return
		} else {
			return
		}
	}

	p.SessionId = SessionId
}

func (p *ProcessInfo) getProcessInfo() {
	c, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, uint32(p.Pid))
	if err != nil {
		return
	}
	defer windows.CloseHandle(c)

	p.getCpuPercent(c)
	p.getUserName(c)
	p.ProcessIdToSessionId()

	if p.Memory, err = getProcessMemoryInfo(c); err != nil {
		p.Memory = 0
	}

	if p.HandleCount, err = getProcessHandleCount(c); err != nil {
		p.HandleCount = 0
	}

	if p.Path, err = getModuleFileName(c); err != nil {
		p.Path = ""
	}

}

func NewProcessMgr() *ProcessMgr {
	ctx, cancel := context.WithCancel(context.Background())
	mgr := &ProcessMgr{
		ProcessEntrys:   make([]windows.ProcessEntry32, 0),
		processes:       make(map[int32]*ProcessInfo),
		IsEnumProcesses: true,
		chanProcess:     make(chan *ProcessInfo, 1000),
		ctx:             ctx,
		exit:            cancel}
	go mgr.run()
	go mgr.runProcess()
	return mgr
}

type ProcessMgr struct {
	m               sync.RWMutex
	IsEnumProcesses bool
	ProcessEntrys   []windows.ProcessEntry32
	processes       map[int32]*ProcessInfo
	chanProcess     chan *ProcessInfo
	ctx             context.Context
	exit            context.CancelFunc
}

func (mgr *ProcessMgr) runProcess() error {
	for {
		select {
		case <-mgr.ctx.Done():
			return mgr.ctx.Err()
		case proc := <-mgr.chanProcess:
			proc.getProcessInfo()
		}
	}
}

func (mgr *ProcessMgr) run() error {
	d := time.Second
	for {
		select {
		case <-mgr.ctx.Done():
			return mgr.ctx.Err()
		case <-time.After(d):
		}
		if mgr.IsEnumProcesses {
			ProcessEntrys, err := mgr.GetProcessEntrys()
			if err != nil {
				continue
			}

			log.Infof("process cnt=%d", len(ProcessEntrys))

			old := make(map[int32]int32)
			mgr.m.RLock()
			for k, _ := range mgr.processes {
				old[k] = k
			}

			for _, v := range ProcessEntrys {
				_, ok := mgr.processes[int32(v.ProcessID)]
				if ok {
					delete(old, int32(v.ProcessID))
				} else {
					if ProcessorCoreNum == 0 {
						ProcessorCoreNum = 4
					}
					processInfo := &ProcessInfo{
						Pid:              int32(v.ProcessID),
						Name:             windows.UTF16ToString(v.ExeFile[:]),
						ThreadCount:      v.Threads,
						dwElepsedTime:    3000,
						ProcessorCoreNum: ProcessorCoreNum,
						ParentProcessID:  v.ParentProcessID}
					go processInfo.getProcessInfo()
					mgr.processes[int32(v.ProcessID)] = processInfo
				}
			}

			mgr.ProcessEntrys = ProcessEntrys
			mgr.IsEnumProcesses = false
			mgr.m.RUnlock()
		}
	}
}

func (mgr *ProcessMgr) Exit() {
	mgr.exit()
}

func (mgr *ProcessMgr) GetProcessEntrys() ([]windows.ProcessEntry32, error) {
	var procEntrys []windows.ProcessEntry32
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(snapshot)
	var procEntry windows.ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = windows.Process32First(snapshot, &procEntry); err != nil {
		return nil, err
	}

	procEntrys = append(procEntrys, procEntry)
	for {
		err = windows.Process32Next(snapshot, &procEntry)
		if err != nil {
			return procEntrys, nil
		}

		procEntrys = append(procEntrys, procEntry)
	}

	return procEntrys, nil
}

func (mgr *ProcessMgr) Lookup(page, size int) ([]*ProcessInfo, int) {
	infos := make([]*ProcessInfo, 0)
	mgr.m.RLock()
	defer mgr.m.RUnlock()
	mgr.IsEnumProcesses = true
	offset := (page - 1) * size
	if offset < len(mgr.ProcessEntrys) {
		if offset+size >= len(mgr.ProcessEntrys) {
			for _, v := range mgr.ProcessEntrys[offset:] {
				if p, ok := mgr.processes[int32(v.ProcessID)]; ok {
					select {
					case mgr.chanProcess <- p:
					default:
					}
					infos = append(infos, p)
				}
			}
		} else {
			for _, v := range mgr.ProcessEntrys[offset : offset+size] {
				if p, ok := mgr.processes[int32(v.ProcessID)]; ok {

					select {
					case mgr.chanProcess <- p:
					default:
					}
					infos = append(infos, p)
				}
			}
		}
	}

	return infos, len(mgr.ProcessEntrys)
}

func (mgr *ProcessMgr) KillProcess(pid int32) string {
	p, _ := process.NewProcess(pid)
	name, _ := p.Name()
	if err := p.Kill(); err != nil {
		return err.Error()
	} else {
		return fmt.Sprintf("进程：%s ,结束成功", name)
	}
}
