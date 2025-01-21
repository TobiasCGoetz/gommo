package main

import (
	"fmt"
	"math/rand"
	"time"
)


var gameMap [mapWidth][mapHeight]Tile
var playerList []Player

func initMap (gameMap *[mapWidth][mapHeight]Tile) {
	for a, column := range gameMap {
		for b, _ := range column {
			choice := rand.Intn(len(tileTypes))
			gameMap[a][b] = tileTypes[choice]
		}
	}
}

func printMap (gameMap *[mapWidth][mapHeight]Tile) {
	for a, row := range gameMap {
		for b, _ := range row {
			fmt.Printf("%c|", gameMap[a][b].toString()[0])
		}
		fmt.Printf("\n")
	}
}

func move() {
	//Set new coordinates per player from move
	for _, player := range playerList {
		switch player.dir {
		case North:
			player.y += 1
		case East:
			player.x += 1
		case South:
			player.y -= 1
		case West:
			player.x -= 1
		}
	//Reset move direction per player
	player.dir = Stay
	}
}

func ressources() {
	//Add card from tile
	//Handle cutoff/selection/blocking
}

func group() {
	//Create groups from position
}

func fight() {
	//Calculate dice + weapon VS zombies per group
}

func spread() {
	//Get all cities
	//Increment neighbours
	//Maybe cutoff
}

func tick() {
	move()
	ressources()
	group()
	fight()
	spread()
}


func main() {
	rand.Seed(time.Now().UnixNano())
	var me Player = Player{"me", 13, 42, North, Weapon, [4]Card{Food, Wood, Wood, None}}
	playerList = append(playerList, me)
	initMap(&gameMap)
	printMap(&gameMap)
	for i := 0; i < 100; i++ {
		tick()
	}
}
