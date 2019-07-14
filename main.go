package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

var (
	snakeBodyStyle, snakeDeadStyle, snakeBoxStyle, snakeBorderStyle, foodStyle tcell.Style
)

func init() {
	snakeBorderStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	snakeBoxStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)
	snakeBodyStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	snakeDeadStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
	foodStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
}

func drawSnake(s tcell.Screen, x, y int, snake *snake, dead bool) {
	pts := snake.getAll()

	style := snakeBodyStyle
	if dead {
		style = snakeDeadStyle
	}

	for i, pt := range pts {
		if i == 0 {
			r := '@'
			if dead {
				r = 'x'
			}
			s.SetContent(pt.x+x+1, pt.y+y+1, r, nil, style)
		} else {
			s.SetContent(pt.x+x+1, pt.y+y+1, 'o', nil, style)
		}
	}
}

func drawFood(s tcell.Screen, x, y int, food *point) {
	screenX, screenY := food.x+x+1, food.y+y+1
	s.SetContent(screenX, screenY, '$', nil, foodStyle)
}

func drawBox(s tcell.Screen, x, y int, w, h int) {
	s.SetContent(x+0, y+0, '┌', nil, snakeBorderStyle)
	s.SetContent(x+w-1, y+0, '┐', nil, snakeBorderStyle)
	s.SetContent(x+0, y+h-1, '└', nil, snakeBorderStyle)
	s.SetContent(x+w-1, y+h-1, '┘', nil, snakeBorderStyle)

	for i := 1; i < w-1; i++ {
		s.SetContent(x+i, y+0, '-', nil, snakeBorderStyle)
		s.SetContent(x+i, y+h-1, '-', nil, snakeBorderStyle)
	}

	for j := 1; j < h-1; j++ {
		s.SetContent(x+0, y+j, '▒', nil, snakeBorderStyle)
		s.SetContent(x+w-1, y+j, '▒', nil, snakeBorderStyle)
	}
}

func clearBox(s tcell.Screen, x, y int, w, h int) {
	for i := 1; i < w-1; i++ {
		for j := 1; j < h-1; j++ {
			s.SetContent(x+i, y+j, ' ', nil, snakeBoxStyle)
		}
	}
}

func keyDirection(key *tcell.EventKey) direction {
	rn := key.Rune()
	if rn == 'w' || rn == 'W' {
		return up
	}
	if rn == 's' || rn == 'S' {
		return down
	}
	if rn == 'd' || rn == 'D' {
		return right
	}
	if rn == 'a' || rn == 'A' {
		return left
	}
	return unknow
}

func main() {
	encoding.Register()

	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e := s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	defer s.Fini()

	eventC := make(chan tcell.Event, 32)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				eventC <- s.PollEvent()
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	boxX, boxY, boxW, boxH := 10, 10, 70, 26
	model := newModel(boxW-2, boxH-2)
	drawBox(s, boxX, boxY, boxW, boxH)

	for {
		s.Show()

		select {
		case ev := <-eventC:
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.Sync()
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEsc {
					return
				}
				r := ev.Rune()
				if r == 'r' || r == 'R' {
					model.restart()
				} else {
					d := keyDirection(ev)
					if d != unknow {
						model.processMove(d)
					}
				}
			}

		case <-ticker.C:
			model.processTimer()
		}
		clearBox(s, boxX, boxY, boxW, boxH)
		drawSnake(s, boxX, boxY, &model.snake, model.dead)
		drawFood(s, boxX, boxY, &model.food)
	}
}
