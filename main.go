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
	for a, player := range playerList {
		switch player.dir {
			case North:
				playerList[a].y += 1
			case East:
				playerList[a].x += 1
			case South:
				playerList[a].y -= 1
			case West:
				playerList[a].x -= 1
		}
		if mapWidth <= playerList[a].x {
			playerList[a].x = mapWidth-1
		}
		if playerList[a].x < 0 {
			playerList[a].x = 0
		}
		if mapHeight <= playerList[a].y {
			playerList[a].y = mapHeight-1
		}
		if playerList[a].y < 0 {
			playerList[a].y = 0
		}
		//Reset move direction per player
		player.dir = Stay
	}
}

func resources() {
	for _, player := range playerList {
		var firstEmpty = 5
		//Find first empty card space
		for f, card := range player.cards {
			if card == None {
				firstEmpty = f
			}
		}
		//Add card from tile
		switch gameMap[player.x][player.y] {
			case Forest:
				player.cards[firstEmpty] = Wood
				//Only add 2. Wood if there's space
				for f, card := range player.cards  {
					if card == None {
						player.cards[f] = Wood
					}
				}
			case City:
				player.cards[firstEmpty] = Weapon
			case Farm:
				player.cards[firstEmpty] = Food
			case Laboratory:
				player.cards[firstEmpty] = Research
		}
	}
}

func getHandSize(player Player) int {
	var count = 0
	for _, card := range player.cards {
		if card != None {
			count++
		}
	}
	return count
}

func limitCards() {
	for _, player := range playerList {
		if getHandSize(player) > 4 {
			if player.discard == None {
				fmt.Println("Cheater:")
				fmt.Println(player.id)
			}
		}
		for f, card := range player.cards {
			if card == player.discard {
				player.cards[f] = None
			}
		}
		player.discard = None
	}
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
	resources()
	limitCards()
	group()
	fight()
	spread()
}


func main() {
	rand.Seed(time.Now().UnixNano())
	var me Player = Player{"me", 2, 5, North, Weapon, None, [5]Card{Food, Wood, Wood, None, None}}
	playerList = append(playerList, me)
	initMap(&gameMap)
	printMap(&gameMap)
	for i := 0; i < 100; i++ {
		tick()
	}
}
