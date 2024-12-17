package etcd

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/mooncake9527/npx/config"
	"github.com/mooncake9527/npx/models"
	"github.com/mooncake9527/npx/scheduling"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	client             *clientv3.Client
	l                  sync.RWMutex
	registered         map[string]*config.RegisterNode
	discovered         map[string][]*models.ServiceNode //已发现的服务
	schedulingHandlers map[string]scheduling.SchedulingHandler
}

func NewClient(cfg *clientv3.Config) (*EtcdClient, error) {
	client, err := clientv3.New(*cfg)
	if err != nil {
		return nil, err
	}
	return &EtcdClient{
		client:             client,
		l:                  sync.RWMutex{},
		discovered:         make(map[string][]*models.ServiceNode),
		registered:         make(map[string]*config.RegisterNode),
		schedulingHandlers: make(map[string]scheduling.SchedulingHandler),
	}, nil
}

func (c *EtcdClient) Register(s *config.RegisterNode) error {
	var err error
	go func() {
		kv := clientv3.NewKV(c.client)
		lease := clientv3.NewLease(c.client)
		var curLeaseId clientv3.LeaseID = 0
		for {
			if curLeaseId == 0 {
				leaseResp, err := lease.Grant(context.TODO(), int64(s.Timeout.Seconds()))
				if err != nil {
					slog.Error("Grant err", err)
					return
				}
				key := fmt.Sprintf("%s%d", s.Name, leaseResp.ID)
				b, err := json.Marshal(s)
				if err != nil {
					slog.Error("json marshal err", err)
				}
				if _, err := kv.Put(context.TODO(), key, string(b), clientv3.WithLease(clientv3.LeaseID(leaseResp.ID))); err != nil {
					slog.Error("put err", err)
					return
				}
				c.registered[key] = s
				curLeaseId = clientv3.LeaseID(leaseResp.ID)
			} else {
				slog.Debug("key:", curLeaseId)
				// 续约租约，如果租约已经过期将curLeaseId复位到0重新走创建租约的逻辑
				if _, err := lease.KeepAliveOnce(context.TODO(), curLeaseId); err == rpctypes.ErrLeaseNotFound {
					slog.Error("keepalive err", err)
					curLeaseId = 0
					continue
				}
			}
			time.Sleep(s.Interval)
		}
	}()
	return err
}

func (c *EtcdClient) Deregister() {
	kv := clientv3.NewKV(c.client)
	for k, _ := range c.registered {
		kv.Delete(context.TODO(), k)
	}
}

func (c *EtcdClient) Watch(s *config.DiscoveryNode) error {
	go func() {
		func() {
			kv := clientv3.NewKV(c.client)
			rangeResp, err := kv.Get(context.TODO(), s.Name, clientv3.WithPrefix())
			if err != nil {
				slog.Error("GetKey err:", err)
			}
			for _, kv := range rangeResp.Kvs {
				c.putServiceNode(kv.Value, s)
			}
		}()

		watcher := clientv3.NewWatcher(c.client)
		// Watch 服务目录下的更新
		watchChan := watcher.Watch(context.TODO(), s.Name, clientv3.WithPrefix())
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				slog.Info("Events ", s.Name, string(event.Kv.Value))
				switch event.Type {
				case mvccpb.PUT: //PUT事件，目录下有了新key
					c.putServiceNode(event.Kv.Value, s)
				case mvccpb.DELETE: //DELETE事件，目录中有key被删掉(Lease过期，key 也会被删掉)
					c.delServiceNode(string(event.Kv.Key), s)
				}
			}
		}
	}()
	return nil
}

func (c *EtcdClient) putServiceNode(data []byte, s *config.DiscoveryNode) {
	c.l.Lock()
	defer c.l.Unlock()
	slog.Debug("-----put ", s.Name, string(data))
	var rs models.ServiceNode
	err := json.Unmarshal(data, &rs)
	if err != nil {
		slog.Error("unmarshal err", err)
	}
	if vs, ok := c.discovered[s.Name]; ok {
		found := false
		for _, v := range vs {
			if v.Id == rs.Id {
				slog.Debug("-----update ", s.Name, rs)
				v.Addr = rs.Addr
				v.Port = rs.Port
				v.Tags = rs.Tags
				v.Weight = rs.Weight
				v.Namespace = rs.Namespace
				v.Protocol = rs.Protocol
				v.SetEnable(true)
				v.ClearFailCnt()
				found = true
				break
			}
		}
		if !found {
			slog.Debug("-----add other ", s.Name, rs)
			rs.SetEnable(true)
			rs.ClearFailCnt()
			vs = append(vs, &rs)
			c.discovered[s.Name] = vs
		}
	} else {
		slog.Debug("-----add first", s.Name, rs)
		rs.SetEnable(true)
		rs.ClearFailCnt()
		c.discovered[s.Name] = []*models.ServiceNode{&rs}
	}
	c.schedulingHandlers[s.Name] = scheduling.GetHandler(s.SchedulingAlgorithm)
}

func (c *EtcdClient) delServiceNode(curId string, s *config.DiscoveryNode) {
	c.l.RLock()
	defer c.l.RUnlock()
	slog.Debug("-----del ", s.Name, curId)
	if vs, ok := c.discovered[s.Name]; ok {
		for i, v := range vs {
			if v.Id == curId {
				slog.Debug("-----del ", s.Name, curId)
				v.Close()
				vs = append(vs[:i], vs[i+1:]...)
				break
			}
		}
		c.discovered[s.Name] = vs
	} else {
		slog.Info("not found")
	}
}

func (c *EtcdClient) GetService(name string, clientIp string) (*models.ServiceNode, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	slog.Debug("-----get ", name)
	if rs, ok := c.discovered[name]; ok && len(rs) > 0 {
		if sh, ok := c.schedulingHandlers[name]; ok {
			return sh.GetServiceNode(rs, name), nil
		}
	}
	return nil, errors.New("no service")
}
