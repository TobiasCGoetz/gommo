package main

import (
	"fmt"
	"math/rand"
	"time"
)


var gameMap [mapWidth][mapHeight]*Tile
var cityList []IntTuple
var playerList []*Player

func initMap (gMap *[mapWidth][mapHeight]*Tile) {
	for a, column := range gMap {
		for b, _ := range column {
			choice := rand.Intn(len(terrainTypes))
			var tile = Tile{ terrainTypes[choice], 0 }
			gMap[a][b] = &tile
		}
	}
}

func printMap (gMap *[mapWidth][mapHeight]*Tile) {
	for a, row := range gMap {
		for b, _ := range row {
			fmt.Printf("%c|", gMap[a][b].terrain.toString()[0])
		}
		fmt.Printf("\n")
	}
}

func printPlayers (pList *[]*Player) {
	for  i:=0; i < mapWidth; i++ {
		for j:=0; j < mapHeight; j++ {
			var coordsFound = false
			for _, player := range playerList {
				if player.x == i && j == player.y {
					coordsFound = true
				}
			}
			if coordsFound {
				fmt.Printf("X|")
			} else {
				fmt.Printf("%c|", gameMap[i][j].terrain.toString()[0])
			}
		}
		fmt.Printf("\n")
	}
}

func printHandCards (player Player) {
	fmt.Printf(player.id)
	fmt.Printf(": ")
	for _, card := range player.cards {
		fmt.Printf(card.toString())	
		fmt.Printf("|")
	}
	fmt.Printf("\n")
}

func createCityList () {
	for a, column := range gameMap {
		for b, tile := range column {
			if tile.terrain == City {
				var coordinates = IntTuple{ a, b }
				cityList = append(cityList, coordinates)
			}
		}
	}
}

func move() {
	//Set new coordinates per player from move
	for a, player := range playerList {
		switch player.direction {
			case North:
				playerList[a].y += 1
			case East:
				playerList[a].x += 1
			case South:
				playerList[a].y -= 1
			case West:
				playerList[a].x -= 1
			case Stay:
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
		player.direction = Stay
	}
}

func resources() {
	for _, player := range playerList {
		//TODO: Function to find first empty card space (reuse below)
		var firstEmpty = -1
		for f, card := range player.cards {
			if card == None {
				firstEmpty = f
			}
		}
		//Add card from tile
		if firstEmpty > -1 {
					printHandCards(*playerList[0])
			switch gameMap[player.x][player.y].terrain {
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
			printHandCards(*playerList[0])
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
		if getHandSize(*player) > 4 {
			if player.discard == None {
				fmt.Println("Cheater:")
				fmt.Println(player.cards)
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

func handleCombat() {
	//Create groups from position
	var combatGroups = make(map[IntTuple][]*Player)
	for _, player := range playerList {
		var pos = IntTuple{ player.x, player.y }
		combatGroups[pos] = append(combatGroups[pos], player)
	}
	for _, group := range combatGroups {
		fight(group)
	}
}

func fight(group []*Player) {
	//Calculate dice + weapon VS zombies per group
	var attackValue = 0
	var x = group[0].x
	var y = group[0].y
	for a, player := range group {
		if player.play == Weapon {
			attackValue += 6
			removeLoop:for b, _ := range player.cards {
				if group[a].cards[b] == Weapon {
					group[a].cards[b] = None
					break removeLoop
				}
			}
		} else {
			attackValue += rand.Intn(6)
		}
	}
	if attackValue < gameMap[x][y].zombies {
		for a, _ := range group {
			group[a].alive = false
		}
		gameMap[x][y].zombies += len(group)
	} else {
		gameMap[x][y].zombies = 0
	}

}

func spread() {
	for _, city := range cityList {
		//North
		if city.y < mapHeight-1 {
			gameMap[city.x][city.y+1].zombies++
		}
		//East
		if city.x < mapWidth-1 {
			gameMap[city.x+1][city.y].zombies++
		}
		//South
		if city.y > 0 {
			gameMap[city.x][city.y-1].zombies++
		}
		//West
		if city.x > 0 {
			gameMap[city.x-1][city.y].zombies++
		}
	}
	//Maybe cutoff?
}

func tick() {
	move()
	resources()
	limitCards()
	handleCombat()
	spread()
}

func playerHasCard (player *Player, card Card) (int, bool) {
	for a, c := range player.cards {
		if c == card {
			return a, true
		}
	}
	return -1, false
}

func randomizePlayerInput(player *Player) {
	//Randomize movement
	player.direction = directions[rand.Intn(len(directions))]
	//Randomize card played
	player.play = Dice
	//Randomize consume
	a, found := playerHasCard(player, Food)
	player.consume = Food
	if !found {
		a, found = playerHasCard(player, Wood)
		player.consume = Wood
		if a == -1 {
			player.consume = None
		}
	}
	//Randomize discard
	_, found = playerHasCard(player, None)
	if !found {
		player.discard = player.cards[0]
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var me Player = Player{"me", 2, 5, North, Weapon, Food, None, [5]Card{Food, Wood, Wood, None, None}, true}
	playerList = append(playerList, &me)
	initMap(&gameMap)
	printMap(&gameMap)
	createCityList()
	for i := 0; i < 30; i++ {
		fmt.Print("\033[H\033[2J")
		printPlayers(&playerList)
		printHandCards(me)
		randomizePlayerInput(&me)
		tick()
		time.Sleep(time.Second/2)
	}
}
