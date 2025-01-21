package main

import (
	"fmt"
	"math/rand"
)

func initMap (gameMap [mapWidth][mapHeight]Tile) {
	for a, column := range gameMap {
		for b, _ := range column {
			choice := rand.Intn(tileTypes.len())
			gameMap[a][b] = tileTypes[choice]
		}
	}
}

func printMap (gameMap [mapWidth][mapHeight]Tile) {
	for a, row := range gameMap {
		for b, _ := range row {
			fmt.Printf("%c|", gameMap[a][b].toString()[0])
		}
		fmt.Printf("\n")
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	//var forest Tile = Forest
	//var food Card = Food
	//var me Player = Player{"first", 13, 42, [4]Card{Food, Wood, Wood, None}}
	//fmt.Println(forest.toString())
	//fmt.Println(food.toString())
	//fmt.Println(me.toString())
	var gameMap [mapWidth][mapHeight]Tile
	initMap(gameMap)
	printMap(gameMap)
}
