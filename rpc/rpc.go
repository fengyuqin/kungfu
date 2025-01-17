/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package rpc

import (
	"time"

	"github.com/fengyuqin/kungfu/v2/config"
	"github.com/fengyuqin/kungfu/v2/discover"
	"github.com/fengyuqin/kungfu/v2/logger"
	"github.com/fengyuqin/kungfu/v2/treaty"
	"github.com/nats-io/nats.go"
)

type HttpHandler interface {
	Run(addr ...string) error
}

type CallbackFunc func(req *MsgRpc) []byte

// ServerRpc rpc interface
type ServerRpc interface {
	RegEncoder(typ string, encoder EncoderRpc)                                        //register encoder
	Subscribe(s RssBuilder) error                                                     //self Subscribe
	SubscribeBroadcast(s RssBuilder) error                                            //broadcast subscribe
	QueueSubscribe(s RssBuilder) error                                                //queue self Subscribe
	SendMsg(s ReqBuilder) error                                                       //send msg direct
	Publish(s ReqBuilder) error                                                       //publish
	QueuePublish(s ReqBuilder) error                                                  //queue publish
	PublishBroadcast(s ReqBuilder) error                                              //broadcast publish
	Request(s ReqBuilder) error                                                       //request
	QueueRequest(s ReqBuilder) error                                                  //queue request
	Response(codeType string, v any) []byte                                           //response the msg
	DecodeMsg(codeType string, data []byte, v any) error                              //decode msg
	GetCoder(codeType string) EncoderRpc                                              //get encoder
	GetServer() *treaty.Server                                                        //get current server
	Find(serverType string, arg any, options ...discover.FilterOption) *treaty.Server //find server
	RemoveFindCache(arg any)                                                          //clear find cache
	Close() error                                                                     //close option
}

// NewRpcServer create rpc server
func NewRpcServer(cfg config.RpcConf, server *treaty.Server) ServerRpc {
	timeout := time.Duration(cfg.DialTimeout) * time.Second
	var r ServerRpc
	switch cfg.UseType {
	case "nats":
		r = NewRpcNats(
			WithNatsEndpoints(cfg.Endpoints),
			WithNatsDialTimeout(timeout),
			WithNatsOptions(nats.Timeout(timeout)),
			WithNatsServer(server),
			WithNatsPrefix(cfg.Prefix),
			WithNatsDebugMsg(cfg.DebugMsg),
		)
	case "rabbitmq":
		r = NewRpcRabbitMq(
			WithRabbitMqEndpoints(cfg.Endpoints),
			WithRabbitMqDialTimeout(timeout),
			WithRabbitMqServer(server),
			WithRabbitMqPrefix(cfg.Prefix),
			WithRabbitMqDebugMsg(cfg.DebugMsg),
		)
	default:
		logger.Fatal("NewRpcConnector failed")
	}
	return r
}
