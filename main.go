package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var playerMap = make(map[string]*Player)
var botList []*Player
var r *rand.Rand
var gMap gameMap

func rollDice() int {
	return rand.Intn(playerMaxAttack) + playerMinAttack
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
			var cardPos, hasCard = hasCardWhere(player.Cards[:], player.Discard)
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

// TODO: Unify order of attributes across functions
func tick(playerMap *map[string]*Player) {
	fmt.Println("Next tick is happening...")
	fmt.Println("Moving players...")
	move(playerMap)
	fmt.Println("Distributing ressources...")
	gMap.resources()
	fmt.Println("Combat is upon us...")
	gMap.handleCombat()
	fmt.Println("The infection is spreading...")
	gMap.spread()
	fmt.Println("Players feeding themselves...")
	gMap.consume(playerMap)
	fmt.Println("Limiting player inventory")
	limitCards(playerMap)
}

func randomizeBots(bots []*Player) {
	fmt.Println("Randomizing bot turns...")
	for _, bot := range bots {
		//Randomize movement
		bot.Direction = Directions[r.Intn(len(Directions))]
		//Randomize card played
		bot.Play = Dice
		//Randomize consume

		if _, foodFound := hasCardWhere(bot.Cards[:], Food); foodFound {
			bot.Consume = Food
		} else if _, woodFound := hasCardWhere(bot.Cards[:], Wood); woodFound {
			bot.Consume = Wood
		} else {
			bot.Consume = None
		}

		//Randomize discard
		_, found := hasCardWhere(bot.Cards[:], None)
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

func main() {
	if len(os.Args) == 2 {
		idSalt = os.Args[1]
		fmt.Println(idSalt)
	}
	gMap = NewGameMap()
	var botID = 0
	var turnTimer = int8(turnLength)
	hasWon = false
	r = rand.New(rand.NewSource(time.Now().Unix()))
	registry := newHandlerRegistry()
	//registry.AddHandler(CreateUserEvent{}.Type(), CreateUserHandler)
	//baseEvent := BaseEvent{"playerId", time.Now(), BaseEvent{}.Type(), false}
	//createUserEvent := CreateUserEvent{baseEvent, "username"}
	//fmt.Println(baseEvent.Type(), createUserEvent.Type())
	//registry.Handle(createUserEvent)

	go setupAPI(registry, &turnTimer, &hasWon)

	var remainingTurns = maxTurns
	for ; remainingTurns > 0; remainingTurns-- {
		fmt.Println("Remaining turns: ", remainingTurns)
		//printMap()
		for turnTimer = int8(turnLength); turnTimer >= 0; turnTimer-- {
			if turnTimer == 0 {
				randomizeBots(botList)
				tick(&playerMap)
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
