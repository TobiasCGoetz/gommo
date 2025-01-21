package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var r *rand.Rand

func initMap(r rand.Rand, gMap *[mapWidth][mapHeight]*Tile) {
	for a, column := range gMap {
		for b := range column {
			choice := r.Intn(len(terrainTypes) - 1)
			var tile = Tile{terrainTypes[choice], 0, []Player{}}
			gMap[a][b] = &tile
		}
	}
}

func createCityList(gMap *[mapWidth][mapHeight]*Tile) []IntTuple {
	var cities []IntTuple
	for a, column := range gMap {
		for b, tile := range column {
			if tile.Terrain == City {
				var coordinates = IntTuple{a, b}
				cities = append(cities, coordinates)
			}
		}
	}
	return cities
}

func getMapTile(x int, y int, gMap *[mapWidth][mapHeight]*Tile) *Tile {
	if x < 0 || x >= mapWidth || y < 0 || y >= mapHeight {
		fmt.Printf("Prevented tile access at %d/%d", x, y)
		return &Tile{Edge, -1, []Player{}}
	}
	var truncX = x % 100
	var truncY = y % 100
	return (*gMap)[truncX][truncY]
}

func move(pList *[]*Player) {
	//Set new coordinates per player from move
	for a, player := range *pList {
		if !player.Alive {
			continue
		}
		switch player.Direction {
		case North:
			(*pList)[a].Y += 1
		case East:
			(*pList)[a].X += 1
		case South:
			(*pList)[a].Y -= 1
		case West:
			(*pList)[a].X -= 1
		case Stay:
			return
		}
		if mapWidth <= (*pList)[a].X {
			(*pList)[a].X = mapWidth - 1
		}
		if (*pList)[a].X < 0 {
			(*pList)[a].X = 0
		}
		if mapHeight <= (*pList)[a].Y {
			(*pList)[a].Y = mapHeight - 1
		}
		if (*pList)[a].Y < 0 {
			(*pList)[a].Y = 0
		}
		//Reset move direction per player
		(*pList)[a].Direction = Stay
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
		if !player.Alive {
			continue
		}
		var firstEmpty = getFirstEmptyHandSlot(player.Cards)
		//Add card from tile
		if firstEmpty > -1 {
			//printHandCards(*playerList[0])
			switch gMap[player.X][player.Y].Terrain {
			case Forest:
				(*pList)[pNr].Cards[firstEmpty] = Wood
				firstEmpty = getFirstEmptyHandSlot(player.Cards)
				if firstEmpty > -1 {
					(*pList)[pNr].Cards[firstEmpty] = Wood
				}
			case City:
				player.Cards[firstEmpty] = Weapon
			case Farm:
				player.Cards[firstEmpty] = Food
			case Laboratory:
				player.Cards[firstEmpty] = Research
			}
			//printHandCards(*playerList[0])
		}
	}
}

