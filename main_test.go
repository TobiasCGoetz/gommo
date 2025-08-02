package main

import (
	"testing"
)

//t.Errorf("Got %d, but wanted 1", got)

func TestMovement(t *testing.T) {
	xPos, yPos := 1, 1
	playerId := setupTest(xPos, yPos)
	playerPtr := pMap.getPlayerPtr(playerId)
	for i := 0; i < 5; i++ {
		playerPtr.Direction = Directions[i]
		pMap.move()
	}
	if playerPtr.CurrentTile.XPos != xPos ||
		playerPtr.CurrentTile.YPos != yPos {
		t.Errorf(
			"Got %d/%d position, but expected %d/%d",
			playerPtr.CurrentTile.XPos,
			playerPtr.CurrentTile.YPos,
			xPos,
			yPos)
	}
}

func TestCombatLose(t *testing.T) {}
func TestCombatWin(t *testing.T)  {}
func TestConsume(t *testing.T)    {}
func TestResources(t *testing.T)  {}
func TestWin(t *testing.T)        {}

func setupTest(xPos int, yPos int) string {
	setupTestState()
	setupTestMap()
	return setupTestPlayer(xPos, yPos)
}

func setupTestState() {
	gState = NewGameState()
}

func setupTestPlayer(xPos int, yPos int) string {
	pMap = NewPlayerMap()
	return pMap.addPlayer("TestPlayer1", gMap.getTileFromPos(xPos, yPos))
}

func setupTestMap() {
	gMap = NewGameMap()
}
