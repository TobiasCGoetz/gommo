package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

var pMap playerMap
var r *rand.Rand
var gMap gameMap
var gState gameState

// eventLogger is the global event logger instance
var eventLogger EventLogger

func rollDice(playerID string) int {
	result := rand.Intn(gameConfig.Combat.PlayerMaxAttack) + gameConfig.Combat.PlayerMinAttack

	// Log the dice roll event
	eventLogger.LogEvent(EventDiceRoll, playerID, map[string]interface{}{
		"result": result,
		"min":    gameConfig.Combat.PlayerMinAttack,
		"max":    gameConfig.Combat.PlayerMaxAttack,
	})

	return result
}

func tick() {
	fmt.Println("# Tick")
	fmt.Println("Moving players...")
	pMap.move()
	fmt.Println("Distributing ressources...")
	gMap.resources()
	fmt.Println("Combat is upon us...")
	gMap.handleCombat()
	fmt.Println("The infection is spreading...")
	gMap.spread()
	fmt.Println("Players feeding themselves...")
	pMap.playersConsume()
	fmt.Println("Limiting player inventory")
	pMap.limitCards()
}

func getPlayerOrNil(id string) *Player {
	return pMap.Players[id] //TODO: Improve
}

func main() {
	if len(os.Args) == 2 {
		gameConfig.Server.IDSalt = os.Args[1]
		fmt.Println(gameConfig.Server.IDSalt)
	}

	r = rand.New(rand.NewSource(time.Now().Unix()))
	gMap = NewGameMap()
	gState = NewGameState()
	pMap = NewPlayerMap()
	eventLogger = NewEventLogger() // Initialize the global event logger

	go setupAPI()

	fmt.Println("Remaining turns: ", gState.getRemainingTurns())
	for !gState.haveWon() {
		if !gState.isTurnOver() {
			time.Sleep(time.Second)
			gState.timerDown()
		} else {
			gState.resetTime()
			tick()
			fmt.Println("Remaining turns: ", gState.getRemainingTurns())
			if pMap.havePlayersWon() {
				fmt.Println("Game over due to win")
				gState.win()
			}
		}
	}
}
