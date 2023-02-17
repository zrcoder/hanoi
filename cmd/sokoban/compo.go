package main

import (
	_ "fmt"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

const (
	Wall  = '#'
	Me    = '@'
	Slot  = 'X'
	Box   = 'O'
	BoxIn = '*'
	Blank = ' '
)

type sokoban struct {
	level int
	x     int
	y     int
}

// TODO
func New() *sokoban {
	return nil
}

/*
func (g *Game) whereami() {
	for y, cels := range g.board {
		for x, cel := range cels {
			if cel.obj == GUY {
				g.x, g.y = x, y
			}
		}
	}
}

func (g *Game) checkMove(dx, dy int) {
	var (
		x, y int
		k    = 1
	)
	for {
		x, y = g.x+dx*k, g.y+dy*k
		cell := g.board[y][x]
		switch cell.obj {
		case FLOOR, SLOT:
			// move obj along the way foward
			for k > 0 {
				xp, yp := x-dx, y-dy
				g.board[y][x].obj = g.board[yp][xp].obj
				x, y = xp, yp
				k--
			}
			g.board[y][x].obj = g.board[y][x].base
			g.x, g.y = x+dx, y+dy
			return
		case WALL:
			return
		case BOX:
			if g.board[y+dy][x+dx].obj == BOX {
				return
			}
		}
		k++
	}
}

func (g *Game) move(dir Direction) {
	switch dir {
	case UP:
		g.checkMove(0, -1)
	case DOWN:
		g.checkMove(0, 1)
	case LEFT:
		g.checkMove(-1, 0)
	case RIGHT:
		g.checkMove(1, 0)
	}
	done := g.checkState()
	if done {
		g.nextLevel()
	}
}

func (g *Game) checkState() bool {
	for _, cells := range g.board {
		for _, cell := range cells {
			if cell.base == SLOT && cell.obj != BOX {
				return false
			}
		}
	}
	return true
}

func (g *Game) nextLevel() {
	if g.level < g.db.maxLevel {
		g.level = g.level + 1
		g.reset()
	}
}

func (g *Game) prevLevel() {
	if g.level > 1 {
		g.level = g.level - 1
		g.reset()
	}
}

func (g *Game) toggleDebug() {
	g.debug = !g.debug
}

// func main() {
// 	var p = func(g *Game) {
// 		for i, row := range g.board {
// 			fmt.Printf("%-2d %q\n", i, row)
// 		}
// 		fmt.Println()
// 	}
//
// 	g := NewGame()
// 	p(g)
// 	g.move(RIGHT)
// 	p(g)
// }
*/
