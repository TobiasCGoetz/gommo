package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

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

func printPlayers (gMap *[mapWidth][mapHeight]*Tile, pList *[]*Player) {
	for  i:=0; i < mapWidth; i++ {
		for j:=0; j < mapHeight; j++ {
			var coordsFound = false
			for _, player := range *pList {
				if player.x == i && j == player.y {
					coordsFound = true
				}
			}
			if coordsFound {
				fmt.Printf("X|")
			} else {
				fmt.Printf("%c|", gMap[i][j].terrain.toString()[0])
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

func createCityList (gMap *[mapWidth][mapHeight]*Tile) []IntTuple  {
	var cities []IntTuple
	fmt.Println("createCityList()")
	for a, column := range gMap {
		for b, tile := range column {
			if tile.terrain == City {
				var coordinates = IntTuple{ a, b }
				cities = append(cities, coordinates)
			}
		}
	}
	return cities
}

func move(pList *[]*Player) {
	//Set new coordinates per player from move
	for a, player := range *pList {
		if !player.alive {
			continue
		}
		switch player.direction {
			case North:
				(*pList)[a].y += 1
			case East:
				(*pList)[a].x += 1
			case South:
				(*pList)[a].y -= 1
			case West:
				(*pList)[a].x -= 1
			case Stay:
				break
		}
		if mapWidth <= (*pList)[a].x {
			(*pList)[a].x = mapWidth-1
		}
		if (*pList)[a].x < 0 {
			(*pList)[a].x = 0
		}
		if mapHeight <= (*pList)[a].y {
			(*pList)[a].y = mapHeight-1
		}
		if (*pList)[a].y < 0 {
			(*pList)[a].y = 0
		}
		//Reset move direction per player
		(*pList)[a].direction = Stay
	}
}

func getFirstEmptyHandSlot(hand [5]Card) int {
	var firstEmpty = -1
	for f, card := range hand {
		if card == None {
			return f
		}
	}
	return firstEmpty
}

func resources(pList *[]*Player, gMap *[mapWidth][mapHeight]*Tile) {
	for pNr, player := range *pList {
		if !player.alive {
			continue
		}
		//TODO: Function to find first empty card space (reuse below)
		var firstEmpty = getFirstEmptyHandSlot(player.cards)
		//Add card from tile
		if firstEmpty > -1 {
			//printHandCards(*playerList[0])
			switch gMap[player.x][player.y].terrain {
				case Forest:
					(*pList)[pNr].cards[firstEmpty] = Wood
					firstEmpty = getFirstEmptyHandSlot(player.cards)
					if firstEmpty > -1 {
						(*pList)[pNr].cards[firstEmpty] = Wood
					}
				case City:
					player.cards[firstEmpty] = Weapon
				case Farm:
					player.cards[firstEmpty] = Food
				case Laboratory:
					player.cards[firstEmpty] = Research
			}
			//printHandCards(*playerList[0])
		}
	}
}

func consume(pList *[]*Player) {
	for a, player := range *pList {
		var playerCards = getHandSize(*player)
		if !player.alive {
			continue
		}
		if (*pList)[a].consume == None {
			(*pList)[a].alive = false
		} else {
			b, hasCard := playerHasCard(player, player.consume)
			if hasCard {
				(*pList)[a].cards[b] = None
			} else {
				(*pList)[a].alive = false
			}
		}
		var playerCards2 = getHandSize(*player)
		if playerCards == playerCards2 {
			fmt.Println("ERROR: Consumed cards should've been:", player.consume.toString())
			fmt.Println("ERROR: No card has been removed.")
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

func limitCards(pList *[]*Player) {
	for a, player := range *pList {
		if getHandSize(*player) > 4 {
			if player.discard == None {
				(*pList)[a].cards[4] = None
			}
		} else {
			for f, card := range player.cards {
				if card == player.discard && card != None {
					fmt.Printf(card.toString())
					fmt.Printf("\n")
					(*pList)[a].cards[f] = None
				}
			}
		}
		player.discard = None
	}
}

func handleCombat(gMap *[mapWidth][mapHeight]*Tile, pList *[]*Player) {
	//Create groups from position
	var combatGroups = make(map[IntTuple][]*Player)
	for _, player := range *pList {
		var pos = IntTuple{ player.x, player.y }
		combatGroups[pos] = append(combatGroups[pos], player)
	}
	for _, group := range combatGroups {
		fight(gMap, group)
	}
}

//TODO: Reevaluate when call by value is okay (argument is not altered)
func fight(gMap *[mapWidth][mapHeight]*Tile, group []*Player) {
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
	if attackValue < gMap[x][y].zombies {
		for a, _ := range group {
			group[a].alive = false
		}
		gMap[x][y].zombies += len(group)
	} else {
		gMap[x][y].zombies = 0
	}

}

func spread(gMap *[mapWidth][mapHeight]*Tile, cities *[]IntTuple) {
	for _, city := range *cities {
		//North
		if city.y < mapHeight-1 {
			gMap[city.x][city.y+1].zombies++
		}
		//East
		if city.x < mapWidth-1 {
			gMap[city.x+1][city.y].zombies++
		}
		//South
		if city.y > 0 {
			gMap[city.x][city.y-1].zombies++
		}
		//West
		if city.x > 0 {
			gMap[city.x-1][city.y].zombies++
		}
	}
	//Maybe cutoff?
}

func tick(gMap *[mapWidth][mapHeight]*Tile, cities *[]IntTuple, pList *[]*Player) {
	move(pList)
	resources(pList, gMap)
	handleCombat(gMap, pList)
	spread(gMap, cities)
	consume(pList)
	limitCards(pList)
}

func playerHasCard (player *Player, card Card) (int, bool) {
	for a, c := range player.cards {
		if c == card {
			return a, true
		}
	}
	return -1, false
}

func randomizeBot(players []*Player) {
	for _, player := range players {
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
}

func addBot(players *[]*Player, bots *[]*Player, bID int) {
	var rX = rand.Intn(mapWidth-1)
	var rY = rand.Intn(mapHeight-1)
	var bot = Player{
		id:        strconv.Itoa(bID),
		x:         rX,
		y:         rY,
		direction: Stay,
		play:      None,
		consume:   None,
		discard:   None,
		cards:     [5]Card{ Food, Wood, Wood, None, None },
		alive:     true,
		isBot:     true,
	}
	*players = append(*players, &bot)
	*bots = append(*bots, &bot)
}

func restockBots(players *[]*Player, bots *[]*Player, bID *int) {
	var botDelta = botNumber - len(*bots)
	for i := 0; i < botDelta; i++ {
		addBot(players, bots, *bID)
		*bID++
	}
}

//TODO: Handle dead players correctly
func main() {
	var gameMap [mapWidth][mapHeight]*Tile
	var cityList []IntTuple
	var playerList []*Player
	var botList []*Player
	var botID = 0
	var isRunning = true
	rand.Seed(time.Now().UnixNano())
	initMap(&gameMap)
	go setupAPI()
	cityList = createCityList(&gameMap)
	for isRunning {
		restockBots(&playerList, &botList, &botID)
		randomizeBot(botList)
		tick(&gameMap, &cityList, &playerList)
	}
}
