package main

import (
	"fmt"
)

func main() {
	var forest Tile = Tile.Forest
	var food Card = Card.Food
	var me Player = Player{"first", 13, 42, Card.Food, Card.Wood, Card.Wood, Card.None}
	fmt.Println(me.String())
}
