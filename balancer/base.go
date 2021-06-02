package balancer

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/jqiris/kungfu/session"
	"github.com/jqiris/kungfu/utils"
	"net/http"
	"net/url"

	"github.com/apex/log"
	"github.com/jqiris/kungfu/coder"
	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/helper"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
)

type BaseBalancer struct {
	ServerId              string
	Server                *treaty.Server
	Rpcx                  rpcx.RpcBalancer
	ClientServer          *http.Server
	ClientCoder           coder.Coder
	EventHandlerSelf      func(req []byte) []byte //处理自己的事件
	EventHandlerBroadcast func(req []byte) []byte //处理广播事件
}

func (b *BaseBalancer) HandleBalance(w http.ResponseWriter, r *http.Request) {
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	serverType, uid := "", 0
	if err == nil {
		if len(queryForm["server_type"]) > 0 {
			serverType = queryForm["server_type"][0]
		}
		if len(queryForm["uid"]) > 0 {
			uid = utils.StringToInt(queryForm["uid"][0])
		}
	}
	if len(serverType) < 1 {
		res := &treaty.BalanceResult{
			Code: treaty.CodeType_CodeChooseBackendLogin,
		}
		b.WriteResponse(w, res)
		return
	}
	connetor, err := b.Balance(r.RemoteAddr)
	if err != nil {
		res := &treaty.BalanceResult{
			Code: treaty.CodeType_CodeFailed,
		}
		b.WriteResponse(w, res)
		return
	}
	backend := discover.GetServerByType(serverType, r.RemoteAddr)
	var backendPre *treaty.Server
	sess := session.GetSession(int32(uid))
	if sess != nil {
		backendPre = sess.Backend
	}
	res := &treaty.BalanceResult{
		Code:       treaty.CodeType_CodeSuccess,
		Connector:  connetor,
		Backend:    backend,
		BackendPre: backendPre,
	}
	b.WriteResponse(w, res)
}

func (b *BaseBalancer) WriteResponse(w http.ResponseWriter, msg proto.Message) {
	if v, e := b.ClientCoder.Marshal(msg); e == nil {
		if _, e2 := w.Write(v); e2 != nil {
			logger.Error(e2)
		}
	}
}
func (b *BaseBalancer) Init() {
	//find the  server config
	if b.Server = helper.FindServerConfig(conf.GetServersConf(), b.GetServerId()); b.Server == nil {
		logger.Fatal("BaseBalancer can find the server config")
	}
	//init the rpcx
	b.Rpcx = rpcx.NewRpcBalancer(conf.GetRpcxConf())
	//init the coder
	b.ClientCoder = coder.NewJsonCoder()
	//set the server
	b.ClientServer = &http.Server{Addr: fmt.Sprintf("%s:%d", b.Server.ServerIp, b.Server.ClientPort)}
	//handle the blance
	http.HandleFunc("/balance", b.HandleBalance)
	//run the server
	go func() {
		err := b.ClientServer.ListenAndServe()
		if err != nil {
			log.Error(err.Error())
		}
	}()
	logger.Infoln("init the balancer:", b.ServerId)
}

func (b *BaseBalancer) AfterInit() {
	//Subscribe event
	if err := b.Rpcx.Subscribe(b.Server, func(req []byte) []byte {
		logger.Infof("BaseBalancer Subscribe received: %+v", req)
		return b.EventHandlerSelf(req)
	}); err != nil {
		logger.Error(err)
	}
	if err := b.Rpcx.SubscribeBalancer(func(req []byte) []byte {
		logger.Infof("BaseBalancer SubscribeBalancer received: %+v", req)
		return b.EventHandlerBroadcast(req)
	}); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseBalancer) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *BaseBalancer) Shutdown() {
	if b.ClientServer != nil {
		if err := b.ClientServer.Close(); err != nil {
			logger.Error(err)
		}
	}
	logger.Infoln("stop the balancer:", b.ServerId)
}

func (b *BaseBalancer) Balance(remoteAddr string) (*treaty.Server, error) {
	if server := discover.GetServerByType("connector", remoteAddr); server != nil {
		return server, nil
	}

	return nil, errors.New("no suitable connector found")
}

func (b *BaseBalancer) GetServer() *treaty.Server {
	return b.Server
}

func (b *BaseBalancer) RegEventHandlerSelf(handler func(req []byte) []byte) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *BaseBalancer) RegEventHandlerBroadcast(handler func(req []byte) []byte) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}
func (b *BaseBalancer) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *BaseBalancer) GetServerId() string {
	return b.ServerId
}