func consume(pList *[]*Player, gMap *[mapWidth][mapHeight]*Tile) {
	for a, player := range *pList {
		var playerCards = getHandSize(*player)
		if !player.Alive {
			continue
		}
		if (*pList)[a].Consume == None {
			_, hasCard := playerHasCard(player, Food)
			if hasCard {
				(*pList)[a].Consume = Food
			} else {
				(*pList)[a].Consume = Wood
			}
		}
		b, hasCard := playerHasCard(player, player.Consume)
		if hasCard {
			if player.Consume == Wood {
				var zombiesAttracted = 0
				var tileNN = getMapTile(player.X, player.Y+1, gMap)
				var tileNE = getMapTile(player.X+1, player.Y+1, gMap)
				var tileEE = getMapTile(player.X+1, player.Y, gMap)
				var tileSE = getMapTile(player.X+1, player.Y-1, gMap)
				var tileSS = getMapTile(player.X, player.Y-1, gMap)
				var tileSW = getMapTile(player.X-1, player.Y-1, gMap)
				var tileWW = getMapTile(player.X-1, player.Y, gMap)
				var tileNW = getMapTile(player.X-1, player.Y+1, gMap)

				if tileNN.Zombies > 0 {
					zombiesAttracted++
					tileNN.Zombies -= 1
				}
				if tileNE.Zombies > 0 {
					zombiesAttracted++
					tileNE.Zombies -= 1
				}
				if tileEE.Zombies > 0 {
					zombiesAttracted++
					tileEE.Zombies -= 1
				}
				if tileSE.Zombies > 0 {
					zombiesAttracted++
					tileSE.Zombies -= 1
				}
				if tileSS.Zombies > 0 {
					zombiesAttracted++
					tileSS.Zombies -= 1
				}
				if tileSW.Zombies > 0 {
					zombiesAttracted++
					tileSW.Zombies -= 1
				}
				if tileWW.Zombies > 0 {
					zombiesAttracted++
					tileWW.Zombies -= 1
				}
				if tileNW.Zombies > 0 {
					zombiesAttracted++
					tileNW.Zombies -= 1
				}
				getMapTile(player.X, player.Y, gMap).Zombies += zombiesAttracted
			}
			(*pList)[a].Cards[b] = None
		} else {
			player.Alive = false
		}
		var playerCards2 = getHandSize(*player)
		if playerCards == playerCards2 && playerCards != 0 && player.Alive {
			fmt.Println("ERROR: Consumed cards should've been:", player.Consume.toString())
			fmt.Println("ERROR: No card has been removed from the", playerCards2, "in hand.")
			fmt.Println("ERROR: PLAYER:")
			fmt.Println(player)
		}
	}
}

func getHandSize(player Player) int {
	var count = 0
	for _, card := range player.Cards {
		if card != None {
			count++
		}
	}
	return count
}

func limitCards(pList *[]*Player) {
	for a, player := range *pList {
		if getHandSize(*player) > 4 {
			if player.Discard == None {
				(*pList)[a].Cards[4] = None
			} else {
				var cardPos, hasCard = playerHasCard((*pList)[a], player.Discard)
				if hasCard {
					(*pList)[a].Cards[cardPos] = None
				} else {
					(*pList)[a].Cards[4] = None
				}
			}
		}
		(*pList)[a].Discard = None
	}
}

func handleCombat(gMap *[mapWidth][mapHeight]*Tile, pList *[]*Player) {
	//Create groups from position
	var combatGroups = make(map[IntTuple][]*Player)
	for _, player := range *pList {
		var pos = IntTuple{player.X, player.Y}
		combatGroups[pos] = append(combatGroups[pos], player)
	}
	for _, group := range combatGroups {
		fight(gMap, group)
	}
}

func fight(gMap *[mapWidth][mapHeight]*Tile, group []*Player) {
	//Calculate dice + weapon VS Zombies per group
	var attackValue = 0
	var x = group[0].X
	var y = group[0].Y
	for a, player := range group {
		if player.Play == Weapon {
			cardIndex, hasCard := playerHasCard(player, Weapon)
			if hasCard {
				attackValue += weaponStrength
				group[a].Cards[cardIndex] = None
			} else {
				attackValue += r.Intn(playerMaxAttack-1) + playerMinAttack
			}
		} else {
			attackValue += r.Intn(playerMaxAttack-1) + playerMinAttack
		}
		group[a].Play = Dice
	}
	if attackValue < gMap[x][y].Zombies {
		for a := range group {
			group[a].Alive = false
		}
		gMap[x][y].Zombies += len(group)
	} else {
		gMap[x][y].Zombies = 0
	}
}

