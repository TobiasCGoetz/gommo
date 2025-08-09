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
	t.Logf("Tile has %d players", len(gMap.getTileFromPos(xPos, yPos).playerPtrs))
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

func TestConsume(t *testing.T) {
	xPos, yPos := 1, 1
	playerId := setupTest(xPos, yPos)
	playerPtr := pMap.getPlayerPtr(playerId)
	playerPtr.Cards = [5]Card{Weapon, Food, Wood, Wood, Wood}
	playerPtr.Consume = Food
	pMap.playersConsume()
	if playerPtr.Cards != [5]Card{Weapon, None, Wood, Wood, Wood} {
		t.Errorf("Cards in hand incorrectly consumed. Wrong result is:")
		for _, value := range playerPtr.Cards {
			t.Errorf("- %s", value.toString())
		}
	}
}

func TestResources(t *testing.T) {
	xPos, yPos := 1, 1
	playerId := setupTest(xPos, yPos)
	playerPtr := pMap.getPlayerPtr(playerId)
	playerPtr.Cards = [5]Card{Weapon, None, None, Wood, Wood}
	gMap.getTileFromPos(xPos, yPos).Terrain = Forest
	gMap.resources()
	for _, value := range playerPtr.Cards {
		t.Logf("- %s", value.toString())
	}
}

func TestWin(t *testing.T) {
	//Test if a player can win if he arrives at a laboratory while holding a research card
	//It's important, that he can't win at the laboratory where he got the research card
	//He has to bring it to a different laboratory than the card originated from
	xPos, yPos := 1, 1
	playerId := setupTest(xPos, yPos)
	playerPtr := pMap.getPlayerPtr(playerId)

	// Set up a laboratory tile and acquire research cards through the proper system
	gMap.getTileFromPos(xPos, yPos).Terrain = Laboratory
	playerPtr.Cards = [5]Card{None, None, None, None, None} // Clear initial cards
	
	// Acquire research cards at the first laboratory (this should be tracked)
	gMap.resources() // This gives research cards and tracks their acquisition location
	gMap.resources() // Get more research cards
	gMap.resources() // Get more research cards
	gMap.resources() // Get more research cards
	gMap.resources() // Get enough research cards for victory

	// Player should not win at the same laboratory where they got the research
	if pMap.havePlayersWon() {
		t.Errorf("Player won at the same laboratory where they got research cards")
	}

	// Move player to a different laboratory
	newX, newY := 2, 2
	newTile := gMap.getTileFromPos(newX, newY)
	newTile.Terrain = Laboratory
	playerPtr.CurrentTile.removePlayer(playerPtr)
	newTile.addPlayer(playerPtr)
	playerPtr.CurrentTile = newTile

	// Now player should be able to win at the different laboratory
	if !pMap.havePlayersWon() {
		t.Errorf("Player should win when at a different laboratory with enough research cards")
	}

}

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
