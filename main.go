package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var pMap playerMap
var botList []*Player
var r *rand.Rand
var gMap gameMap

func rollDice() int {
	return rand.Intn(playerMaxAttack) + playerMinAttack
}

// TODO: Unify order of attributes across functions
func tick() {
	fmt.Println("Next tick is happening...")
	fmt.Println("Moving players...")
	pMap.move()
	fmt.Println("Distributing ressources...")
	gMap.resources()
	fmt.Println("Combat is upon us...")
	gMap.handleCombat()
	fmt.Println("The infection is spreading...")
	gMap.spread()
	fmt.Println("Players feeding themselves...")
	gMap.consume(&pMap.Players)
	fmt.Println("Limiting player inventory")
	pMap.limitCards()
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

func addBot(playerMap *map[string]*Player, bots *[]*Player, bID int) {
	var bot = Player{
		ID:          strconv.Itoa(bID),
		CurrentTile: nil,
		Direction:   Stay,
		Play:        None,
		Consume:     None,
		Discard:     None,
		Cards:       [5]Card{Food, Wood, Wood, None, None},
		Alive:       true,
		IsBot:       true,
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

func getPlayerOrNil(id string) *Player {
	return pMap.Players[id] //TODO: Improve
}

func main() {
	if len(os.Args) == 2 {
		idSalt = os.Args[1]
		fmt.Println(idSalt)
	}
	//var botID = 0
	var turnTimer = int8(turnLength)
	hasWon = false
	r = rand.New(rand.NewSource(time.Now().Unix()))
	gMap = NewGameMap()
	//registry.AddHandler(CreateUserEvent{}.Type(), CreateUserHandler)
	//baseEvent := BaseEvent{"playerId", time.Now(), BaseEvent{}.Type(), false}
	//createUserEvent := CreateUserEvent{baseEvent, "username"}
	//fmt.Println(baseEvent.Type(), createUserEvent.Type())
	//registry.Handle(createUserEvent)

	go setupAPI()

	var remainingTurns = maxTurns
	for ; remainingTurns > 0; remainingTurns-- {
		fmt.Println("Remaining turns: ", remainingTurns)
		//printMap()
		for turnTimer = int8(turnLength); turnTimer >= 0; turnTimer-- {
			if turnTimer == 0 {
				randomizeBots(botList)
				tick()
				hasWon = pMap.havePlayersWon()
				//restockBots(&playerMap, &botList, &botID)
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