// TODO: citiesList can be value instead of reference
// TODO: decide if spread is 4 or 8 directions
func spread(gMap *[mapWidth][mapHeight]*Tile, cities *[]IntTuple) {
	for _, city := range *cities {
		if gMap[city.X][city.Y].Zombies < zombieCutoff {
			gMap[city.X][city.Y].Zombies++
			continue
		}
		//North
		if city.Y < mapHeight-1 && gMap[city.X][city.Y+1].Zombies < zombieCutoff {
			gMap[city.X][city.Y+1].Zombies++
		}
		//East
		if city.X < mapWidth-1 && gMap[city.X+1][city.Y].Zombies < zombieCutoff {
			gMap[city.X+1][city.Y].Zombies++
		}
		//South
		if city.Y > 0 && gMap[city.X][city.Y-1].Zombies < zombieCutoff {
			gMap[city.X][city.Y-1].Zombies++
		}
		//West
		if city.X > 0 && gMap[city.X-1][city.Y].Zombies < zombieCutoff {
			gMap[city.X-1][city.Y].Zombies++
		}
	}
}

func updatePlayerPositions(pList *[]*Player, gMap *[mapWidth][mapHeight]*Tile) {
	//Clear all player info in tiles
	for _, tileRow := range *gMap {
		for _, tile := range tileRow {
			(*tile).Players = []Player{}
		}
	}
	//Reposition all players
	for _, player := range *pList {
		var tilePList = gMap[player.X][player.Y].Players
		gMap[player.X][player.Y].Players = append(tilePList, *player)
	}
}

// TODO: Unify order of attributes across functions
func tick(gMap *[mapWidth][mapHeight]*Tile, cities *[]IntTuple, pList *[]*Player) {
	move(pList)
	resources(pList, gMap)
	handleCombat(gMap, pList)
	spread(gMap, cities)
	consume(pList, gMap)
	limitCards(pList)
	updatePlayerPositions(pList, gMap)
}

func playerHasCard(player *Player, card Card) (int, bool) {
	for a, c := range player.Cards {
		if c == card {
			return a, true
		}
	}
	return -1, false
}

func randomizeBot(bots []*Player) {
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
// TODO: Make sure ID has no /
func addPlayer(players *[]*Player, playerName string) string {
	var rX = r.Intn(mapWidth - 1)
	var rY = r.Intn(mapHeight - 1)
	var nowString = strconv.Itoa(int(time.Now().UnixNano() << 2))
	var playerID = ""
	for i := 0; i < len(nowString); i++ {
		playerID += string(nowString[i] ^ idSalt[i])
	}
	playerID = base64.StdEncoding.EncodeToString([]byte(playerID))
	var player = Player{
		ID:        playerID,
		Name:      playerName,
		X:         rX,
		Y:         rY,
		Direction: Stay,
		Play:      None,
		Consume:   None,
		Discard:   None,
		Cards:     [5]Card{Food, Wood, Wood, None, None},
		Alive:     true,
		IsBot:     false,
	}
	*players = append(*players, &player)
	return playerID
}

func addBot(players *[]*Player, bots *[]*Player, bID int) {
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

func havePlayersWon(players *[]*Player) bool {
	for _, player := range *players {
		if player.hasWinCondition() {
			return true
		}
	}
	return false
}

func main() {
	if len(os.Args) == 2 {
		idSalt = os.Args[1]
		fmt.Println(idSalt)
	}
	var gameMap [mapWidth][mapHeight]*Tile
	var cityList []IntTuple
	var playerList []*Player
	var botList []*Player
	var botID = 0
	var turnTimer = int8(turnLength)
	hasWon = false
	go setupAPI(&playerList, &gameMap, &turnTimer, &hasWon)
	for {
		r = rand.New(rand.NewSource(time.Now().Unix()))
		initMap(*r, &gameMap)
		cityList = createCityList(&gameMap)
		var remainingTurns = maxTurns
		for ; remainingTurns > 0; remainingTurns-- {
			fmt.Println("Remaining turns: ", remainingTurns)
			for turnTimer = int8(turnLength); turnTimer >= 0; turnTimer-- {
				if turnTimer == 0 {
					randomizeBot(botList)
					tick(&gameMap, &cityList, &playerList)
					hasWon = havePlayersWon(&playerList)
					restockBots(&playerList, &botList, &botID)
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
		fmt.Println("Restarting game.")
		time.Sleep(120 * time.Second)
	}
}
