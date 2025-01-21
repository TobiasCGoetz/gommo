package main

import (
	"fmt"
)

func main() {
	var forest Tile = Forest
	var food Card = Food
	var me Player = Player{"first", 13, 42, [4]Card{Food, Wood, Wood, None}}
	fmt.Println(forest.toString())
	fmt.Println(food.toString())
	fmt.Println(me.toString())
}
