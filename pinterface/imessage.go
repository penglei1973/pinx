package pinterface

type IMessage interface {
	GetDataLen() uint32 // 获取消息数据段的长度
	GetMsgId() uint32   // 获取消息ID
	GetData() []byte    // 获取消息内容

	SetMsgId(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
