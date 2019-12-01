package pinterface

type IRouter interface {
	PreHandle(request IRequest)  // 处理conn业务之前的hook的方法
	Handle(request IRequest)     // 处理conn业务方法
	PostHandle(request IRequest) //处理完conn的hook
}
