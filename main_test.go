package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovement(t *testing.T) {
	t.Run("player returns to original position after all directions", func(t *testing.T) {
		// Arrange
		ts := setupTestSuite(t)
		originalX, originalY := 1, 1
		playerID := ts.createPlayerAt(originalX, originalY)
		player := ts.getPlayer(playerID)

		// Act - move in all directions (should return to original position)
		for i := 0; i < 5; i++ {
			player.Direction = Directions[i]
			ts.playerMap.move()
		}

		// Assert
		assert.Equal(t, originalX, player.CurrentTile.XPos, "Player X position should match original")
		assert.Equal(t, originalY, player.CurrentTile.YPos, "Player Y position should match original")
	})
}

func TestCombat(t *testing.T) {
	tests := []struct {
		name           string
		zombieCount    int
		playerCount    int
		hasWeapon      bool
		expectedAlive  bool
		description    string
	}{
		{
			name:          "player loses against overwhelming zombies",
			zombieCount:   100,
			playerCount:   1,
			hasWeapon:     false,
			expectedAlive: false,
			description:   "Single player should not survive against 100 zombies",
		},
		{
			name:          "multiple players win against single zombie",
			zombieCount:   1,
			playerCount:   2,
			hasWeapon:     false,
			expectedAlive: true,
			description:   "Two players should survive against one zombie",
		},
		{
			name:          "player with weapon wins against multiple zombies",
			zombieCount:   5,
			playerCount:   1,
			hasWeapon:     true,
			expectedAlive: true,
			description:   "Player with weapon should survive against 5 zombies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ts := setupTestSuite(t)
			xPos, yPos := 1, 1
			playerID := ts.createPlayerAt(xPos, yPos)
			player := ts.getPlayer(playerID)

			// Add additional players if needed
			for i := 1; i < tt.playerCount; i++ {
				ts.playerMap.addPlayer("TestPlayer"+string(rune(i+1)), ts.gameMap.getTileFromPos(xPos, yPos))
			}

			// Set up weapon if needed
			if tt.hasWeapon {
				player.Cards[0] = Weapon
				player.Play = Weapon
			}

			// Add zombies
			ts.gameMap.addZombiesToTile(xPos, yPos, tt.zombieCount)

			// Act
			ts.gameMap.handleCombat()

			// Assert
			assert.Equal(t, tt.expectedAlive, player.Alive, tt.description)
		})
	}
}





func TestPlayerConsume(t *testing.T) {
	tests := []struct {
		name         string
		initialCards [5]Card
		consumeCard  Card
		expectedCards [5]Card
		description  string
	}{
		{
			name:         "consume food card",
			initialCards: [5]Card{Weapon, Food, Wood, Wood, Wood},
			consumeCard:  Food,
			expectedCards: [5]Card{Weapon, None, Wood, Wood, Wood},
			description:  "Food card should be consumed and replaced with None",
		},
		{
			name:         "consume wood card",
			initialCards: [5]Card{Weapon, Food, Wood, Wood, None},
			consumeCard:  Wood,
			expectedCards: [5]Card{Weapon, Food, None, Wood, None},
			description:  "First wood card should be consumed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ts := setupTestSuite(t)
			playerID := ts.createPlayerAt(1, 1)
			player := ts.getPlayer(playerID)
			player.Cards = tt.initialCards
			player.Consume = tt.consumeCard

			// Act
			ts.playerMap.playersConsume()

			// Assert
			assert.Equal(t, tt.expectedCards, player.Cards, tt.description)
		})
	}
}

func TestResourceDistribution(t *testing.T) {
	tests := []struct {
		name            string
		terrain         Terrain
		initialCards    [5]Card
		expectedNewCard Card
		expectedCount   int
		description     string
	}{
		{
			name:            "forest gives wood",
			terrain:         Forest,
			initialCards:    [5]Card{Weapon, None, None, Wood, Wood},
			expectedNewCard: Wood,
			expectedCount:   2, // Forest gives 2 wood
			description:     "Forest terrain should provide wood cards",
		},
		{
			name:            "laboratory gives research",
			terrain:         Laboratory,
			initialCards:    [5]Card{None, None, None, None, None},
			expectedNewCard: Research,
			expectedCount:   1,
			description:     "Laboratory terrain should provide research cards",
		},
		{
			name:            "city gives weapon",
			terrain:         City,
			initialCards:    [5]Card{None, None, None, None, None},
			expectedNewCard: Weapon,
			expectedCount:   1,
			description:     "City terrain should provide weapon cards",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			ts := setupTestSuite(t)
			xPos, yPos := 1, 1
			playerID := ts.createPlayerAt(xPos, yPos)
			player := ts.getPlayer(playerID)
			player.Cards = tt.initialCards
			ts.gameMap.getTileFromPos(xPos, yPos).Terrain = tt.terrain

			// Count initial cards of expected type
			initialCount := ts.countCards(player.Cards, tt.expectedNewCard)

			// Act
			ts.gameMap.resources()

			// Assert
			finalCount := ts.countCards(player.Cards, tt.expectedNewCard)
			expectedFinalCount := initialCount + tt.expectedCount
			assert.Equal(t, expectedFinalCount, finalCount, tt.description)
		})
	}
}

