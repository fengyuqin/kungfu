package rpc

import (
	"github.com/jqiris/kungfu/logger"
	"github.com/jqiris/kungfu/treaty"
)

func DefaultCallback(req *MsgRpc) []byte {
	logger.Info("DefaultCallback")
	return nil
}

type RssBuilder struct {
	queue    string
	server   *treaty.Server
	callback CallbackFunc
	codeType string
	suffix   string
	parallel bool
}

func NewRssBuilder(server *treaty.Server) *RssBuilder {
	return &RssBuilder{
		queue:    DefaultQueue,
		server:   server,
		callback: DefaultCallback,
		codeType: CodeTypeProto,
		suffix:   DefaultSuffix,
		parallel: true,
	}
}

func (r *RssBuilder) SetQueue(queue string) *RssBuilder {
	r.queue = queue
	return r
}

func (r *RssBuilder) SetServer(server *treaty.Server) *RssBuilder {
	r.server = server
	return r
}
func (r *RssBuilder) SetCallback(callback CallbackFunc) *RssBuilder {
	r.callback = callback
	return r
}
func (r *RssBuilder) SetCodeType(codeType string) *RssBuilder {
	r.codeType = codeType
	return r
}
func (r *RssBuilder) SetSuffix(suffix string) *RssBuilder {
	r.suffix = suffix
	return r
}
func (r *RssBuilder) SetParallel(parallel bool) *RssBuilder {
	r.parallel = parallel
	return r
}

func (r *RssBuilder) Build() RssBuilder {
	return RssBuilder{
		queue:    r.queue,
		server:   r.server,
		callback: r.callback,
		codeType: r.codeType,
		suffix:   r.suffix,
		parallel: r.parallel,
	}
}

type ReqBuilder struct {
	codeType   string
	suffix     string
	server     *treaty.Server
	msgId      int32
	req        interface{}
	resp       interface{}
	serverType string
}

func NewReqBuilder(server *treaty.Server) *ReqBuilder {
	serverType := ""
	if server != nil {
		serverType = server.ServerType
	}
	return &ReqBuilder{
		codeType:   CodeTypeProto,
		suffix:     DefaultSuffix,
		server:     server,
		serverType: serverType,
	}
}

func (r *ReqBuilder) SetCodeType(codeType string) *ReqBuilder {
	r.codeType = codeType
	return r
}
func (r *ReqBuilder) SetSuffix(suffix string) *ReqBuilder {
	r.suffix = suffix
	return r
}
func (r *ReqBuilder) SetServer(server *treaty.Server) *ReqBuilder {
	r.server = server
	return r
}
func (r *ReqBuilder) SetMsgId(msgId int32) *ReqBuilder {
	r.msgId = msgId
	return r
}
func (r *ReqBuilder) SetReq(req interface{}) *ReqBuilder {
	r.req = req
	return r
}
func (r *ReqBuilder) SetResp(resp interface{}) *ReqBuilder {
	r.resp = resp
	return r
}
func (r *ReqBuilder) SetServerType(serverType string) *ReqBuilder {
	r.serverType = serverType
	return r
}

func (r *ReqBuilder) Build() ReqBuilder {
	return ReqBuilder{
		codeType:   r.codeType,
		suffix:     r.suffix,
		server:     r.server,
		msgId:      r.msgId,
		req:        r.req,
		resp:       r.resp,
		serverType: r.serverType,
	}
}
