package systeminfo

import (
	"context"
	"github.com/StackExchange/wmi"
	"sync"
	"time"
)

func NewServiceMgr() *ServiceMgr {
	ctx, cancel := context.WithCancel(context.Background())
	mgr := &ServiceMgr{services: make([]Win32_Service, 0), ctx: ctx, exit: cancel, IsEnumService: true}
	go mgr.run()
	return mgr
}

type ServiceMgr struct {
	m             sync.RWMutex
	services      []Win32_Service
	IsEnumService bool
	ctx           context.Context
	exit          context.CancelFunc
}

type Win32_Service struct {
	Name        string `json:"name"`
	Caption     string `json:"caption"`
	State       string `json:"state"`
	Description string `json:"description"`
	StartMode   string `json:"startMode"`
}

func (mgr *ServiceMgr) run() error {
	d := time.Second
	for {
		select {
		case <-mgr.ctx.Done():
			return mgr.ctx.Err()
		case <-time.After(d):
		}
		d = time.Second * 5

		mgr.m.Lock()
		if mgr.IsEnumService {
			var s []Win32_Service
			q := wmi.CreateQuery(&s, "")
			wmi.Query(q, &s)
			mgr.IsEnumService = false
			mgr.services = s
		}

		mgr.m.Unlock()
	}
}

func (mgr *ServiceMgr) Lookup() []Win32_Service {
	mgr.m.RLock()
	defer mgr.m.RUnlock()
	mgr.IsEnumService = true
	return mgr.services
}
