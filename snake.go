package main

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

type Point struct {
	//TODO: make this float, so its easier to scale with speeds and FPS
	x int
	y int
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
