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
}

func newModel(w, h int) *model {
	c := model{
		snake:  snake{points: []point{{w / 2, h / 2}, {w/2 - 1, h / 2}}},
		width:  w,
		height: h,
	}
	rand.Seed(time.Now().Unix())
	c.randomFood()
	return &c
}

func (c *model) randomFood() {
	c.food = point{
		rand.Intn(c.width), rand.Intn(c.height),
	}
}

func (c *model) restart() {
	c.snake = snake{points: []point{{c.width / 2, c.height / 2}, {c.width/2 - 1, c.height / 2}}}
	c.dead = false
	c.randomFood()
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
		c.randomFood()
	} else {
		c.snake.move(d)
	}
}

func (c *model) processTimer() {
	d := right
	all := c.snake.getAll()
	if len(all) >= 2 {
		d = getDirection(all[0], all[1])
	}

	c.processMove(d)
}
