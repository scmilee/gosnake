package main

import (
	"fmt"
	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/nsf/termbox-go"
	"time"
)

const table_height = 50
const table_width = 50
const FPS = 5

const head_char = 'O'
const body_char = 'o'
const wall_char = 'X'

const game_title = "GO-SNAKE GAME"

// enums for snake movement
/*
		0
	3		1
		2
*/
const (
	up    int8 = 0
	right int8 = 1
	down  int8 = 2
	left  int8 = 3
)

type point struct {
	//TODO: make this float, so its easier to scale with speeds and FPS
	x int8
	y int8
}

type Snake struct {
	//TODO: maybe add custom heda/body_char for multiple snakes?
	head_position point
	body          [][]point
	//speed is a value between 0 and 100
	speed     uint8
	direction int8
}

func (s *Snake) resetSnake() {
	s.head_position = point{table_height / 2, table_width / 2}
	s.speed = 1
	s.direction = right
}

func (s *Snake) move() {
	switch s.direction {
	//move up
	case 0:
		s.head_position.x = s.head_position.x - 1
	//move right
	case 1:
		s.head_position.y = s.head_position.y + 1
	//move down
	case 2:
		s.head_position.x = s.head_position.x + 1
	//move left
	case 3:
		s.head_position.y = s.head_position.y - 1
	}
}

func (s Snake) checkMove(dir int8) bool {
	switch dir {
	case 0:
		if s.direction == 0 || s.direction == 2 {
			return false
		} else {
			return true
		}
	case 1:
		if s.direction == 1 || s.direction == 3 {
			return false
		} else {
			return true
		}
	case 2:
		if s.direction == 2 || s.direction == 0 {
			return false
		} else {
			return true
		}
	case 3:
		if s.direction == 3 || s.direction == 1 {
			return false
		} else {
			return true
		}
	}
	return true
}

//Game Table of default 50x50 Size
type Game_table [table_height][table_width]byte

func (table *Game_table) updateGameTable(s Snake) {
	for x := range table {
		for y := range table[x] {
			// left wall
			if x == 0 {
				table[x][y] = wall_char

				// right wall
			} else if x == table_height-1 {
				table[x][y] = wall_char

				// top wall
			} else if y == 0 {
				table[x][y] = wall_char

				// botoom wall
			} else if y == table_width-1 {
				table[x][y] = wall_char

				// fill rest with spaces
				//TODO replace with current player and objects in the map
			} else {
				table[x][y] = ' '
				table[s.head_position.x][s.head_position.y] = head_char

			}
		}
	}
}

func main() {
	//this has to called first to init and finalize key eventing
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	fmt.Printf("\t\t%v\n", game_title)

	//init stuff
	var table Game_table
	table.resetGameTable()
	table.drawTable()

	var s Snake
	s.resetSnake()
	//Tick Intervall s = 600 msecs/ 30 FPS
	ticker := time.NewTicker(600 / FPS * time.Millisecond)

	c := make(chan int8)

	go checkKeyInput(c)

mainloop:
	for range ticker.C {
		select {
		case dir := <-c:
			switch dir {
			//abort
			case -1:
				break mainloop
				// up
			case 0:
				if s.checkMove(dir) == true {
					s.direction = dir
				}
				update(&s, &table)
				// right
			case 1:
				if s.checkMove(dir) == true {
					s.direction = dir
				}
				update(&s, &table)
				// down
			case 2:
				if s.checkMove(dir) == true {
					s.direction = dir
				}
				update(&s, &table)
				//left
			case 3:
				if s.checkMove(dir) == true {
					s.direction = dir
				}
				update(&s, &table)
			}
		default:
			update(&s, &table)
		}
	}
}

func update(s *Snake, table *Game_table) {
	s.move()
	table.updateGameTable(*s)
	table.drawTable()
}

func checkKeyInput(c chan int8) {
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyArrowUp {
				c <- 0
			} else if ev.Key == termbox.KeyArrowRight {
				c <- 1
			} else if ev.Key == termbox.KeyArrowDown {
				c <- 2
			} else if ev.Key == termbox.KeyArrowLeft {
				c <- 3
			} else if ev.Key == termbox.KeyF2 {
				c <- -1
				break
			} else if ev.Key == termbox.KeyEsc {
				c <- -1
				break
			}
		}
	}
}

func (table Game_table) drawTable() {
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveTo(0, 0))
	fmt.Printf("\t\t%v\n", game_title)

	for x := range table {
		for y := range table[x] {
			fmt.Printf("%c", table[x][y])
		}
		fmt.Println()
	}
}

func (table *Game_table) resetGameTable() {
	for x := range table {
		for y := range table[x] {
			// left wall
			if x == 0 {
				table[x][y] = wall_char

				// right wall
			} else if x == table_height-1 {
				table[x][y] = wall_char

				// top wall
			} else if y == 0 {
				table[x][y] = wall_char

				// botoom wall
			} else if y == table_width-1 {
				table[x][y] = wall_char

				// fill rest with spaces
			} else {
				table[x][y] = ' '
			}
		}
	}
}
