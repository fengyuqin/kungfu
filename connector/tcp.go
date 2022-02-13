package connector

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpc"
	"github.com/jqiris/kungfu/tcpface"
	"github.com/jqiris/kungfu/tcpserver"
	"github.com/jqiris/kungfu/treaty"
)

type TcpConnector struct {
	ServerId              string
	Server                *treaty.Server
	Rpc                   rpc.ServerRpc
	EventJsonSelf         rpc.CallbackFunc        //处理自己的json事件
	EventHandlerSelf      rpc.CallbackFunc        //处理自己的事件
	EventHandlerBroadcast rpc.CallbackFunc        //处理广播事件
	ClientServer          tcpface.IServer         //zinx server
	RouteHandler          func(s tcpface.IServer) //注册路由
}

func (b *TcpConnector) Init() {
	if b.Server == nil {
		panic("服务配置信息不能为空")
		return
	}
	if b.RouteHandler == nil {
		panic("路由配置信息不能为空")
		return
	}
	//赋值id
	b.ServerId = b.Server.ServerId
	//init the rpc
	b.Rpc = rpc.NewRpcServer(config.GetRpcConf(), b.Server)
	//run the front server
	b.ClientServer = tcpserver.NewServer(b.Server)
	b.RouteHandler(b.ClientServer)
	go b.ClientServer.Serve()

	logger.Info("init the connector:", b.ServerId)
}

func (b *TcpConnector) AfterInit() {
	if b.Server == nil {
		panic("服务配置信息不能为空")
		return
	}
	if b.EventJsonSelf == nil {
		panic("EventJsonSelf不能为空")
		return
	}
	if b.EventHandlerSelf == nil {
		panic("EventHandlerSelf不能为空")
		return
	}
	if b.EventHandlerBroadcast == nil {
		panic("EventHandlerBroadcast不能为空")
		return
	}
	builder := rpc.NewSubscriberRpc(b.Server).SetCodeType(rpc.CodeTypeProto).SetCallback(func(req *rpc.MsgRpc) []byte {
		return b.EventHandlerSelf(req)
	})
	//Subscribe event
	if err := b.Rpc.Subscribe(builder.Build()); err != nil {
		logger.Error(err)
	}
	builder = builder.SetSuffix("json").SetCodeType(rpc.CodeTypeJson).SetCallback(func(req *rpc.MsgRpc) []byte {
		return b.EventJsonSelf(req)
	})
	//Subscribe event
	if err := b.Rpc.Subscribe(builder.Build()); err != nil {
		logger.Error(err)
	}
	builder = builder.SetSuffix(rpc.DefaultSuffix).SetCodeType(rpc.CodeTypeProto).SetCallback(func(req *rpc.MsgRpc) []byte {
		return b.EventHandlerBroadcast(req)
	})
	if err := b.Rpc.SubscribeConnector(builder.Build()); err != nil {
		logger.Error(err)
	}
	//register the service
	if err := discover.Register(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *TcpConnector) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *TcpConnector) Shutdown() {
	//stop the server
	if b.ClientServer != nil {
		b.ClientServer.Stop()
	}
	logger.Info("stop the connector:", b.ServerId)
}

func (b *TcpConnector) GetServer() *treaty.Server {
	return b.Server
}

func (b *TcpConnector) SetRouteHandler(handler func(s tcpface.IServer)) {
	b.RouteHandler = handler
}
