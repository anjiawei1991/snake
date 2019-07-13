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
	snakeBodyStyle, snakeBoxStyle, foodStyle, snakeBoxStyleDead tcell.Style
)

func init() {
	snakeBoxStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorBlack)
	snakeBoxStyleDead = tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorBlack)
	snakeBodyStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	foodStyle = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
}

func drawSnake(s tcell.Screen, x, y int, snake *snake) {
	pts := snake.getAll()

	for _, pt := range pts {
		screenX, screenY := pt.x*2+x, pt.y+y
		s.SetContent(screenX, screenY, tcell.RuneBlock, nil, snakeBodyStyle)
		s.SetContent(screenX+1, screenY, tcell.RuneBlock, nil, snakeBodyStyle)
	}
}

func drawFood(s tcell.Screen, x, y int, food *point) {
	screenX, screenY := food.x*2+x, food.y+y
	s.SetContent(screenX, screenY, tcell.RuneBlock, nil, foodStyle)
	s.SetContent(screenX+1, screenY, tcell.RuneBlock, nil, foodStyle)
}

func drawBox(s tcell.Screen, x, y int, w, h int, style tcell.Style) {
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			s.SetContent(x+i, y+j, ' ', nil, style)
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

	model := newModel()
	boxX, boxY, boxW, boxH := 20, 5, 100, 50

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
		s.Fill(' ', tcell.StyleDefault.Background(tcell.ColorGray))

		if model.dead {
			drawBox(s, boxX, boxY, boxW, boxH, snakeBoxStyleDead)
		} else {
			drawBox(s, boxX, boxY, boxW, boxH, snakeBoxStyle)
		}
		drawSnake(s, boxX, boxY, &model.snake)
		drawFood(s, boxX, boxY, &model.food)

		if model.dead {
			snakeBoxStyle = snakeBoxStyle.Background(tcell.ColorRed)
		}
	}
}
