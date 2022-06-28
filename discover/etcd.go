package discover

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"stathat.com/c/consistent"

	"github.com/jqiris/kungfu/v2/logger"
	"github.com/jqiris/kungfu/v2/treaty"
	"github.com/jqiris/kungfu/v2/utils"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// EtcdDiscoverer etcd discoverer
type EtcdDiscoverer struct {
	Config           clientv3.Config
	Client           *clientv3.Client
	ServerList       map[string]*treaty.Server  //serverId=>server
	ServerTypeMap    map[string]*ServerTypeItem //serverType=>serverTypeItem
	ServerLock       *sync.RWMutex
	EventHandlerList []EventHandler
	Prefix           string
}
type EtcdOption func(e *EtcdDiscoverer)

func WithEtcdEndpoints(endpoints []string) EtcdOption {
	return func(e *EtcdDiscoverer) {
		e.Config.Endpoints = endpoints
	}
}

func WithEtcdDialTimeOut(d time.Duration) EtcdOption {
	return func(e *EtcdDiscoverer) {
		e.Config.DialTimeout = d
	}
}
func WithEtcdPrefix(prefix string) EtcdOption {
	return func(e *EtcdDiscoverer) {
		prefix = "/" + prefix + "/"
		e.Prefix = prefix
	}
}

// NewEtcdDiscoverer init EtcdDiscoverer
func NewEtcdDiscoverer(opts ...EtcdOption) *EtcdDiscoverer {
	e := &EtcdDiscoverer{
		ServerList:       make(map[string]*treaty.Server),
		ServerTypeMap:    make(map[string]*ServerTypeItem),
		ServerLock:       new(sync.RWMutex),
		EventHandlerList: make([]EventHandler, 0),
		Prefix:           "/server/",
	}
	for _, opt := range opts {
		opt(e)
	}
	cli, err := clientv3.New(e.Config)
	if err != nil {
		logger.Fatal(err)
		return nil
	}
	e.Client = cli
	e.Init()
	return e
}

func (e *EtcdDiscoverer) RegEventHandlers(handlers ...EventHandler) {
	e.EventHandlerList = append(e.EventHandlerList, handlers...)
}

func (e *EtcdDiscoverer) EventHandlerExec(ev *clientv3.Event, server *treaty.Server) {
	for _, handler := range e.EventHandlerList {
		handler(ev, server)
	}
}

// Init init
func (e *EtcdDiscoverer) Init() {
	//监听服务器变化
	go utils.SafeRun(func() {
		e.Watcher()
	})
	//统计所有服务器
	list := e.FindServerList()
	if len(list) > 0 {
		for k, v := range list {
			e.ServerLock.Lock()
			if _, ok := e.ServerTypeMap[k]; !ok {
				item := NewServerTypeItem()
				for _, vv := range v {
					e.ServerList[vv.ServerId] = vv
					item.List[vv.ServerId] = vv
					item.hash.Add(vv.ServerId)
				}
				e.ServerTypeMap[k] = item
			}
			e.ServerLock.Unlock()
		}
	}
	e.DumpServers()
}

func (e *EtcdDiscoverer) IsCurEvent(key string) bool {
	prefix := ""
	keyArr := strings.Split(key, "/")
	if len(keyArr) > 1 {
		prefix = "/" + keyArr[1] + "/"
	}
	return e.Prefix == prefix
}

