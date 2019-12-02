package core

import "fmt"

const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

type AOIManager struct {
	MinX  int           // 区域左边界坐标
	MaxX  int           // 区域有边界坐标
	CntsX int           // x方向格子的数量
	MinY  int           // 区域上边界坐标
	MaxY  int           // 区域下边界坐标
	CntsY int           // y方向的格子数量
	Grids map[int]*Grid // 当前区域中都有哪些格子 * key  = 格子ID, value = 格子对象
}

func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		Grids: make(map[int]*Grid),
	}

	// 给AOI初始化区域中所有的格子
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 计算格子ID
			// 格子编号 : id = idy * nx + idx (利用格子坐标得到格子编号)
			gid := y*cntsX + x

			// 初始化一个格子放在AOI中的map里， key是当前格子的ID
			aoiMgr.Grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}

	return aoiMgr
}

// 得到每个格子在x轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.CntsX
}

func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.CntsY
}

// 打印信息方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, cntsX:%d, minY:%d, maxY:%d, cntsY:%d\n Grids in AOI Manager:\n",
		m.MinX, m.MaxX, m.CntsX, m.MinY, m.MaxY, m.CntsY)
	for _, grid := range m.Grids {
		s += fmt.Sprintln(grid)
	}

	return s
}
