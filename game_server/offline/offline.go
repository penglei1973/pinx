package offline

import (
	"fmt"
	"pinx/game_server/core"
	"pinx/pinterface"
)

func Offline(conn pinterface.IConnection) {
	//获取当前连接的Pid属性
	pid, _ := conn.GetProperty("pid")
	//根据pid获取对应的玩家对象
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	//触发玩家下线业务
	if pid != nil {
		player.LostConnection()
	}
	fmt.Println("====> Player ", pid, " left =====")
}