func (e *EtcdDiscoverer) Watcher() {
	for {
		rch := e.Client.Watch(context.Background(), e.Prefix, clientv3.WithPrefix())
		var err error
		for wResp := range rch {
			err = wResp.Err()
			if err != nil {
				logger.Errorf("etcd watch err:%v", err)
			}

			for _, ev := range wResp.Events {
				//logger.Infof("%s %q %q", ev.Type, ev.Kv.Key, ev.Kv.Value)
				if !e.IsCurEvent(string(ev.Kv.Key)) {
					continue
				}
				var silent int32 = 0
				switch ev.Type {
				case clientv3.EventTypePut:
					if server, err := treaty.RegUnSerialize(ev.Kv.Value); err == nil {
						e.ServerLock.Lock()
						e.ServerList[server.ServerId] = server
						if item, ok := e.ServerTypeMap[server.ServerType]; ok {
							if _, okv := item.List[server.ServerId]; !okv {
								item.hash.Add(server.ServerId)
							}
							item.List[server.ServerId] = server
						} else {
							item = NewServerTypeItem()
							item.hash.Add(server.ServerId)
							item.List[server.ServerId] = server
							e.ServerTypeMap[server.ServerType] = item
						}
						e.ServerLock.Unlock()
						e.EventHandlerExec(ev, server)
						silent = atomic.LoadInt32(&server.Silent)
					}
				case clientv3.EventTypeDelete:
					ks := strings.Split(string(ev.Kv.Key), "/")
					if len(ks) > 2 {
						sType, sid := ks[len(ks)-2], ks[len(ks)-1]
						e.ServerLock.Lock()
						delete(e.ServerList, sid)
						if item, ok := e.ServerTypeMap[sType]; ok {
							if server, okv := item.List[sid]; okv {
								item.hash.Remove(sid)
								delete(item.List, sid)
								e.EventHandlerExec(ev, server)
							}
							if len(item.List) == 0 {
								delete(e.ServerTypeMap, sType)
							}

						}
						e.ServerLock.Unlock()
					}
				}
				if silent == 0 {
					e.DumpServers()
				}
			}
		}
	}
}

func (e *EtcdDiscoverer) DumpServers() {
	logger.Info("#####################################DUMP SERVERS BEGIN#################################")
	for typ, list := range e.ServerTypeMap {
		logger.Info("------------------------------------------------------------------------------------")
		for _, server := range list.List {
			logger.Infof("type:%v, server:%+v", typ, server)
		}
	}
	logger.Info("#####################################DUMP SERVERS END###################################")
}

// Register register
func (e *EtcdDiscoverer) Register(server *treaty.Server) error {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	key, val := e.Prefix+treaty.RegSeverItem(server), treaty.RegSerialize(server)
	resp, err := kv.Put(ctx, key, val)
	if err != nil {
		return err
	}
	if atomic.LoadInt32(&server.Silent) == 0 {
		logger.Infof("discover Register server,k=>v,%s=>%s,resp:%v", key, val, resp.Header)
	}
	return nil
}

func (e *EtcdDiscoverer) UnRegister(server *treaty.Server) error {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	if resp, err := kv.Delete(ctx, e.Prefix+treaty.RegSeverItem(server), clientv3.WithPrevKV()); err != nil {
		return err
	} else if atomic.LoadInt32(&server.Silent) == 0 {
		logger.Infof("EtcdDiscoverer unregister resp:%+v", resp)
	}
	return nil
}

func (e *EtcdDiscoverer) IncreLoad(serverId string, load int64) error {
	server := e.GetServerById(serverId, false)
	if server == nil {
		return fmt.Errorf("IncreLoad can't find server %s", serverId)
	}
	atomic.AddInt64(&server.Load, load)
	atomic.StoreInt32(&server.Silent, 1)
	return e.Register(server)
}

func (e *EtcdDiscoverer) DecreLoad(serverId string, load int64) error {
	server := e.GetServerById(serverId, false)
	if server == nil {
		return fmt.Errorf("IncreLoad can't find server %s", serverId)
	}
	atomic.AddInt64(&server.Load, -load)
	atomic.StoreInt32(&server.Silent, 1)
	return e.Register(server)
}

func (e *EtcdDiscoverer) FindServer(serverType string) []*treaty.Server {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	if resp, err := kv.Get(ctx, e.Prefix+serverType+"/", clientv3.WithPrefix()); err != nil {
		logger.Errorf("EtcdDiscoverer FindServer err:%v", err)
		return nil
	} else {
		if resp.Count > 0 {
			res := make([]*treaty.Server, 0)
			for _, v := range resp.Kvs {
				if server, err := treaty.RegUnSerialize(v.Value); err == nil {
					res = append(res, server)
				} else {
					logger.Errorf("EtcdDiscoverer FindServer err:%+v", err)
				}
			}
			return res
		}
	}
	return nil
}

