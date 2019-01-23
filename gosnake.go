package main

import (
	"fmt"
	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/nsf/termbox-go"
	"math/rand"
	"time"
)

//make this variable
const table_height = 30
const table_width = 30
const FPS = 7
const body_gain = 3
const fruit_spawn_timer = 4000
const head_char = 'O'
const body_char = 'o'
const wall_char = 'X'
const fruit_char = '0'

const game_title = "GO-SNAKE"

const (
	easy       int8 = 1
	medium     int8 = 2
	hard       int8 = 3
	impossible int8 = 4
)

var game_difficulty int8 = easy

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

// game state indicates the current state of the game
//	-1 -> game aborted/closed
//	0	-> game not started
//	1	-> game started
//	2	-> game paused
//	3	-> game lost
const (
	st_closed      int8 = -2
	st_not_started int8 = 0
	st_running     int8 = 1
	st_paused      int8 = -1
	st_lost        int8 = -3
)

var game_state int8 = st_not_started

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
	score     int
}

func (s *Snake) resetSnake() {
	s.speed = 1
	s.score = 0
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

func (s *Snake) equalPoints(p *Point) bool {
	if (*s).body[0].x == p.x && (*s).body[0].y == p.y {
		return true
	} else {
		return false
	}
}

func (s *Snake) checkCollisionBody() bool {
	if len(s.body) > 2 {
		for _, pos := range s.body[2:] {
			// collision with own body
			if s.equalPoints(&pos) {
				return true
			}
		}
	}
	return false
}

func (s *Snake) addToBody(p Point, numb int) {
	for i := 0; i < numb; i++ {
		s.body = append(s.body, p)
	}
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
	case s.body[0].y >= table_height-1:
		return true
	case s.body[0].y <= 0:
		return true
	case s.body[0].x >= table_width-1:
		return true
	case s.body[0].x <= 0:
		return true
	default:
		return false
	}
}

func (s *Snake) checkCollisionFruit(f *Fruit) {
	for i, pos := range *f {
		if pos.x == s.body[0].x && pos.y == s.body[0].y {
			*f = append((*f)[:i], (*f)[i+1:]...)
			s.addToBody(s.body[len(s.body)-1], body_gain)
			s.score = s.score + (1 * int(game_difficulty))
		}
	}
}

//Game Table of default 50x50 Size
type Game_table [table_height][table_width]byte

func (table *Game_table) updateGameTable(s Snake, fruit *Fruit) {
	for x := range *table {
		for y := range (*table)[x] {
			// left wall
			if x == 0 {
				(*table)[x][y] = wall_char
				// right wall
			} else if x == table_height-1 {
				(*table)[x][y] = wall_char
				// top wall
			} else if y == 0 {
				(*table)[x][y] = wall_char
				// botoom wall
			} else if y == table_width-1 {
				(*table)[x][y] = wall_char
				// fill rest with spaces
				//TODO replace with current player and objects in the map
			} else {
				(*table)[x][y] = ' '
				for i := 1; i < len(s.body); i++ {
					(*table)[s.body[i].y][s.body[i].x] = body_char
				}
				for _, f := range *fruit {
					(*table)[f.y][f.x] = fruit_char
				}
				(*table)[s.body[0].y][s.body[0].x] = head_char
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
	game_state = st_running
	game_difficulty = easy

	// seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	//init stuff
	var table Game_table
	var s Snake
	var f Fruit
	s.resetSnake()
	table.resetGameTable()
	f.spawnFruit()
	table.drawTable(&s)

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
		if game_state != st_lost {
			select {
			case dir := <-c:
				switch dir {
				//pause
				case st_paused:
					if game_state == st_paused {
						game_state = st_running
					} else if game_state == st_running {
						game_state = st_paused

					}
					update(&s, &table, &f)
					//abort
				case st_closed:
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
		} else {
			break
		}
	}
	ticker_fps.Stop()

	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveTo(0, 0))
	fmt.Printf("\n\n\t\tYOU LOST\n\t\t---------\n\t\tScore: %v\n\n", s.score)
	fmt.Printf("\nPress any key to close the game...")
	time.Sleep(10 * time.Millisecond)
	termbox.PollEvent()

}

func update(s *Snake, table *Game_table, f *Fruit) {
	if game_state == st_running {
		s.move()
		//collision detection -> if false ==> collision detected
		if s.checkCollisionWall() || s.checkCollisionBody() {
			game_state = st_lost
		}
		s.checkCollisionFruit(f)
		table.updateGameTable(*s, f)
		table.drawTable(s)
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
			} else if ev.Ch == 'p' {
				c <- st_paused
			} else if ev.Key == termbox.KeyEsc {
				c <- st_closed
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

func (table Game_table) drawTable(s *Snake) {
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveTo(0, 0))
	fmt.Printf("\t\t%v\nScore: %v\n", game_title, s.score)

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
