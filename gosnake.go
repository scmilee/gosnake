package main

import (
	"fmt"
	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

const table_height = 50
const table_width = 50
const FPS = 5
const fruit_spawn_timer = 4000

const head_char = 'O'
const body_char = 'o'
const wall_char = 'X'
const fruit_char = '0'

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

var pause = false

type Point struct {
	//TODO: make this float, so its easier to scale with speeds and FPS
	x int
	y int
}

type Fruit []Point

func (f *Fruit) spawnFruit() {
	p := Point{1 + rand.Intn(table_height-1), 1 + rand.Intn(table_width-1)}
	*f = append(*f, p)
}

type Snake struct {
	body []Point
	//speed is a value between 0 and 100
	speed     uint8
	direction int8
}

func (s *Snake) resetSnake() {
	s.speed = 1
	s.direction = right
	s.body = append(s.body, Point{table_height / 2, table_width / 2})
}

func (s *Snake) move() {
	last_move := make([]Point, len(s.body))
	copy(last_move, s.body)

	switch s.direction {
	//move up
	case up:
		s.body[0].y = s.body[0].y - 1
	//move right
	case right:
		s.body[0].x = s.body[0].x + 1
	//move down
	case down:
		s.body[0].y = s.body[0].y + 1
	//move left
	case left:
		s.body[0].x = s.body[0].x - 1
	}

	if len(s.body) > 1 {
		for i := 0; i < len(s.body); i++ {
			if i > 0 {
				s.body[i] = last_move[i-1]
			}
		}
	}
}

func (s *Snake) addToBody(p Point) {
	s.body = append(s.body, p)
}

// Check if the given dir is a valid direction to move
func (s Snake) checkMove(dir int8) bool {
	switch dir {
	case up:
		if s.direction == up || s.direction == down {
			return false
		} else {
			return true
		}
	case right:
		if s.direction == right || s.direction == left {
			return false
		} else {
			return true
		}
	case down:
		if s.direction == down || s.direction == up {
			return false
		} else {
			return true
		}
	case left:
		if s.direction == left || s.direction == right {
			return false
		} else {
			return true
		}
	}
	return true
}

func (s Snake) checkCollisionWall() bool {
	switch {
	case s.body[0].x >= table_height-1:
		return false
	case s.body[0].x <= 0:
		return false
	case s.body[0].y >= table_width-1:
		return false
	case s.body[0].y <= 0:
		return false
	default:
		return true
	}
}

func (s *Snake) checkCollisionFruit(f *Fruit) {
	for i, pos := range *f {
		if pos.x == s.body[0].x && pos.y == s.body[0].y {
			*f = append((*f)[:i], (*f)[i+1:]...)
			s.addToBody(s.body[len(s.body)-1])
		}
	}
}

//Game Table of default 50x50 Size
type Game_table [table_height][table_width]byte

func (table *Game_table) updateGameTable(s Snake, fruit *Fruit) {
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
				for i := 1; i < len(s.body); i++ {
					table[s.body[i].y][s.body[i].x] = body_char
				}
				for _, f := range *fruit {
					table[f.y][f.x] = fruit_char
				}
				table[s.body[0].y][s.body[0].x] = head_char
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
	var f Fruit
	s.resetSnake()
	//Tick Intervall s = 600 msecs/ 30 FPS
	ticker_fps := time.NewTicker(600 / FPS * time.Millisecond)
	/*
		sends/receives different signals
			 0 - change direction to: up
			 1 - change direction to: right
			 2 - change direction to: down
			 3 - change direction to: left
			 -1 - break mainloop aka. stop program
			 -2 - pause the game
	*/
	c := make(chan int8)
	fruit_spawn := make(chan bool)
	go checkKeyInput(c)
	go spawnFruitTicker(fruit_spawn)

mainloop:
	for range ticker_fps.C {
		select {
		case dir := <-c:
			switch dir {
			//pause
			case -2:
				pause = !pause
				update(&s, &table, &f)
				//abort
			case -1:
				break mainloop
				// up
			case up:
				if s.checkMove(dir) == true {
					s.direction = dir
				}
				update(&s, &table, &f)
				// right
			case right:
				if s.checkMove(dir) == true {
					s.direction = dir
				}
				update(&s, &table, &f)
				// down
			case down:
				if s.checkMove(dir) == true {
					s.direction = dir
				}
				update(&s, &table, &f)
				//left
			case left:
				if s.checkMove(dir) == true {
					s.direction = dir
				}
				update(&s, &table, &f)
			}
		case <-fruit_spawn:
			f.spawnFruit()
			update(&s, &table, &f)
		default:
			update(&s, &table, &f)
		}
	}
	ticker_fps.Stop()
}

func update(s *Snake, table *Game_table, f *Fruit) {
	if pause == false {
		s.move()
		//collision detection -> if false ==> collision detected
		if !s.checkCollisionWall() {
			pause = true
		}
		s.checkCollisionFruit(f)
		table.updateGameTable(*s, f)
		table.drawTable()
	}
}

func checkKeyInput(c chan int8) {
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			if ev.Key == termbox.KeyArrowUp {
				c <- up
			} else if ev.Key == termbox.KeyArrowRight {
				c <- right
			} else if ev.Key == termbox.KeyArrowDown {
				c <- down
			} else if ev.Key == termbox.KeyArrowLeft {
				c <- left
			} else if ev.Key == termbox.KeyCtrlP {
				c <- -2
			} else if ev.Key == termbox.KeyEsc {
				c <- -1
				break
			}
		}
	}
}

func spawnFruitTicker(fruit_chan chan bool) {
	fruit_ticker := time.NewTicker(fruit_spawn_timer * time.Millisecond)
	for range fruit_ticker.C {
		fruit_chan <- true
	}
	fruit_ticker.Stop()
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
