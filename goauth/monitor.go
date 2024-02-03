package goauth

import (
	"fmt"
	"time"

	"github.com/curtisnewbie/gocommon/common"
	"github.com/curtisnewbie/miso/miso"
)

const (
	DefaultMonitorPath = "/auth/resource"
)

var (
	monitorServiceTickser = []*miso.TickRunner{}
	monitorPool           = miso.NewAsyncPool(500, 20)
)

type MonitorConf struct {
	Monitor []MonitoredService
}

type MonitoredService struct {
	Service string
	Path    string
	All     bool
}

func LoadMonitoredServices() []MonitoredService {
	var c MonitorConf
	miso.UnmarshalFromProp(&c)
	for i, m := range c.Monitor {
		if m.Path == "" {
			m.Path = DefaultMonitorPath
			c.Monitor[i] = m
		}
	}
	return c.Monitor
}

type QueryResourcePathRes struct {
	Resources []CreateResReq
	Paths     []CreatePathReq
}

func QueryResourcePath(rail miso.Rail, server miso.Server, service string, path string) (QueryResourcePathRes, error) {
	var resp miso.GnResp[QueryResourcePathRes]
	err := miso.NewTClient(rail, server.BuildUrl(path)).
		Require2xx().
		Get().
		Json(&resp)
	if err != nil {
		return QueryResourcePathRes{}, fmt.Errorf("failed to query resource path from monitored service, server: %+v, service: %v, %w",
			server, service, err)
	}
	return resp.Res()
}

func CreateMonitoredServiceWatches(rail miso.Rail) error {
	services := LoadMonitoredServices()
	for i := range services {
		s := services[i]
		if err := CreateMonitoredServiceWatch(rail, s); err != nil {
			return err
		}
	}
	return nil
}

func CreateMonitoredServiceWatch(rail miso.Rail, m MonitoredService) error {
	triggered := func() {
		servers := miso.ListServers(m.Service)
		miso.Debugf("Servers for %v: %+v", m.Service, servers)
		if len(servers) < 1 {
			return
		}

		doMonitor := func(server miso.Server) {
			monitorPool.Go(func() {
				rail := miso.EmptyRail()
				res, err := QueryResourcePath(rail, server, m.Service, m.Path)
				if err != nil {
					rail.Errorf("monitor service %v failed, %v", m.Service, err)
				} else {
					rail.Debugf("service %v (%v:%v), returned resouces/paths: %+v", m.Service, server.Address, server.Port, res)
					user := common.GetUser(rail) // just to satisfy the method, it's always a zero value
					for _, r := range res.Resources {
						if err := CreateResourceIfNotExist(rail, r, user); err != nil {
							rail.Errorf("failed to create resource, req: %+v, %v", r, err)
						}
					}
					for _, r := range res.Paths {
						if err := CreatePathIfNotExist(rail, r, user); err != nil {
							rail.Errorf("failed to create path, req: %+v, %v", r, err)
						}
					}
				}
			})
		}

		if m.All {
			for i := range servers {
				server := servers[i]
				doMonitor(server)
			}
		} else {
			server := servers[miso.RandomServerSelector(servers)]
			doMonitor(server)
		}
	}

	if err := miso.SubscribeServerChanges(rail, m.Service, triggered); err != nil {
		return fmt.Errorf("failed to subscribe server chagnes, service: %v, %v", m.Service, err)
	}

	tr := miso.NewTickRuner(time.Minute*1, triggered)
	monitorServiceTickser = append(monitorServiceTickser, tr)
	tr.Start()
	return nil
}
