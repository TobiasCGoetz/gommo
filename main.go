package main

import (
	"fmt"
	"math/rand"
	"time"
)

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
	//Reset move direction per player
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

func printSlice(s []Player) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	//var forest Tile = Forest
	//var food Card = Food
	var me Player = Player{"me", 13, 42, North, Weapon, [4]Card{Food, Wood, Wood, None}}
	//fmt.Println(forest.toString())
	//fmt.Println(food.toString())
	//fmt.Println(me.toString())
	var gameMap [mapWidth][mapHeight]Tile
	var playerList []Player
	printSlice(playerList)
	playerList = append(playerList, me)
	printSlice(playerList)
	initMap(&gameMap)
	printMap(&gameMap)
}
