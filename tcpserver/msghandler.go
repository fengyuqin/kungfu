package tcpserver

import (
	"fmt"
	"github.com/jqiris/kungfu/config"
	"strconv"
)

// MsgHandle -
type MsgHandle struct {
	Apis           map[int32]IHandler   //存放每个MsgID 所对应的处理方法的map属性
	WorkerPoolSize int32                //业务工作Worker池的数量
	TaskQueue      []chan IRequest      //Worker负责取任务的消息队列
	Config         config.ConnectorConf //客户端配置
}

//NewMsgHandle 创建MsgHandle
func NewMsgHandle(cfg config.ConnectorConf) *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[int32]IHandler),
		WorkerPoolSize: cfg.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan IRequest, cfg.WorkerPoolSize),
		Config:    cfg,
	}
}

//SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	//fmt.Println("Add ConnID=", request.GetConnection().GetConnID()," request msgID=", request.GetMsgID(), "to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}

//DoMsgHandler 马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is not FOUND!")
		return
	}

	//执行对应处理方法
	handler(request)
}

//AddRouter 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID int32, handler IHandler) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		panic("repeated api , msgID = " + strconv.Itoa(int(msgID)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgID] = handler
	fmt.Println("Add api msgID = ", msgID)
}

//AddRouters 批量注册消息
func (mh *MsgHandle) AddRouters(routers map[int32]IHandler) {
	for msgID, handler := range routers {
		mh.Apis[msgID] = handler
	}
}

//StartOneWorker 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan IRequest) {
	//fmt.Println("Worker ID = ", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//StartWorkerPool 启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan IRequest, mh.Config.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
