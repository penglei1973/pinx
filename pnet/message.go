package pnet

type Message struct {
	Id      uint32 // 消息的ID
	DataLen uint32 //消息的长度
	Data    []byte // 消息的内容
}

func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// 获取消息数据段的长度
func (msg *Message) GetDataLen() uint32 {
	return msg.DataLen
}

// 获取消息ID
func (msg *Message) GetMsgId() uint32 {
	return msg.Id
}

// 获取消息内容
func (msg *Message) GetData() []byte {
	return msg.Data
}

func (msg *Message) SetMsgId(Id uint32) {
	msg.Id = Id
}

func (msg *Message) SetData(Data []byte) {
	msg.Data = Data
}

func (msg *Message) SetDataLen(DataLen uint32) {
	msg.DataLen = DataLen
}
