package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"github.com/ahmetalpbalkan/go-cursor"
	"github.com/nsf/termbox-go"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//make this variable
var table_height int = 0
var table_width int = 0

var FPS int = 5
var body_gain int = 3
var fruit_spawn_timer = 4000

const standard_width = 40
const standard_height = 20
const head_char = 'O'
const body_char = 'o'
const wall_char = 'X'
const fruit_char = '0'

var game_title = "GO-SNAKE"

const (
	easy       int8 = 1
	medium     int8 = 2
	hard       int8 = 3
	impossible int8 = 4
)

type Point struct {
	x int
	y int
}

var game_difficulty int8 = easy

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

//Game Table of default 50x50 Size
type Game_table [][]byte

func (table *Game_table) updateGameTable(s Snake, fruit *Fruit) {
	for y := range *table {
		for x := range (*table)[y] {
			// left wall
			if x == 0 {
				(*table)[y][x] = wall_char
				// right wall
			} else if x == table_width-1 {
				(*table)[y][x] = wall_char
				// top wall
			} else if y == 0 {
				(*table)[y][x] = wall_char
				// botoom wall
			} else if y == table_height-1 {
				(*table)[y][x] = wall_char
				// fill rest with spaces
				//TODO replace with current player and objects in the map
			} else {
				(*table)[y][x] = ' '
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
	
func check(e error) {
	if e != nil {
			panic(e)
	}
}

func main() {
	//get command line flags
	//	-difficulty=...
	//	-size=AxB
	diffFlag := flag.String("difficulty", "medium", "This sets the difficulty for the game. Possible difficulties are: easy, medium, hard, impossible. Default is medium.")
	sizeFlag := flag.String("size", "30x30", "This sets the size of the game table. It has to be formatted like this: 'WIDTH'x'HEIGHT'!")
	flag.Parse()

	if flag.Parsed() {
		switch *diffFlag {
		case "easy":
			game_difficulty = easy

		case "medium":
			game_difficulty = medium

		case "hard":
			game_difficulty = hard

		case "impossible":
			game_difficulty = impossible
		default:
			game_difficulty = medium
		}

		sizeSlice := strings.Split(*sizeFlag, "x")
		if len(sizeSlice) > 1 {
			width, err1 := strconv.Atoi(sizeSlice[0])
			height, err2 := strconv.Atoi(sizeSlice[1])
			if err1 != nil || err2 != nil {
				table_width = standard_width
				table_height = standard_height
			} else {
				if width > 0 && height > 0 {
					table_width = width
					table_height = height
				} else {
					table_width = standard_width
					table_height = standard_height
				}
			}
		} else {
			table_width = standard_width
			table_height = standard_height
		}
	} else {
		panic("error")
	}

	//this has to called first to init and finalize key eventing
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	game_state = st_running

	switch game_difficulty {
	case easy:
		body_gain = 1
		FPS = 4
		game_title += " - EASY"
	case medium:
		body_gain = 3
		FPS = 6
		game_title += " - MEDIUM"
	case hard:
		body_gain = 10
		FPS = 8
		game_title += " - HARD"
	case impossible:
		body_gain = 15
		FPS = 10
		game_title += " - IMPOSSIBLE"
	}

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
	ticker_fps := time.NewTicker(time.Duration(600/FPS) * time.Millisecond)

	/*
		sends/receives different signals
			 0 - change direction to: up
			 1 - change direction to: right
			 2 - change direction to: down
			 3 - change direction to: left
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
	time.Sleep(100 * time.Millisecond)
	updateHighScore(&s)
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
	fruit_ticker := time.NewTicker(time.Duration(fruit_spawn_timer) * time.Millisecond)
	for range fruit_ticker.C {
		fruit_chan <- true
	}
	fruit_ticker.Stop()
}

func (table Game_table) drawTable(s *Snake) {
	fmt.Print(cursor.ClearEntireScreen())
	fmt.Print(cursor.MoveTo(0, 0))
	fmt.Printf("\t\t%v\nHighScore: %s\nScore: %v\n", game_title,string(s.highScore),s.score)

	for x := range table {
		for y := range table[x] {
			fmt.Printf("%c", table[x][y])
		}
		fmt.Println()
	}
}

func (table *Game_table) resetGameTable() {
	//set size of Game_table
	*table = make([][]byte, table_height)
	for i := range *table {
		(*table)[i] = make([]byte, table_width)
	}
}

func updateHighScore(s *Snake) {
	intScore, err := strconv.Atoi(string(s.highScore))
	check(err)
	if s.score > intScore {
		d1 := []byte(fmt.Sprintf("%v", s.score))
		errd1 := ioutil.WriteFile("./highscore", d1, 0644)
		check(errd1)
	}
}
