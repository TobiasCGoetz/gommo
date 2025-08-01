package main

import (
	"testing"
)

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

func TestCombatLose(t *testing.T) {
	xPos, yPos := 1, 1
	playerId := setupTest(xPos, yPos)
	playerPtr := pMap.getPlayerPtr(playerId)
	gMap.addZombiesToTile(xPos, yPos, 100)
	gMap.handleCombat()
	if playerPtr.Alive {
		t.Errorf("Player survived impossible combat encounter")
	}
}

func TestCombatWin(t *testing.T) {
	xPos, yPos := 1, 1
	playerId := setupTest(xPos, yPos)
	playerPtr := pMap.getPlayerPtr(playerId)
	//Add another player for guaranteed win against single Zombie
	pMap.addPlayer("TestPlayer2", gMap.getTileFromPos(xPos, yPos))
	gMap.getTileFromPos(xPos, yPos).addZombies(1)
	t.Logf(
		"%d/%d has %d Zombies and %d players",
		xPos,
		yPos,
		gMap.getTileFromPos(xPos, yPos).Zombies,
		len(gMap.getTileFromPos(xPos, yPos).playerPtrs))
	gMap.handleCombat()
	if !playerPtr.Alive {
		t.Errorf("Player died in safe combat encounter")
	}
}

func TestCombatItemWin(t *testing.T) {
	xPos, yPos := 1, 1
	playerId := setupTest(xPos, yPos)
	playerPtr := pMap.getPlayerPtr(playerId)
	playerPtr.Cards[0] = Weapon
	playerPtr.Play = Weapon
	gMap.getTileFromPos(xPos, yPos).addZombies(5)
	if !playerPtr.Alive {
		t.Errorf("Player died regardless of adequate weapon use")
	}
}

func TestConsume(t *testing.T)   {}
func TestResources(t *testing.T) {}
func TestWin(t *testing.T)       {}

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
