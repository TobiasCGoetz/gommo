package main

type Direction int

const (
	North Direction = iota
	East
	South
	West
	Stay
)

var Directions = [5]Direction{North, East, South, West, Stay}

func (d Direction) toString() string {
	return []string{"North", "East", "South", "West", "Stay"}[d]
}

var directions = map[string]Direction{
	"north": North,
	"east":  East,
	"south": South,
	"west":  West,
	"stay":  Stay,
}
