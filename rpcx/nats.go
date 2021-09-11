package rpcx

import (
	"path"
	"strings"
	"time"

	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/treaty"
	"github.com/jqiris/kungfu/utils"
	"github.com/nats-io/nats.go"
)

const (
	Balancer  = "balancer"
	Connector = "connector"
	Server    = "server"
	Database  = "database"
)

type RpcNats struct {
	Endpoints   []string
	Options     []nats.Option
	Client      *nats.Conn
	DialTimeout time.Duration
	RpcCoder    *RpcEncoder
	Server      *treaty.Server
	DebugMsg    bool
	Prefix      string
}
type RpcNatsOption func(r *RpcNats)

func WithNatsDebugMsg(debug bool) RpcNatsOption {
	return func(r *RpcNats) {
		r.DebugMsg = debug
	}
}
func WithNatsEndpoints(endpoints []string) RpcNatsOption {
	return func(r *RpcNats) {
		r.Endpoints = endpoints
	}
}
func WithNatsDialTimeout(timeout time.Duration) RpcNatsOption {
	return func(r *RpcNats) {
		r.DialTimeout = timeout
	}
}
func WithNatsServer(server *treaty.Server) RpcNatsOption {
	return func(r *RpcNats) {
		r.Server = server
	}
}
func WithNatsOptions(opts ...nats.Option) RpcNatsOption {
	return func(r *RpcNats) {
		r.Options = opts
	}
}
func WithNatsPrefix(prefix string) RpcNatsOption {
	return func(r *RpcNats) {
		r.Prefix = prefix
	}
}

func NewRpcNats(opts ...RpcNatsOption) *RpcNats {
	r := &RpcNats{
		Prefix: "RpcX",
	}
	for _, opt := range opts {
		opt(r)
	}
	url := strings.Join(r.Endpoints, ",")
	conn, err := nats.Connect(url, r.Options...)
	if err != nil {
		logger.Fatal(err)
	}
	r.Client = conn
	r.RpcCoder = NewRpcEncoder()
	return r
}

func (r *RpcNats) Subscribe(server *treaty.Server, callback CallbackFunc) error {
	sub := path.Join(r.Prefix, treaty.RegSeverItem(server))
	if _, err := r.Client.Subscribe(sub, func(msg *nats.Msg) {
		go utils.SafeRun(func() {
			r.DealMsg(msg, callback)
		})
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) QueueSubscribe(queue string, server *treaty.Server, callback CallbackFunc) error {
	sub := path.Join(r.Prefix, treaty.RegSeverItem(server))
	if _, err := r.Client.QueueSubscribe(sub, queue, func(msg *nats.Msg) {
		go utils.SafeRun(func() {
			r.DealMsg(msg, callback)
		})
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) SubscribeBalancer(callback CallbackFunc) error {
	sub := path.Join(r.Prefix, Balancer)
	if _, err := r.Client.Subscribe(sub, func(msg *nats.Msg) {
		go utils.SafeRun(func() {
			r.DealMsg(msg, callback)
		})
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) SubscribeConnector(callback CallbackFunc) error {
	sub := path.Join(r.Prefix, Connector)
	if _, err := r.Client.Subscribe(sub, func(msg *nats.Msg) {
		go utils.SafeRun(func() {
			r.DealMsg(msg, callback)
		})
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) SubscribeServer(callback CallbackFunc) error {
	sub := path.Join(r.Prefix, Server)
	if _, err := r.Client.Subscribe(sub, func(msg *nats.Msg) {
		go utils.SafeRun(func() {
			r.DealMsg(msg, callback)
		})
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) SubscribeDatabase(callback CallbackFunc) error {
	sub := path.Join(r.Prefix, Database)
	if _, err := r.Client.Subscribe(sub, func(msg *nats.Msg) {
		utils.SafeRun(func() {
			r.DealMsg(msg, callback)
		})
	}); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) DealMsg(msg *nats.Msg, callback CallbackFunc) {
	req := &RpcMsg{}
	err := r.RpcCoder.Decode(msg.Data, req)
	if err != nil {
		logger.Error(err)
		return
	}
	resp := callback(r, req)
	if resp != nil {
		if err := msg.Respond(resp); err != nil {
			logger.Error(err)
		}
	}
	if r.DebugMsg {
		logger.Infof("DealMsg,msgType: %v, msgId: %v", req.MsgType, req.MsgId)
	}
}

func (r *RpcNats) Request(server *treaty.Server, msgId int32, req, resp interface{}) error {
	var msg *nats.Msg
	var err error
	var data []byte
	data, err = r.EncodeMsg(Request, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverItem(server))
	if msg, err = r.Client.Request(sub, data, r.DialTimeout); err == nil {
		respMsg := &RpcMsg{MsgData: resp}
		err = r.RpcCoder.Decode(msg.Data, respMsg)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func (r *RpcNats) Publish(server *treaty.Server, msgId int32, req interface{}) error {
	data, err := r.EncodeMsg(Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, treaty.RegSeverItem(server))
	if err = r.Client.Publish(sub, data); err != nil {
		return err
	}
	return nil
}

func (r *RpcNats) PublishBalancer(msgId int32, req interface{}) error {
	data, err := r.EncodeMsg(Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Balancer)
	return r.Client.Publish(sub, data)
}

func (r *RpcNats) PublishConnector(msgId int32, req interface{}) error {
	data, err := r.EncodeMsg(Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Connector)
	return r.Client.Publish(sub, data)
}

func (r *RpcNats) PublishServer(msgId int32, req interface{}) error {
	data, err := r.EncodeMsg(Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Server)
	return r.Client.Publish(sub, data)
}

func (r *RpcNats) PublishDatabase(msgId int32, req interface{}) error {
	data, err := r.EncodeMsg(Publish, msgId, req)
	if err != nil {
		return err
	}
	sub := path.Join(r.Prefix, Database)
	return r.Client.Publish(sub, data)
}

func (r *RpcNats) EncodeMsg(msgType MessageType, msgId int32, req interface{}) ([]byte, error) {
	rpcMsg := &RpcMsg{
		MsgType: msgType,
		MsgId:   msgId,
		MsgData: req,
	}
	data, err := r.RpcCoder.Encode(rpcMsg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RpcNats) DecodeMsg(data []byte, v interface{}) error {
	return r.RpcCoder.DecodeMsg(data, v)
}

func (r *RpcNats) GetCoder() *RpcEncoder {
	return r.RpcCoder
}

func (r *RpcNats) Response(v interface{}) []byte {
	return r.RpcCoder.Response(v)
}

func (r *RpcNats) GetServer() *treaty.Server {
	return r.Server
}
