package core

import (
	"fmt"
	"sync"
)

type Grid struct {
	GID       int          // 格子ID
	MinX      int          // 格子左边界坐标
	MaxX      int          // 格子右边界坐标
	MinY      int          // 格子上边界坐标
	MaxY      int          // 格子下边界坐标
	playerIDs map[int]bool //当前格子内的玩家或者物体成员ID
	pIDLock   sync.RWMutex // playerIDs的保护map的锁
}

// 初始化一个格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 向当前格子中添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

// 从格子中删除一个玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// 得到当前格子中所有的玩家
func (g *Grid) GetPlyerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

// 打印信息方法
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d,maxY:%d, playerIDs:%v", g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}

//根据格子的gID得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundGridsByGid(gID int) (grids []*Grid) {
	//判断gID是否存在
	if _, ok := m.Grids[gID]; !ok {
		return
	}
	//将当前gid添加到九宫格中
	grids = append(grids, m.Grids[gID])
	//根据gid得到当前格子所在的X轴编号
	idx := gID % m.CntsX
	//判断当前idx左边是否还有格子
	if idx > 0 {
		grids = append(grids, m.Grids[gID-1])
	}
	//判断当前的idx右边是否还有格子
	if idx < m.CntsX-1 {
		grids = append(grids, m.Grids[gID+1])
	}
	//将x轴当前的格子都取出，进行遍历，再分别得到每个格子的上下是否有格子
	//得到当前x轴的格子id集合
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}
	//遍历x轴格子
	for _, v := range gidsX {
		//计算该格子处于第几列
		idy := v / m.CntsX
		//判断当前的idy上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.Grids[v-m.CntsX])
		}
		//判断当前的idy下边是否还有格子
		if idy < m.CntsY-1 {
			grids = append(grids, m.Grids[v+m.CntsX])
		}
	}
	return
}

//通过横纵坐标获取对应的格子ID
func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(x) - m.MinY) / m.gridLength()
	return gy*m.CntsX + gx
}

//通过横纵坐标得到周边九宫格内的全部PlayerIDs
func (m *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	//根据横纵坐标得到当前坐标属于哪个格子ID
	gID := m.GetGIDByPos(x, y)
	//根据格子ID得到周边九宫格的信息
	grids := m.GetSurroundGridsByGid(gID)
	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlyerIDs()...)
		fmt.Printf("===> grid ID : %d, pids : %v ====", v.GID,
			v.GetPlyerIDs())
	}
	return
}

//通过GID获取当前格子的全部playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.Grids[gID].GetPlyerIDs()
	return
}

//移除一个格子中的PlayerID
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.Grids[gID].Remove(pID)
}

//添加一个PlayerID到一个格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.Grids[gID].Add(pID)
}

//通过横纵坐标添加一个Player到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.Grids[gID]
	grid.Add(pID)
}

//通过横纵坐标把一个Player从对应的格子中删除
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.Grids[gID]
	grid.Remove(pID)
}
