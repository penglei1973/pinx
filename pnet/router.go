package pnet

import "pinx/pinterface"

// 嵌入基类， 重写方法
type BaseRouter struct {
}

// 处理conn业务之前的hook的方法
func (br *BaseRouter) PreHandle(request pinterface.IRequest) {

}

// 处理conn业务方法
func (br *BaseRouter) Handle(request pinterface.IRequest) {

}

//处理完conn的hook
func (br *BaseRouter) PostHandle(request pinterface.IRequest) {

}
