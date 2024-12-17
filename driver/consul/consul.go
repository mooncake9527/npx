package consul

import (
	"github.com/pkg/errors"
	"log/slog"
	"sync"
	"time"

	"github.com/mooncake9527/npx/config"
	"github.com/mooncake9527/npx/models"
	"github.com/mooncake9527/npx/scheduling"

	"github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	client             *api.Client
	l                  sync.RWMutex
	registered         []*config.RegisterNode
	discovered         map[string][]*models.ServiceNode //已发现的服务
	schedulingHandlers map[string]scheduling.SchedulingHandler
}

func NewClient(cfg *api.Config) (*ConsulClient, error) {
	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	return &ConsulClient{
		client:             client,
		l:                  sync.RWMutex{},
		discovered:         make(map[string][]*models.ServiceNode),
		registered:         make([]*config.RegisterNode, 0),
		schedulingHandlers: make(map[string]scheduling.SchedulingHandler),
	}, nil
}

func (c *ConsulClient) Register(s *config.RegisterNode) error {
	meta := map[string]string{
		"protocol": string(s.Protocol),
	}
	r := &api.AgentServiceRegistration{
		Namespace: s.Namespace,
		ID:        s.Id,
		Name:      s.Name,
		Port:      s.Port,
		Address:   s.Addr,
		Tags:      s.Tags,
		Meta:      meta,
	}
	var check *api.AgentServiceCheck
	if s.HealthCheck != "" {
		check = &api.AgentServiceCheck{
			Timeout:                        s.Timeout.String(),
			Interval:                       s.Interval.String(),
			DeregisterCriticalServiceAfter: (s.Timeout * 3).String(), //超过3倍超时时间，自动注销
		}
		if s.Protocol == "http" {
			check.HTTP = s.HealthCheck
		} else if s.Protocol == "grpc" {
			check.GRPC = s.HealthCheck
		}
		r.Check = check
	}

	err := c.client.Agent().ServiceRegister(r)
	if err != nil {
		slog.Error("register", "err", err)
		return err
	}
	c.registered = append(c.registered, s)
	return nil
}

func (c *ConsulClient) Deregister() {
	for _, r := range c.registered {
		c.client.Agent().ServiceDeregister(r.Id)
	}
}

func (c *ConsulClient) Watch(s *config.DiscoveryNode) error {
	var lastIndex uint64 = 0
	c.schedulingHandlers[s.Name] = scheduling.GetHandler(s.SchedulingAlgorithm)
	go func(s *config.DiscoveryNode) {
		for {
			entries, qmeta, err := c.client.Health().Service(s.Name, s.Tag, false, &api.QueryOptions{
				Namespace: s.Namespace,
				WaitIndex: lastIndex,
			})
			if err != nil {
				slog.Error("watch", "err", err)
				time.Sleep(time.Second * 1)
				continue
			}
			lastIndex = qmeta.LastIndex
			slog.Debug("watch", "entries", entries, "qmeta", qmeta)
			for _, entry := range entries {
				status := entry.Checks.AggregatedStatus()
				slog.Debug("watch", "-------status--", status, "entry.Service.ID--------", entry.Service.ID)
				switch status {
				case api.HealthPassing:
					c.putServiceNode(s, entry)
				//case api.HealthMaint, api.HealthCritical, api.HealthWarning:
				default:
					c.delServiceNode(s, entry)
				}
			}
			time.Sleep(time.Second * time.Duration(s.RetryTime))
		}
	}(s)
	return nil
}

func (c *ConsulClient) GetService(name string, clientIp string) (*models.ServiceNode, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	if rs, ok := c.discovered[name]; ok && len(rs) > 0 {
		if sh, ok := c.schedulingHandlers[name]; ok {
			return sh.GetServiceNode(rs, name), nil
		}
	}
	return nil, errors.New("no service")
}

func (c *ConsulClient) putServiceNode(s *config.DiscoveryNode, entry *api.ServiceEntry) {
	c.l.Lock()
	defer c.l.Unlock()
	if vs, ok := c.discovered[s.Name]; ok {
		found := false
		for _, v := range vs {
			if v.Id == entry.Service.ID {
				slog.Debug("watch", "update---", entry.Service.ID)
				found = true
				v.Addr = entry.Service.Address
				v.Port = entry.Service.Port
				v.Protocol = entry.Service.Meta["protocol"]
				v.Namespace = entry.Service.Namespace
				v.Tags = entry.Service.Tags
				v.SetEnable(true)
				v.ClearFailCnt()
				break
			}
		}
		if !found {
			slog.Debug("watch", "add reset---", entry.Service.ID)
			ds := c.entryToServiceNode(entry, s)
			vs = append(vs, ds)
			c.discovered[s.Name] = vs
		}
	} else {
		slog.Debug("watch", "add 1---", entry.Service.ID)
		ds := c.entryToServiceNode(entry, s)
		c.discovered[s.Name] = []*models.ServiceNode{ds}
	}
}

func (c *ConsulClient) delServiceNode(s *config.DiscoveryNode, entry *api.ServiceEntry) {
	c.l.Lock()
	defer c.l.Unlock()
	slog.Debug("watch", "entry.Service.ID-----del---", entry.Service.ID)
	if vs, ok := c.discovered[s.Name]; ok {
		for i, v := range vs {
			if v.Id == entry.Service.ID {
				v.Close()
				vs = append(vs[:i], vs[i+1:]...)
				break
			}
		}
		c.discovered[s.Name] = vs
	} else {
		slog.Debug("not found")
	}
}

func (c *ConsulClient) entryToServiceNode(entry *api.ServiceEntry, s *config.DiscoveryNode) *models.ServiceNode {
	r := config.RegisterNode{
		Id:        entry.Service.ID,
		Namespace: entry.Service.Namespace,
		Name:      entry.Service.Service,
		Tags:      entry.Service.Tags,
		Addr:      entry.Service.Address,
		Port:      entry.Service.Port,
		Protocol:  entry.Service.Meta["protocol"],
		FailLimit: s.FailLimit,
	}
	n := models.ServiceNode{
		RegisterNode: r,
	}
	n.SetEnable(true)

	return &n
}
