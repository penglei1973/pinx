package online

import (
	"fmt"
	"pinx/game_server/core"
	"pinx/pinterface"
)

func Online(conn pinterface.IConnection) {
	player := core.NewPlayer(conn)
	player.SyncPid()
	//同步当前玩家的初始化坐标信息给客户端，走MsgID:200消息
	player.BroadCastStartPosition()
	core.WorldMgrObj.AddPlayer(player)
	conn.SetProperty("pid", player.Pid)
	//==============同步周边玩家上线信息，与现实周边玩家信息========
	player.SyncSurrounding()
	fmt.Println("=====> Player pidId = ", player.Pid, " arrived====")
}
