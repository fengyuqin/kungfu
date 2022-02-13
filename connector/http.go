package connector

import (
	"github.com/jqiris/kungfu/config"
	"github.com/jqiris/kungfu/discover"
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/rpc"
	"github.com/jqiris/kungfu/treaty"
)

type HttpConnector struct {
	ServerId              string
	Server                *treaty.Server
	Rpc                   rpc.ServerRpc
	EventJsonSelf         rpc.CallbackFunc //处理自己的json事件
	EventHandlerSelf      rpc.CallbackFunc //处理自己的事件
	EventHandlerBroadcast rpc.CallbackFunc //处理广播事件
	ConnectorConf         config.ConnectorConf
}

func (b *HttpConnector) Init() {
	if b.Server == nil {
		panic("服务配置信息不能为空")
		return
	}
	//赋值id
	b.ServerId = b.Server.ServerId
	b.ConnectorConf = config.GetConnectorConf()
	//init the rpc
	b.Rpc = rpc.NewRpcServer(config.GetRpcConf(), b.Server)
}

func (b *HttpConnector) AfterInit() {
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

func (b *HttpConnector) BeforeShutdown() {
	//unregister the service
	if err := discover.UnRegister(b.Server); err != nil {
		logger.Error(err)
	}
}

func (b *HttpConnector) Shutdown() {
	logger.Info("stop the connector:", b.ServerId)
}

func (b *HttpConnector) GetServer() *treaty.Server {
	return b.Server
}
func (b *HttpConnector) RegEventJsonSelf(handler rpc.CallbackFunc) { //注册自己事件处理器
	b.EventJsonSelf = handler
}
func (b *HttpConnector) RegEventHandlerSelf(handler rpc.CallbackFunc) { //注册自己事件处理器
	b.EventHandlerSelf = handler
}

func (b *HttpConnector) RegEventHandlerBroadcast(handler rpc.CallbackFunc) { //注册广播事件处理器
	b.EventHandlerBroadcast = handler
}
func (b *HttpConnector) SetServerId(serverId string) {
	b.ServerId = serverId
}

func (b *HttpConnector) GetServerId() string {
	return b.ServerId
}
