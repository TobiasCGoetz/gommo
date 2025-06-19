package main

type gamestate struct {
	turnTimer      int8
	havePlayersWon bool
}

func newGameState() gamestate {
	return gamestate{0, false}
}

func (gs gamestate) haveWon() bool {
	return gs.havePlayersWon
}

func (gs *gamestate) win() {
	gs.havePlayersWon = true
}

func (gs *gamestate) timeTick() {
	gs.turnTimer--
}

func (gs *gamestate) timeTickWithCheck() bool {
	gs.turnTimer--
	if gs.turnTimer <= 0 {
		gs.turnTimer = turnLength
		return true
	}
	return false
}

func (gs gamestate) isTurnOver() bool {
	if gs.turnTimer <= 0 {
		return true
	}
	return false
}

func (gs *gamestate) resetTime() {
	gs.turnTimer = turnLength
}
