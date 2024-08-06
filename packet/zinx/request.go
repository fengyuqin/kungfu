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

package zinx

import (
	"github.com/fengyuqin/kungfu/v2/tcpface"
)

type Request struct {
	agent *Agent
	msg   *Message
}

func (r *Request) GetConnID() int {
	return r.agent.connId
}

func (r *Request) GetMsgID() int32 {
	return r.msg.Id
}

func (r *Request) GetMsgData() []byte {
	return r.msg.Data
}

func (r *Request) GetConnection() tcpface.IConnection {
	return r.agent
}

func (r *Request) GetServerID() string {
	return r.agent.server.GetServerID()
}