func TestWinConditions(t *testing.T) {
	t.Run("player cannot win at same laboratory where research was acquired", func(t *testing.T) {
		// Arrange
		ts := setupTestSuite(t)
		xPos, yPos := 1, 1
		playerID := ts.createPlayerAt(xPos, yPos)
		player := ts.getPlayer(playerID)

		// Set up laboratory and clear initial cards
		ts.gameMap.getTileFromPos(xPos, yPos).Terrain = Laboratory
		player.Cards = [5]Card{None, None, None, None, None}

		// Act - acquire research cards at the first laboratory
		for i := 0; i < victoryNumber; i++ {
			ts.gameMap.resources()
		}

		// Assert - player should not win at same laboratory
		assert.False(t, ts.playerMap.havePlayersWon(), "Player should not win at the same laboratory where research was acquired")
	})

	t.Run("player can win at different laboratory with enough research cards", func(t *testing.T) {
		// Arrange
		ts := setupTestSuite(t)
		xPos, yPos := 1, 1
		playerID := ts.createPlayerAt(xPos, yPos)
		player := ts.getPlayer(playerID)

		// Set up first laboratory and acquire research cards
		ts.gameMap.getTileFromPos(xPos, yPos).Terrain = Laboratory
		player.Cards = [5]Card{None, None, None, None, None}
		for i := 0; i < victoryNumber; i++ {
			ts.gameMap.resources()
		}

		// Move player to a different laboratory
		newX, newY := 2, 2
		ts.movePlayerTo(player, newX, newY)
		ts.gameMap.getTileFromPos(newX, newY).Terrain = Laboratory

		// Assert - player should be able to win at different laboratory
		assert.True(t, ts.playerMap.havePlayersWon(), "Player should win when at a different laboratory with enough research cards")
	})

	t.Run("player cannot win without enough research cards", func(t *testing.T) {
		// Arrange
		ts := setupTestSuite(t)
		playerID := ts.createPlayerAt(1, 1)
		player := ts.getPlayer(playerID)
		ts.gameMap.getTileFromPos(1, 1).Terrain = Laboratory
		player.Cards = [5]Card{Research, None, None, None, None} // Only 1 research card

		// Act & Assert
		assert.False(t, ts.playerMap.havePlayersWon(), "Player should not win with insufficient research cards")
	})
}

// TestSuite provides isolated test environment
type TestSuite struct {
	t         *testing.T
	gameMap   *gameMap
	playerMap *playerMap
	gameState *gameState
}

// setupTestSuite creates a new isolated test environment
func setupTestSuite(t *testing.T) *TestSuite {
	// Initialize global variables that the existing code depends on
	gMap = NewGameMap()
	pMap = NewPlayerMap()
	gState = NewGameState()
	
	return &TestSuite{
		t:         t,
		gameMap:   &gMap,
		playerMap: &pMap,
		gameState: &gState,
	}
}

// createPlayerAt creates a test player at the specified position
func (ts *TestSuite) createPlayerAt(x, y int) string {
	// Ensure the game map is properly initialized
	if ts.gameMap.gMap[0][0] == nil {
		ts.gameMap.init()
	}
	return ts.playerMap.addPlayer("TestPlayer1", ts.gameMap.getTileFromPos(x, y))
}

// getPlayer returns a player by ID
func (ts *TestSuite) getPlayer(playerID string) *Player {
	return ts.playerMap.getPlayerPtr(playerID)
}

// movePlayerTo moves a player to a new position
func (ts *TestSuite) movePlayerTo(player *Player, x, y int) {
	oldTile := player.CurrentTile
	newTile := ts.gameMap.getTileFromPos(x, y)
	oldTile.removePlayer(player)
	newTile.addPlayer(player)
	player.CurrentTile = newTile
}

// countCards counts how many cards of a specific type are in the hand
func (ts *TestSuite) countCards(cards [5]Card, cardType Card) int {
	count := 0
	for _, card := range cards {
		if card == cardType {
			count++
		}
	}
	return count
}

// Legacy setup functions for backward compatibility
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
