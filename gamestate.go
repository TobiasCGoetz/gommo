package main

type gameState struct {
	turnTimer      int
	remainingTurns int
	havePlayersWon bool
}

func NewGameState() gameState {
	return gameState{turnLength, maxTurns, false}
}

func (gs gameState) haveWon() bool {
	return gs.havePlayersWon
}

func (gs *gameState) win() {
	gs.havePlayersWon = true
}

func (gs *gameState) timerDown() {
	gs.turnTimer--
}

func (gs gameState) isTurnOver() bool {
	if gs.turnTimer <= 0 {
		return true
	}
	return false
}

func (gs *gameState) resetTime() {
	gs.turnTimer = turnLength
	gs.remainingTurns--
}

func (gs gameState) getRemainingTurns() int {
	return gs.remainingTurns
}

func (gs gameState) isGameOver() bool {
	if gs.remainingTurns < 0 || gs.haveWon() {
		return true
	}
	return false
}
