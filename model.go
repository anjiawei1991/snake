package main

import (
	"math/rand"
	"time"
)

type model struct {
	snake         snake
	food          point
	width, height int
	dead          bool
	score         int
	level         int
	lastMoveTime  time.Time
}

func newModel(w, h int) *model {
	c := model{
		width:  w,
		height: h,
	}
	c.init()
	return &c
}

func (c *model) init() {
	rand.Seed(time.Now().Unix())
	c.snake = snake{points: []point{{c.width / 2, c.height / 2}, {c.width/2 - 1, c.height / 2}}}
	c.dead = false
	c.level = 1
	c.score = 0
	c.lastMoveTime = time.Now()
	c.randomFood()
}

func (c *model) randomFood() {
	c.food = point{
		rand.Intn(c.width), rand.Intn(c.height),
	}
}

func (c *model) restart() {
	c.init()
}

func (c *model) processMove(d direction) {
	if c.dead {
		return
	}
	all := c.snake.getAll()
	head := all[0]
	pt := nextPos(head, d)
	if len(all) >= 2 && pt == all[1] {
		return
	}
	if pt.x < 0 || pt.x >= c.width || pt.y < 0 || pt.y >= c.height {
		c.dead = true
		return
	}

	for _, spt := range all {
		if spt == pt {
			c.dead = true
			return
		}
	}

	if c.food == pt {
		c.snake.grow(d)
		c.score++
		c.level = c.score/5 + 1
		c.randomFood()
	} else {
		c.snake.move(d)
	}

	c.lastMoveTime = time.Now()
}

func (c *model) processTick() {
	if time.Now().Sub(c.lastMoveTime) < time.Second/time.Duration(c.level) {
		return
	}

	d := right
	all := c.snake.getAll()
	if len(all) >= 2 {
		d = getDirection(all[0], all[1])
	}

	c.processMove(d)
}
