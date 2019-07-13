package main

type direction int

const (
	unknow direction = 0
	up     direction = 1
	down   direction = 2
	left   direction = 3
	right  direction = 4
)

type point struct {
	x, y int
}

func nextPos(pt point, d direction) point {
	switch d {
	case up:
		pt.y--
	case down:
		pt.y++
	case left:
		pt.x--
	case right:
		pt.x++
	}
	return pt
}

// 获取从 pt2 到 pt1 的方向(需要相临)
func getDirection(pt1, pt2 point) direction {
	for _, d := range []direction{up, down, left, right} {
		if pt1 == nextPos(pt2, d) {
			return d
		}
	}
	return unknow
}

type snake struct {
	points []point // 第 1 个坐标为头，中间坐标为拐点，最后一个坐标是尾巴
}

// grow 增加长度
func (m *snake) grow(d direction) bool {
	// 相反方向成长，不做处理
	if len(m.points) > 1 &&
		nextPos(m.points[0], d) == m.points[1] {
		return false
	}
	m.points = append(m.points, point{})
	copy(m.points[1:], m.points[:len(m.points)-1])
	m.points[0] = nextPos(m.points[0], d)
	return true
}

// 朝指定方向移动一格
// 先朝指定方向成长一格，再将尾巴去掉
func (m *snake) move(d direction) {
	if m.grow(d) {
		m.points = m.points[:len(m.points)-1]
	}
}

// getAll 获取所有的点，从头到尾部排列
func (m *snake) getAll() []point {
	return m.points
}
