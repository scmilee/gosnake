package main

import (
	"math/rand"
)

type Fruit []Point

func (f *Fruit) spawnFruit() {
	p := Point{1 + rand.Intn(table_height-1), 1 + rand.Intn(table_width-1)}
	*f = append(*f, p)
}
