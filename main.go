package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
)

var gameMap [mapWidth][mapHeight]*Tile
var playerMap = make(map[string]*Player)
var botList []*Player
var r *rand.Rand

func initMap(r rand.Rand, gMap *[mapWidth][mapHeight]*Tile) {
	fmt.Println("Initializing game map...")
	for a, column := range gMap {
		for b := range column {
			choice := r.Intn(len(terrainTypes) - 1)
			gMap[a][b] = &Tile{terrainTypes[choice], 0, []string{}}
		}
	}
}

func getMapTile(x int, y int, gMap *[mapWidth][mapHeight]*Tile) *Tile {
	if x < 0 || x >= mapWidth || y < 0 || y >= mapHeight {
		fmt.Printf("Prevented tile access at %d/%d", x, y)
		return &Tile{Edge, -1, []string{}}
	}
	var truncX = x % 100
	var truncY = y % 100
	return (*gMap)[truncX][truncY]
}

func move(playerMap *map[string]*Player) {
	//Set new coordinates per player from move
	for _, player := range *playerMap {
		if !player.Alive {
			continue
		}
		//Fetch current player state
		var player = (*playerMap)[player.ID]

		//Perform move
		switch player.Direction {
		case North:
			player.Y += 1
		case East:
			player.X += 1
		case South:
			player.Y -= 1
		case West:
			player.X -= 1
		case Stay:
			return
		}

		//Write new coordinates
		(*playerMap)[player.ID] = player

		//Prevent out-of-map moves
		if mapWidth <= player.X {
			player.X = mapWidth - 1
		}
		if player.X < 0 {
			player.X = 0
		}
		if mapHeight <= player.Y {
			player.Y = mapHeight - 1
		}
		if player.Y < 0 {
			player.Y = 0
		}
		//Reset move direction
		player.Direction = defaultDirection
		//Write new state
		(*playerMap)[player.ID] = player
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

func resources(playerMap *map[string]*Player, gMap [mapWidth][mapHeight]*Tile) {
	for playerID := range *playerMap {
		var player = (*playerMap)[playerID]
		if !player.Alive {
			continue
		}
		var firstEmpty = getFirstEmptyHandSlot(player.Cards)
		//Add card from tile
		if firstEmpty > -1 {
			switch gMap[player.X][player.Y].Terrain {
			case Forest: //Special case, rewards 2 cards
				player.Cards[firstEmpty] = Wood
				firstEmpty = getFirstEmptyHandSlot(player.Cards)
				if firstEmpty > -1 {
					player.Cards[firstEmpty] = Wood
				}
			case City:
				player.Cards[firstEmpty] = Weapon
			case Farm:
				player.Cards[firstEmpty] = Food
			case Laboratory:
				player.Cards[firstEmpty] = Research
			}
		}
		(*playerMap)[playerID] = player
	}
}

func consume(playerMap *map[string]*Player, gMap *[mapWidth][mapHeight]*Tile) {
	for playerID := range *playerMap {
		//Fetch current player state
		var player = (*playerMap)[playerID]
		if !player.Alive {
			continue
		}

		//We don't allow death by indecision
		if player.Consume == None {
			_, hasCard := playerHasCard(player, Food)
			if hasCard {
				player.Consume = Food
			} else {
				player.Consume = Wood
			}
		}

		//Now remove that card or kill the player
		cardPos, hasCard := playerHasCard(player, player.Consume)
		if hasCard {
			if player.Consume == Wood {

				var zombiesAttracted = 0

				var tileNW = getMapTile(player.X-1, player.Y+1, gMap)
				var tileNN = getMapTile(player.X, player.Y+1, gMap)
				var tileNE = getMapTile(player.X+1, player.Y+1, gMap)
				var tileWW = getMapTile(player.X-1, player.Y, gMap)
				var tileEE = getMapTile(player.X+1, player.Y, gMap)
				var tileSW = getMapTile(player.X-1, player.Y-1, gMap)
				var tileSS = getMapTile(player.X, player.Y-1, gMap)
				var tileSE = getMapTile(player.X+1, player.Y-1, gMap)

				var tileArray = []*Tile{tileNW, tileNN, tileNE, tileWW, tileEE, tileSW, tileSS, tileSE}

				//Remove zombies from surrounding tiles
				for _, nextTile := range tileArray {
					if nextTile.Zombies > 0 {
						zombiesAttracted++
						nextTile.Zombies--
					}
				}

				//Add to players tile
				getMapTile(player.X, player.Y, gMap).Zombies += zombiesAttracted
			}
			player.Cards[cardPos] = None //Remove card from hand
			(*playerMap)[player.ID] = player
		} else {
			player.Alive = false //Card not in hand, kill the player
			(*playerMap)[player.ID] = player
		}
	}
}

func getHandSize(player *Player) int {
	var count = 0
	for _, card := range player.Cards {
		if card != None {
			count++
		}
	}
	return count
}

func limitCards(playerMap *map[string]*Player) {
	for mapKey := range *playerMap {
		var player = (*playerMap)[mapKey]
		if getHandSize(player) > 4 {
			var cardPos, hasCard = playerHasCard(player, player.Discard)
			if hasCard && player.Discard != None && cardPos > -1 { //Better safe...
				player.Cards[cardPos] = None
			} else {
				player.Cards[4] = None
			}
		}
		player.Discard = None
		(*playerMap)[mapKey] = player
	}
}

func handleCombat() {
	var wg = sync.WaitGroup{}
	for x, _ := range gameMap {
		for y, _ := range gameMap[x] {
			wg.Add(1)
			go tileWorker(gameMap[x][y], &wg)
		}
	}
	wg.Wait()
}

func spreadFromSpreader(gMap *[mapWidth][mapHeight]*Tile, xCoord int, yCoord int) {
	// TODO: decide if spread is 4 or 8 directions
	var xOffsets = []int{0, -1, 0, 1, 0}
	var yOffsets = []int{-1, 0, 0, 0, 1} //TODO: Check y-axis direction again!
	for neighbor := 0; neighbor < len(xOffsets); neighbor++ {
		var xTarget = xCoord + xOffsets[neighbor]
		var yTarget = yCoord + yOffsets[neighbor]
		if xTarget < 0 || xTarget >= mapWidth || yTarget < 0 || yTarget >= mapHeight {
			continue
		}
		gMap[xTarget][yTarget].spreadTo()
	}
}

func spread(gMap *[mapWidth][mapHeight]*Tile) {
	for x, _ := range gMap {
		for y, tile := range gMap[x] {
			if tile.isSpreader() {
				spreadFromSpreader(gMap, x, y)
			}
		}
	}
}

// TODO: Unify order of attributes across functions
func tick(gMap *[mapWidth][mapHeight]*Tile, playerMap *map[string]*Player) {
	fmt.Println("Next tick is happening...")
	fmt.Println("Moving players...")
	move(playerMap)
	fmt.Println("Distributing ressources...")
	resources(playerMap, *gMap)
	fmt.Println("Combat is upon us...")
	handleCombat()
	fmt.Println("The infection is spreading...")
	spread(gMap)
	fmt.Println("Players feeding themselves...")
	consume(playerMap, gMap)
	fmt.Println("Limiting player inventory")
	limitCards(playerMap)
}

func playerHasCard(player *Player, card Card) (int, bool) {
	for a, c := range player.Cards {
		if c == card {
			return a, true
		}
	}
	return -1, false
}

func randomizeBots(bots []*Player) {
	fmt.Println("Randomizing bot turns...")
	for _, bot := range bots {
		//Randomize movement
		bot.Direction = Directions[r.Intn(len(Directions))]
		//Randomize card played
		bot.Play = Dice
		//Randomize consume

		if _, foodFound := playerHasCard(bot, Food); foodFound {
			bot.Consume = Food
		} else if _, woodFound := playerHasCard(bot, Wood); woodFound {
			bot.Consume = Wood
		} else {
			bot.Consume = None
		}

		//Randomize discard
		_, found := playerHasCard(bot, None)
		if !found {
			bot.Discard = bot.Cards[0]
		}
	}
}

// TODO: Somehow remove inactive players
func addPlayer(playerName string) string {
	var rX = r.Intn(mapWidth - 1)
	var rY = r.Intn(mapHeight - 1)
	playerID, _ := uuid.NewV7()
	idString := playerID.String()
	var player = Player{
		ID:        idString,
		Name:      playerName,
		X:         rX,
		Y:         rY,
		Direction: defaultDirection,
		Play:      None,
		Consume:   None,
		Discard:   None,
		Cards:     [5]Card{Food, Wood, Wood, None, None},
		Alive:     true,
		IsBot:     false,
	}
	playerMap[idString] = &player
	return idString
}

func addBot(playerMap *map[string]*Player, bots *[]*Player, bID int) {
	var rX = r.Intn(mapWidth - 1)
	var rY = r.Intn(mapHeight - 1)
	var bot = Player{
		ID:        strconv.Itoa(bID),
		X:         rX,
		Y:         rY,
		Direction: Stay,
		Play:      None,
		Consume:   None,
		Discard:   None,
		Cards:     [5]Card{Food, Wood, Wood, None, None},
		Alive:     true,
		IsBot:     true,
	}
	(*playerMap)[strconv.Itoa(bID)] = &bot
	*bots = append(*bots, &bot)
}

func restockBots(playerMap *map[string]*Player, bots *[]*Player, bID *int) {
	var botDelta = botNumber - len(*bots)
	for i := 0; i < botDelta; i++ {
		addBot(playerMap, bots, *bID)
		*bID++
	}
}

func havePlayersWon(playerMap map[string]*Player) bool {
	for _, player := range playerMap {
		if player.hasWinCondition() {
			return true
		}
	}
	return false
}

func getPlayerOrNil(id string) *Player {
	return playerMap[id]
}

func getSurroundingsOfPlayer(id string) (Surroundings, bool) {
	player := getPlayerOrNil(id)
	if player == nil { //TODO: If nil else function or invert? Make them all identical!
		return Surroundings{}, false
	} else {
		var NW = gameMap[player.X-1][player.Y-1].getMapPiece()
		var NN = gameMap[player.X][player.Y-1].getMapPiece()
		var NE = gameMap[player.X+1][player.Y-1].getMapPiece()
		var WW = gameMap[player.X-1][player.Y].getMapPiece()
		var CE = gameMap[player.X][player.Y].getMapPiece()
		var EE = gameMap[player.X+1][player.Y].getMapPiece()
		var SW = gameMap[player.X-1][player.Y+1].getMapPiece()
		var SS = gameMap[player.X][player.Y+1].getMapPiece()
		var SE = gameMap[player.X+1][player.Y+1].getMapPiece()

		var miniMap = Surroundings{
			NW: NW,
			NN: NN,
			NE: NE,
			WW: WW,
			CE: CE,
			EE: EE,
			SW: SW,
			SS: SS,
			SE: SE,
		}
		return miniMap, true
	}
}

func printMap() {
	for x, _ := range gameMap {
		for _, tile := range gameMap[x] {
			fmt.Print(fmt.Sprintf("%s ", tile.Terrain.toChar()))
		}
		fmt.Print("\n")
		for _, tile := range gameMap[x] {
			if len(tile.playerIds) > 0 {
				fmt.Print(fmt.Sprintf("%d", len(tile.playerIds)))
			} else {
				fmt.Print(" ")
			}
			if tile.Zombies > 0 {
				fmt.Print(fmt.Sprintf("%d", tile.Zombies))
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func main() {
	if len(os.Args) == 2 {
		idSalt = os.Args[1]
		fmt.Println(idSalt)
	}
	var botID = 0
	var turnTimer = int8(turnLength)
	hasWon = false
	r = rand.New(rand.NewSource(time.Now().Unix()))
	initMap(*r, &gameMap)
	//Data set up, now we can start the API

	go setupAPI(&gameMap, &turnTimer, &hasWon)

	var remainingTurns = maxTurns
	for ; remainingTurns > 0; remainingTurns-- {
		fmt.Println("Remaining turns: ", remainingTurns)
		//printMap()
		for turnTimer = int8(turnLength); turnTimer >= 0; turnTimer-- {
			if turnTimer == 0 {
				randomizeBots(botList)
				tick(&gameMap, &playerMap)
				hasWon = havePlayersWon(playerMap)
				restockBots(&playerMap, &botList, &botID)
				if hasWon {
					fmt.Println("Game over due to win")
					remainingTurns = 0
					turnTimer = -1
					break
				}
			} else {
				time.Sleep(time.Second)
			}
		}
	}
	time.Sleep(120 * time.Second)
}
