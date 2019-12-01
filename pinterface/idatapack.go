package pinterface

type IDataPack interface {
	GetHeadLen() uint32                // 获取包头长度的方法
	Pack(msg IMessage) ([]byte, error) //封包
	Unpack([]byte) (IMessage, error)   //拆包
}
