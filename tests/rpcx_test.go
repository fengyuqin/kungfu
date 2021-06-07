package tests

import (
	"fmt"
	"testing"

	"github.com/jqiris/kungfu/conf"
	"github.com/jqiris/kungfu/rpcx"
	"github.com/jqiris/kungfu/treaty"
)

func TestRpc(t *testing.T) {
	cfg := config.GetRpcXConf()
	//gate
	s1 := &treaty.Server{
		ServerId:   "1001",
		ServerType: "string_Balancer",
		ServerName: "gate",
		ServerIp:   "127.0.0.1",
		ClientPort: 123,
	}
	w1 := rpcx.NewRpcBalancer(cfg)
	if err := w1.Subscribe(s1, func(req []byte) []byte {
		logger.Infof("gate received: %v", string(req))
		return []byte(fmt.Sprintf("gate received: %v", string(req)))
	}); err != nil {
		logger.Errorf("gate err:%v", err)
	}
	if err := w1.SubscribeBalancer(func(req []byte) []byte {
		logger.Infof("gate2 received: %v", string(req))
		return []byte(fmt.Sprintf("gate2 received: %v", string(req)))
	}); err != nil {
		logger.Errorf("gate2 err:%v", err)
	}
	//connector
	s2 := &treaty.Server{
		ServerId:   "1002",
		ServerType: "string_Connector",
		ServerName: "connector",
		ServerIp:   "127.0.0.1",
		ClientPort: 456,
	}
	w2 := rpcx.NewRpcConnector(cfg)
	if err := w2.Subscribe(s2, func(req []byte) []byte {
		logger.Infof("connector received: %v", string(req))
		return []byte(fmt.Sprintf("connector received: %v", string(req)))
	}); err != nil {
		logger.Errorf("connector err:%v", err)
	}
	if err := w2.SubscribeConnector(func(req []byte) []byte {
		logger.Infof("connector2 received: %v", string(req))
		return []byte(fmt.Sprintf("connector2 received: %v", string(req)))
	}); err != nil {
		logger.Errorf("connector2 err:%v", err)
	}
	//connector
	s3 := &treaty.Server{
		ServerId:   "1003",
		ServerType: "string_Game",
		ServerName: "game",
		ServerIp:   "127.0.0.1",
		ClientPort: 789,
	}
	w3 := rpcx.NewRpcServer(cfg)
	if err := w3.Subscribe(s3, func(req []byte) []byte {
		logger.Infof("server received: %v", string(req))
		return []byte(fmt.Sprintf("server received: %v", string(req)))
	}); err != nil {
		logger.Errorf("server err:%v", err)
	}
	if err := w3.SubscribeServer(func(req []byte) []byte {
		logger.Infof("server2 received: %v", string(req))
		return []byte(fmt.Sprintf("server2 received: %v", string(req)))
	}); err != nil {
		logger.Errorf("server2 err:%v", err)
	}

	//s1 请求 s2
	reply, err := w1.Request(s2, []byte("from gate"))
	logger.Infof("s1=>s2, reply:%v, err:%v", string(reply), err)
}
