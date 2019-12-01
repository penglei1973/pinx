package pnet

import (
	"fmt"
	"pinx/pinterface"
	"pinx/utils"
	"strconv"
)

type MsgHandle struct {
	Apis           map[uint32]pinterface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan pinterface.IRequest
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]pinterface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan pinterface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(request pinterface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgid = ", request.GetConnection().GetConnID(), " is not Found")
		return
	}

	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router pinterface.IRouter) {
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api, msgid = " + strconv.Itoa(int(msgId)))
	}

	mh.Apis[msgId] = router
	fmt.Println("add api msgid = ", msgId)
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan pinterface.IRequest) {
	fmt.Println("worker ID = ", workerID, " is started.")
	// 不断的等待队列中的消息
	for {
		select {
		// 有消息则去除队列的Request 并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {
	// 遍历需要启动worker数量，依次启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan pinterface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		// 启动当前worker， 阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

func (mh *MsgHandle) SendMsgToTaskQueue(request pinterface.IRequest) {
	// 根据connid来分配当前的连接应该由哪个worker负责
	// 轮询平均分配

	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Printf("Add connid = %d\nrequest msgid = %d\n to worker id = %d\n", request.GetConnection().GetConnID(), request.GetMsgId(), workerID)
	// 将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}