func (e *EtcdDiscoverer) FindServerList() map[string][]*treaty.Server {
	kv := clientv3.NewKV(e.Client)
	ctx, cancel := context.WithTimeout(context.TODO(), e.Config.DialTimeout)
	defer cancel()
	if resp, err := kv.Get(ctx, e.Prefix, clientv3.WithPrefix()); err != nil {
		logger.Errorf("EtcdDiscoverer FindServerList err:%v", err)
		return nil
	} else {
		if resp.Count > 0 {
			res := make(map[string][]*treaty.Server)
			for _, v := range resp.Kvs {
				if server, err := treaty.RegUnSerialize(v.Value); err == nil {
					res[server.ServerType] = append(res[server.ServerType], server)
				} else {
					logger.Errorf("EtcdDiscoverer FindServerList err:%+v", err)
				}
			}
			return res
		}
	}
	return nil
}

func (e *EtcdDiscoverer) GetServerList(args ...bool) map[string]*treaty.Server {
	e.ServerLock.RLock()
	defer e.ServerLock.RUnlock()
	filterMaintain := true
	if len(args) > 0 {
		filterMaintain = args[0]
	}
	if filterMaintain {
		list := make(map[string]*treaty.Server)
		for k, v := range e.ServerList {
			if filterMaintain && v.Maintained {
				continue
			}
			list[k] = v
		}
		return list
	}
	return e.ServerList
}

func (e *EtcdDiscoverer) GetServerById(serverId string, args ...bool) *treaty.Server {
	e.ServerLock.RLock()
	defer e.ServerLock.RUnlock()
	filterMaintain := true
	if len(args) > 0 {
		filterMaintain = args[0]
	}
	if v, ok := e.ServerList[serverId]; ok {
		if filterMaintain && v.Maintained {
			return nil
		}
		return v
	}
	return nil
}

func (e *EtcdDiscoverer) GetServerByType(serverType, serverArg string, args ...bool) *treaty.Server {
	e.ServerLock.RLock()
	defer e.ServerLock.RUnlock()
	filterMaintain := true
	if len(args) > 0 {
		filterMaintain = args[0]
	}
	if item, ok := e.ServerTypeMap[serverType]; ok {
		if filterMaintain {
			var srvs []string
			hasMaintained := false
			for k, v := range item.List {
				if !v.Maintained {
					srvs = append(srvs, k)
				} else {
					hasMaintained = true
				}
			}
			if hasMaintained {
				cons := consistent.New()
				cons.Set(srvs)
				if sid, err := cons.Get(serverArg); err == nil {
					return item.List[sid]
				} else {
					return nil
				}
			}
		}
		if sid, err := item.hash.Get(serverArg); err == nil {
			return item.List[sid]
		}
	}
	return nil
}
func (e *EtcdDiscoverer) GetServerByTypeLoad(serverType string, args ...bool) *treaty.Server {
	e.ServerLock.RLock()
	defer e.ServerLock.RUnlock()
	filterMaintain := true
	if len(args) > 0 {
		filterMaintain = args[0]
	}
	if item, ok := e.ServerTypeMap[serverType]; ok {
		var server *treaty.Server
		for _, v := range item.List {
			if filterMaintain && v.Maintained {
				continue
			}
			if server == nil {
				server = v
			} else if v.Load >= 0 && server.Load > v.Load {
				server = v
			}
		}
		return server
	}
	return nil
}
func (e *EtcdDiscoverer) GetServerTypeList(serverType string, args ...bool) map[string]*treaty.Server {
	e.ServerLock.RLock()
	defer e.ServerLock.RUnlock()
	filterMaintain := true
	if len(args) > 0 {
		filterMaintain = args[0]
	}
	if item, ok := e.ServerTypeMap[serverType]; ok {
		if filterMaintain {
			list := make(map[string]*treaty.Server)
			for k, v := range item.List {
				if !v.Maintained {
					list[k] = v
				}
			}
			return list
		}
		return item.List
	}
	return nil
}
