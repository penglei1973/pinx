package pnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"pinx/pinterface"
	"pinx/utils"
)

// 封包拆包类实例
type DataPack struct {
}

// 封包拆包初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包头长度的方法
func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

//封包
func (dp *DataPack) Pack(msg pinterface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的buff
	dataBuff := bytes.NewBuffer([]byte{})

	// 写datalen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil

}

//拆包
func (dp *DataPack) Unpack(binaryData []byte) (pinterface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head的信息， 得到dataLen和msgID
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断datalen长度是否符合要求
	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too large msg")
	}

	return msg, nil
}
