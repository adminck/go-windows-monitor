package systeminfo

import (
	"context"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"sync"
	"time"
)

func NewPercentMgr() *PercentMgr {
	ctx, cancel := context.WithCancel(context.Background())
	mgr := &PercentMgr{CpuMemory: CpuMemoryPercent{}, IsEnumPercent: true, ctx: ctx, exit: cancel}
	go mgr.run()
	return mgr
}

type PercentMgr struct {
	m             sync.RWMutex
	IsEnumPercent bool
	CpuMemory     CpuMemoryPercent
	ctx           context.Context
	exit          context.CancelFunc
}

type CpuMemoryPercent struct {
	Cpu    float64
	Memory *mem.VirtualMemoryStat
}

func (mgr *PercentMgr) run() error {
	d := time.Second
	for {
		select {
		case <-mgr.ctx.Done():
			return mgr.ctx.Err()
		case <-time.After(d):
		}
		if mgr.IsEnumPercent {
			var cpuMemory CpuMemoryPercent
			if percent, err := cpu.Percent(d, false); err == nil {
				cpuMemory.Cpu = percent[0]
			}

			if memInfo, err := mem.VirtualMemory(); err == nil {
				cpuMemory.Memory = memInfo
			}
			mgr.CpuMemory = cpuMemory
			mgr.IsEnumPercent = false
		}

	}
}

func (mgr *PercentMgr) Lookup() CpuMemoryPercent {
	mgr.m.RLock()
	defer mgr.m.RUnlock()
	mgr.IsEnumPercent = true
	return mgr.CpuMemory
}
