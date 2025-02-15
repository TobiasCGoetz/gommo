package main

import (
	"fmt"
	"testing"
	"math/rand"
	"github.com/stretchr/testify/assert"
)

// --- terrain_test.go ---
func TestTerrain_toString(t *testing.T) {
	tests := []struct {
		terrain Terrain
		want    string
	}{
		{terrain: Forest, want: "Forest"},
		{terrain: Farm, want: "Farm"},
		{terrain: City, want: "City"},
		{terrain: Laboratory, want: "Laboratory"},
		{terrain: Edge, want: "Edge"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.terrain.toString(); got != tt.want {
				t.Errorf("Terrain.toString() for %v = %v, want %v", tt.terrain, got, tt.want)
			}
		})
	}
}

func TestTerrain_toChar(t *testing.T) {
	tests := []struct {
		terrain Terrain
		want    string
	}{
		{terrain: Forest, want: "🌲"},
		{terrain: Farm, want: "🌱"},
		{terrain: City, want: "🏠"},
		{terrain: Laboratory, want: "🧬"},
		{terrain: Edge, want: "⌘"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.terrain.toChar(); got != tt.want {
				t.Errorf("Terrain.toChar() for %v = %v, want %v", tt.terrain, got, tt.want)
			}
		})
	}
}

func TestTerrain_isCity(t *testing.T) {
	tests := []struct {
		terrain Terrain
		want    bool
	}{
		{terrain: Forest, want: false},
		{terrain: Farm, want: false},
		{terrain: City, want: true},
		{terrain: Laboratory, want: false},
		{terrain: Edge, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.terrain.toString(), func(t *testing.T) {
			if got := tt.terrain.isCity(); got != tt.want {
				t.Errorf("Terrain.isCity() for %v = %v, want %v", tt.terrain, got, tt.want)
			}
		})
	}
}

func TestTerrain_offersResource(t *testing.T) {
	tests := []struct {
		terrain     Terrain
		wantCard    Card
		wantAmount  int
		testName    string
		expectedNil bool
	}{
		{terrain: City, wantCard: Weapon, wantAmount: 1, testName: "CityOffersWeapon"},
		{terrain: Forest, wantCard: Wood, wantAmount: 2, testName: "ForestOffersWood"},
		{terrain: Farm, wantCard: Food, wantAmount: 1, testName: "FarmOffersFood"},
		{terrain: Laboratory, wantCard: Research, wantAmount: 1, testName: "LaboratoryOffersResearch"},
		{terrain: Edge, wantCard: None, wantAmount: 0, testName: "EdgeOffersNone", expectedNil: true}, // Edge offers no resource
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			card, amount := tt.terrain.offersResource()
			if card != tt.wantCard {
				t.Errorf("Terrain.offersResource() for %v returned card = %v, want %v", tt.terrain, card, tt.wantCard)
			}
			if amount != tt.wantAmount {
				t.Errorf("Terrain.offersResource() for %v returned amount = %v, want %v", tt.terrain, amount, tt.wantAmount)
			}
		})
	}
}

// --- card_test.go ---
func TestCard_toString(t *testing.T) {
	tests := []struct {
		card Card
		want string
	}{
		{card: Food, want: "Food"},
		{card: Wood, want: "Wood"},
		{card: Weapon, want: "Weapon"},
		{card: Dice, want: "Dice"},
		{card: Research, want: "Research"},
		{card: None, want: "None"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.card.toString(); got != tt.want {
				t.Errorf("Card.toString() for %v = %v, want %v", tt.card, got, tt.want)
			}
		})
	}
}

func TestCard_MarshalJSON(t *testing.T) {
	tests := []struct {
		card Card
		want string
	}{
		{card: Food, want: "\"Food\""},
		{card: Wood, want: "\"Wood\""},
		{card: Weapon, want: "\"Weapon\""},
		{card: Dice, want: "\"Dice\""},
		{card: Research, want: "\"Research\""},
		{card: None, want: "\"None\""},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			jsonBytes, _ := tt.card.MarshalJSON()
			jsonString := string(jsonBytes)
			if jsonString != tt.want {
				t.Errorf("Card.MarshalJSON() for %v = %v, want %v", tt.card, jsonString, tt.want)
			}
		})
	}
}

func TestCardsMap(t *testing.T) {
	expectedCards := map[string]Card{
		"food":     Food,
		"wood":     Wood,
		"weapon":   Weapon,
		"dice":     Dice,
		"research": Research,
		"none":     None,
	}
	assert.Equal(t, expectedCards, cards, "Cards map should be initialized correctly")
}

// --- direction_test.go ---
func TestDirection_toString(t *testing.T) {
	tests := []struct {
		direction Direction
		want      string
	}{
		{direction: North, want: "North"},
		{direction: East, want: "East"},
		{direction: South, want: "South"},
		{direction: West, want: "West"},
		{direction: Stay, want: "Stay"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.direction.toString(); got != tt.want {
				t.Errorf("Direction.toString() for %v = %v, want %v", tt.direction, got, tt.want)
			}
		})
	}
}

func TestDirection_MarshalJSON(t *testing.T) {
	tests := []struct {
		direction Direction
		want      string
	}{
		{direction: North, want: "\"North\""},
		{direction: East, want: "\"East\""},
		{direction: South, want: "\"South\""},
		{direction: West, want: "\"West\""},
		{direction: Stay, want: "\"Stay\""},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			jsonBytes, _ := tt.direction.MarshalJSON()
			jsonString := string(jsonBytes)
			if jsonString != tt.want {
				t.Errorf("Direction.MarshalJSON() for %v = %v, want %v", tt.direction, jsonString, tt.want)
			}
		})
	}
}

func TestDirectionsMap(t *testing.T) {
	expectedDirections := map[string]Direction{
		"north": North,
		"east":  East,
		"south": South,
		"west":  West,
		"stay":  Stay,
	}
	assert.Equal(t, expectedDirections, directions, "Directions map should be initialized correctly")
}

// --- tile_test.go ---

func TestTile_resolveCombat_PlayerWinsWeapon(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Alive: true, Cards: [5]Card{Weapon, None, None, None, None}, Play: Weapon}}
	tile := &Tile{Terrain: City, Zombies: 5, playerIds: []string{playerID}}

	originalRollDice := rollDice
	rollDice = func() int { return 1 } // Mock rollDice to always return a low value, weapon should override
	defer func() { rollDice = originalRollDice }() // Restore original rollDice

	tile.resolveCombat()

	assert.Equal(t, 0, tile.Zombies, "Zombies should be 0 after player wins with weapon")
	assert.True(t, playerMap[playerID].Alive, "Player should still be alive")
	weaponCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], Weapon)
	assert.Equal(t, -1, weaponCardIndex, "Weapon card should be removed after use (-1 means not found, None is overwritten)") // -1 because None overwrites card
	noneCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], None)
	assert.NotEqual(t, -1, noneCardIndex, "Should have a None card after weapon used")
}

func TestTile_resolveCombat_PlayerWinsDice(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Alive: true, Cards: [5]Card{Dice, None, None, None, None}, Play: Dice}}
	tile := &Tile{Terrain: City, Zombies: 5, playerIds: []string{playerID}}

	originalRollDice := rollDice
	rollDice = func() int { return 10 } // Mock rollDice to always return a high value
	defer func() { rollDice = originalRollDice }()

	tile.resolveCombat()

	assert.Equal(t, 0, tile.Zombies, "Zombies should be 0 after player wins with dice")
	assert.True(t, playerMap[playerID].Alive, "Player should still be alive")
}

func TestTile_resolveCombat_PlayerLoses(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Alive: true, Cards: [5]Card{None, None, None, None, None}, Play: Dice}} // No weapon
	tile := &Tile{Terrain: City, Zombies: 10, playerIds: []string{playerID}}

	originalRollDice := rollDice
	rollDice = func() int { return 1 } // Mock rollDice to always return a low value
	defer func() { rollDice = originalRollDice }()

	tile.resolveCombat()

	assert.Greater(t, tile.Zombies, 0, "Zombies should remain or increase after player loses") // Zombies might increase depending on logic
	assert.False(t, playerMap[playerID].Alive, "Player should be dead")
	assert.Contains(t, tile.playerIds, playerID, "Player ID should still be on the tile") // Player ID is still there after death
}

func TestTile_resolveCombat_NoPlayers(t *testing.T) {
	tile := &Tile{Terrain: City, Zombies: 5, playerIds: []string{}} // No players
	initialZombies := tile.Zombies
	tile.resolveCombat()
	assert.Equal(t, initialZombies, tile.Zombies, "Zombie count should not change if no players")
}

func TestTile_giveResources_Farm(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Alive: true, Cards: [5]Card{None, None, None, None, None}}}
	tile := &Tile{Terrain: Farm, Zombies: 0, playerIds: []string{playerID}}

	tile.giveResources()

	foodCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], Food)
	assert.NotEqual(t, -1, foodCardIndex, "Player should receive Food on Farm")
}

func TestTile_giveResources_Forest(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Alive: true, Cards: [5]Card{None, None, None, None, None}}}
	tile := &Tile{Terrain: Forest, Zombies: 0, playerIds: []string{playerID}}

	tile.giveResources()

	woodCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], Wood)
	woodCardIndex2, _ := hasCardWhere(playerMap[playerID].Cards[woodCardIndex+1:], Wood) // Check for second wood card
	assert.NotEqual(t, -1, woodCardIndex, "Player should receive Wood on Forest")
	assert.NotEqual(t, -1, woodCardIndex2, "Player should receive two Woods on Forest") // Check for second wood
}

func TestTile_giveResources_City_NoSpace(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Alive: true, Cards: [5]Card{Food, Food, Food, Food, Food}}} // Full hand
	tile := &Tile{Terrain: City, Zombies: 0, playerIds: []string{playerID}}

	tile.giveResources()

	weaponCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], Weapon)
	assert.Equal(t, -1, weaponCardIndex, "Player should NOT receive Weapon on City if hand full")
}

func TestTile_isSpreader_City(t *testing.T) {
	tile := &Tile{Terrain: City, Zombies: 0}
	assert.True(t, tile.isSpreader(), "City tile should be a spreader")
}

func TestTile_isSpreader_ZombieCutoff(t *testing.T) {
	tile := &Tile{Terrain: Farm, Zombies: zombieCutoff, playerIds: []string{}}
	assert.True(t, tile.isSpreader(), "Tile with zombies >= cutoff should be a spreader")
	tileUnder := &Tile{Terrain: Farm, Zombies: zombieCutoff - 1, playerIds: []string{}}
	assert.False(t, tileUnder.isSpreader(), "Tile with zombies < cutoff should not be a spreader")
}

func TestTile_spreadTo(t *testing.T) {
	tile := &Tile{Terrain: Farm, Zombies: 5}
	tile.spreadTo()
	assert.Equal(t, 6, tile.Zombies, "Zombies should increase by 1 after spreadTo if below cutoff")
	tile2 := &Tile{Terrain: Farm, Zombies: zombieCutoff}
	tile2.spreadTo()
	assert.Equal(t, zombieCutoff, tile2.Zombies, "Zombies should NOT increase if at cutoff")
}

func TestTile_spreadToUnbound(t *testing.T) {
	tile := &Tile{Terrain: Farm, Zombies: 5}
	tile.spreadToUnbound()
	assert.Equal(t, 6, tile.Zombies, "Zombies should always increase by 1 after spreadToUnbound")
}

func TestTile_addRemovePlayer(t *testing.T) {
	tile := &Tile{Terrain: Farm, Zombies: 0, playerIds: []string{}}
	playerID := "player1"
	tile.addPlayer(playerID)
	assert.Contains(t, tile.playerIds, playerID, "Player should be added to tile")
	tile.removePlayer(playerID)
	assert.NotContains(t, tile.playerIds, playerID, "Player should be removed from tile")
}

func TestTile_findPlayerIdIndex(t *testing.T) {
	playerID1 := "player1"
	playerID2 := "player2"
	tile := &Tile{Terrain: Farm, Zombies: 0, playerIds: []string{playerID1, playerID2}}

	index, found := tile.findPlayerIdIndex(playerID2)
	assert.True(t, found, "Player ID should be found")
	assert.Equal(t, 1, index, "Correct index should be returned")

	_, found = tile.findPlayerIdIndex("nonexistent")
	assert.False(t, found, "Nonexistent player ID should not be found")
}

func TestTile_getMapPiece(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Direction: North}} // Player planning to move North
	tile := &Tile{Terrain: City, Zombies: 3, playerIds: []string{playerID}}

	mapPiece := tile.getMapPiece()

	assert.Equal(t, "City", mapPiece.TileType, "MapPiece should have correct TileType")
	assert.Equal(t, 3, mapPiece.ZombieCount, "MapPiece should have correct ZombieCount")
	assert.Equal(t, 1, mapPiece.PlayerCount, "MapPiece should have correct PlayerCount")
	assert.Equal(t, 1, mapPiece.PlayersPlanMoveNorth, "MapPiece should reflect player's planned North move")
	assert.Equal(t, 0, mapPiece.PlayersPlanMoveEast, "MapPiece should reflect player's planned East move")
	assert.Equal(t, 0, mapPiece.PlayersPlanMoveSouth, "MapPiece should reflect player's planned South move")
	assert.Equal(t, 0, mapPiece.PlayersPlanMoveWest, "MapPiece should reflect player's planned West move")
}

func TestTile_toStringTile(t *testing.T) {
	tile := Tile{Terrain: Forest, Zombies: 7, playerIds: []string{"player1", "player2"}}
	expectedString := "Forest 7 player1,player2"
	assert.Equal(t, expectedString, tile.toString(), "toString should return the correct format")
}

// --- player_test.go ---
func TestPlayer_hasWinCondition_Wins(t *testing.T) {
	player := Player{Cards: [5]Card{Research, Research, Research, None, None}} // Victory number is 2, player has 3
	assert.True(t, player.hasWinCondition(), "Player should win with enough research cards")
}

func TestPlayer_hasWinCondition_NoWin(t *testing.T) {
	player := Player{Cards: [5]Card{Research, Food, None, None, None}} // Victory number is 2, player has 1
	assert.False(t, player.hasWinCondition(), "Player should not win with insufficient research cards")
}

func TestPlayer_hasCardWhere_CardPresent(t *testing.T) {
	cards := [5]Card{Food, Wood, Weapon, None, None}
	index, found := hasCardWhere(cards[:], Wood)
	assert.True(t, found, "Card should be found")
	assert.Equal(t, 1, index, "Correct index should be returned")
}

func TestPlayer_hasCardWhere_CardAbsent(t *testing.T) {
	cards := [5]Card{Food, Wood, Weapon, None, None}
	index, found := hasCardWhere(cards[:], Dice)
	assert.False(t, found, "Card should not be found")
	assert.Equal(t, -1, index, "Index -1 should be returned when card not found")
}

func TestPlayer_hasCardWhere_MultipleCards(t *testing.T) {
	cards := [5]Card{Food, Wood, Wood, Weapon, None}
	index, found := hasCardWhere(cards[:], Wood)
	assert.True(t, found, "Card should be found")
	assert.Equal(t, 1, index, "Index of first occurrence should be returned")
}

func TestPlayer_firstIndexOfCardType_Present(t *testing.T) {
	player := Player{Cards: [5]Card{None, Wood, Food, Wood, None}}
	index := player.firstIndexOfCardType(Wood)
	assert.Equal(t, 1, index, "Should return index of first Wood card")
}

func TestPlayer_firstIndexOfCardType_Absent(t *testing.T) {
	player := Player{Cards: [5]Card{None, Food, None, None, None}}
	index := player.firstIndexOfCardType(Wood)
	assert.Equal(t, -1, index, "Should return -1 if Wood card is absent")
}

func TestPlayer_toStringPlayer(t *testing.T) {
	player := Player{ID: "testID", Name: "TestPlayer", X: 5, Y: 10, Cards: [5]Card{Food, Wood, None, None, None}}
	expectedString := "testIDTestPlayer: 5|10 FoodWoodNoneNone" // Expected string format
	assert.Equal(t, expectedString, player.toString(), "toString should return the correct format")
}

// --- game_test.go ---
func TestInitMap(t *testing.T) {
	r := newRandSource(1) // Deterministic random source for testing
	var testGameMap [mapWidth][mapHeight]*Tile
	initMap(r, &testGameMap)

	for x := 0; x < mapWidth; x++ {
		for y := 0; y < mapHeight; y++ {
			assert.NotNil(t, testGameMap[x][y], fmt.Sprintf("Tile at [%d][%d] should not be nil", x, y))
			assert.Contains(t, terrainTypes[:len(terrainTypes)-1], testGameMap[x][y].Terrain, "Tile at [%d][%d] should have a valid terrain type (excluding Edge)", x, y) // Exclude Edge
			assert.Equal(t, 0, testGameMap[x][y].Zombies, fmt.Sprintf("Tile at [%d][%d] should have 0 initial zombies", x, y))
			assert.Empty(t, testGameMap[x][y].playerIds, fmt.Sprintf("Tile at [%d][%d] should have no initial players", x, y))
		}
	}
}

func TestGetMapTile_WithinBounds(t *testing.T) {
	var testGameMap [mapWidth][mapHeight]*Tile
	testGameMap[50][50] = &Tile{Terrain: City, Zombies: 0, playerIds: []string{}} // Set a known tile
	gameMap = testGameMap // Set global gameMap for this test

	tile := getMapTile(50, 50, &gameMap)
	assert.Equal(t, City, tile.Terrain, "Should return the correct tile within bounds")
}

func TestGetMapTile_OutOfBounds(t *testing.T) {
	gameMap = [mapWidth][mapHeight]*Tile{} // Reset global gameMap
	tile := getMapTile(-1, -1, &gameMap)   // Out of bounds
	assert.Equal(t, Edge, tile.Terrain, "Should return an Edge tile when out of bounds")
	assert.Equal(t, -1, tile.Zombies, "Out of bounds tile zombies should be -1")
}

func TestMove_North(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, X: 50, Y: 50, Direction: North, Alive: true}} // Start at 50,50
	move(&playerMap)
	assert.Equal(t, 50, playerMap[playerID].X, "Player X coordinate should not change on North move")
	assert.Equal(t, 51, playerMap[playerID].Y, "Player Y coordinate should increase on North move")
	assert.Equal(t, defaultDirection, playerMap[playerID].Direction, "Player direction should reset after move")
}

func TestMove_SouthOutOfBounds(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, X: 50, Y: 0, Direction: South, Alive: true}} // Start at Y=0
	move(&playerMap)
	assert.Equal(t, 50, playerMap[playerID].X, "Player X coordinate should not change on South move")
	assert.Equal(t, 0, playerMap[playerID].Y, "Player Y coordinate should be clamped at 0") // Clamped to 0
}

func TestResources_FarmAndCity(t *testing.T) {
	playerID1 := "player1"
	playerID2 := "player2"
	playerMap = map[string]*Player{
		playerID1: {ID: playerID1, X: 50, Y: 50, Alive: true, Cards: [5]Card{None, None, None, None, None}},
		playerID2: {ID: playerID2, X: 60, Y: 60, Alive: true, Cards: [5]Card{None, None, None, None, None}},
	}
	gameMap = [mapWidth][mapHeight]*Tile{} // Reset gameMap
	gameMap[50][50] = &Tile{Terrain: Farm, Zombies: 0, playerIds: []string{playerID1}}   // Player1 on Farm
	gameMap[60][60] = &Tile{Terrain: City, Zombies: 0, playerIds: []string{playerID2}}   // Player2 on City
	resources()

	foodCardIndex, _ := hasCardWhere(playerMap[playerID1].Cards[:], Food)
	assert.NotEqual(t, -1, foodCardIndex, "Player1 should receive Food on Farm")
	weaponCardIndex, _ := hasCardWhere(playerMap[playerID2].Cards[:], Weapon)
	assert.NotEqual(t, -1, weaponCardIndex, "Player2 should receive Weapon on City")
}

func TestConsume_FoodCard(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Alive: true, Cards: [5]Card{Food, None, None, None, None}, Consume: Food}}
	gameMap = [mapWidth][mapHeight]*Tile{} // Reset gameMap
	consume(&playerMap, &gameMap)

	foodCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], Food)
	assert.Equal(t, -1, foodCardIndex, "Food card should be consumed")
	assert.True(t, playerMap[playerID].Alive, "Player should remain alive after consuming food")
}

func TestConsume_NoFoodCardDies(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Alive: true, Cards: [5]Card{None, None, None, None, None}, Consume: Food}} // No food
	gameMap = [mapWidth][mapHeight]*Tile{}                                                                         // Reset gameMap
	consume(&playerMap, &gameMap)

	assert.False(t, playerMap[playerID].Alive, "Player should die without food")
}

func TestConsume_WoodCardAttractsZombies(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Alive: true, Cards: [5]Card{Wood, None, None, None, None}, Consume: Wood, X: 50, Y: 50}}
	gameMap = [mapWidth][mapHeight]*Tile{} // Reset gameMap
	gameMap[50][50] = &Tile{Terrain: Farm, Zombies: 0, playerIds: []string{playerID}} // Player's tile
	gameMap[49][51] = &Tile{Terrain: Farm, Zombies: 2, playerIds: []string{}}         // SW tile with zombies
	gameMap[51][49] = &Tile{Terrain: Farm, Zombies: 3, playerIds: []string{}}         // NE tile with zombies

	consume(&playerMap, &gameMap)

	woodCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], Wood)
	assert.Equal(t, -1, woodCardIndex, "Wood card should be consumed")
	assert.Equal(t, 2, gameMap[50][50].Zombies, "Player's tile should gain attracted zombies") // 2 from SW, 0 from NN, 0 from NE, 0 from WW, 0 from EE, 3 from NE, 0 from SS, 0 from SE. Total 2, reduced by 1 each.
	assert.Equal(t, 1, gameMap[49][51].Zombies, "SW tile zombies should reduce by 1")       // 2-1 = 1
	assert.Equal(t, 2, gameMap[51][49].Zombies, "NE tile zombies should reduce by 1")       // 3-1 = 2

}

func TestGetHandSize_EmptyHand(t *testing.T) {
	player := Player{Cards: [5]Card{None, None, None, None, None}}
	assert.Equal(t, 0, getHandSize(&player), "Hand size should be 0 for empty hand")
}

func TestGetHandSize_FullHand(t *testing.T) {
	player := Player{Cards: [5]Card{Food, Wood, Weapon, Dice, Research}}
	assert.Equal(t, 5, getHandSize(&player), "Hand size should be 5 for full hand")
}

func TestLimitCards_OverLimitDiscardChosen(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Cards: [5]Card{Food, Wood, Weapon, Dice, Research}, Discard: Weapon}} // Hand over limit, discard Weapon
	limitCards(&playerMap)

	weaponCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], Weapon)
	assert.Equal(t, -1, weaponCardIndex, "Weapon card should be discarded")
	assert.Equal(t, 4, getHandSize(playerMap[playerID]), "Hand size should be limited to 4")
	assert.Equal(t, None, playerMap[playerID].Discard, "Discard should be reset to None")
}

func TestLimitCards_OverLimitDiscardNone(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, Cards: [5]Card{Food, Wood, Weapon, Dice, Research}, Discard: None}} // Hand over limit, no discard chosen
	limitCards(&playerMap)

	researchCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], Research) // Last card should be discarded if no discard chosen
	assert.Equal(t, -1, researchCardIndex, "Research card (last) should be discarded if no discard chosen")
	assert.Equal(t, 4, getHandSize(playerMap[playerID]), "Hand size should be limited to 4")
	assert.Equal(t, None, playerMap[playerID].Discard, "Discard should be reset to None")
}

func TestHandleCombat_Integration(t *testing.T) {
	playerID1 := "player1"
	playerID2 := "player2"
	playerMap = map[string]*Player{
		playerID1: {ID: playerID1, X: 50, Y: 50, Alive: true, Cards: [5]Card{Weapon, None, None, None, None}, Play: Weapon}, // Player1 wins
		playerID2: {ID: playerID2, X: 50, Y: 50, Alive: true, Cards: [5]Card{None, None, None, None, None}, Play: Dice},     // Player2 loses
	}
	gameMap = [mapWidth][mapHeight]*Tile{} // Reset gameMap
	gameMap[50][50] = &Tile{Terrain: City, Zombies: 15, playerIds: []string{playerID1, playerID2}} // Tile with players and zombies

	originalRollDice := rollDice
	rollDice = func() int { return 1 } // Mock dice to low value, Player2 will lose anyway
	defer func() { rollDice = originalRollDice }()

	handleCombat()

	assert.Equal(t, 0, gameMap[50][50].Zombies, "Zombies should be cleared after combat (player win)")
	assert.True(t, playerMap[playerID1].Alive, "Player1 should be alive (winner)")
	assert.False(t, playerMap[playerID2].Alive, "Player2 should be dead (loser)")
}

func TestSpreadFromSpreader_City(t *testing.T) {
	gameMap = [mapWidth][mapHeight]*Tile{} // Reset gameMap
	gameMap[50][50] = &Tile{Terrain: City, Zombies: 5}      // Spreader tile (City)
	gameMap[50][49] = &Tile{Terrain: Farm, Zombies: 2}      // North
	gameMap[51][50] = &Tile{Terrain: Farm, Zombies: 3}      // East
	gameMap[50][51] = &Tile{Terrain: Farm, Zombies: 4}      // South
	gameMap[49][50] = &Tile{Terrain: Farm, Zombies: 5}      // West

	spreadFromSpreader(&gameMap, 50, 50)

	assert.Equal(t, 3, gameMap[50][49].Zombies, "North tile zombies should increase by 1") // 2+1 = 3
	assert.Equal(t, 4, gameMap[51][50].Zombies, "East tile zombies should increase by 1")  // 3+1 = 4
	assert.Equal(t, 5, gameMap[50][51].Zombies, "South tile zombies should increase by 1") // 4+1 = 5
	assert.Equal(t, 6, gameMap[49][50].Zombies, "West tile zombies should increase by 1")  // 5+1 = 6
}

func TestSpreadFromSpreader_ZombieCutoff(t *testing.T) {
	gameMap = [mapWidth][mapHeight]*Tile{} // Reset gameMap
	gameMap[50][50] = &Tile{Terrain: Farm, Zombies: zombieCutoff} // Spreader tile (ZombieCutoff)
	gameMap[50][49] = &Tile{Terrain: Forest, Zombies: 2}    // North

	spreadFromSpreader(&gameMap, 50, 50)

	assert.Equal(t, 3, gameMap[50][49].Zombies, "North tile zombies should increase by 1 from zombie spreader") // 2+1 = 3
}

func TestSpread_Integration(t *testing.T) {
	gameMap = [mapWidth][mapHeight]*Tile{} // Reset gameMap
	gameMap[50][50] = &Tile{Terrain: City, Zombies: 5}      // Spreader 1
	gameMap[70][70] = &Tile{Terrain: Farm, Zombies: zombieCutoff} // Spreader 2
	gameMap[50][49] = &Tile{Terrain: Farm, Zombies: 2}      // Neighbor of Spreader 1
	gameMap[70][69] = &Tile{Terrain: Farm, Zombies: 2}      // Neighbor of Spreader 2

	spread(&gameMap)

	assert.GreaterOrEqual(t, gameMap[50][49].Zombies, 3, "Neighbor of City should have increased zombies")   // 2+1 = 3 or more if spread more than once
	assert.GreaterOrEqual(t, gameMap[70][69].Zombies, 3, "Neighbor of ZombieCutoff should have increased zombies") // 2+1 = 3 or more if spread more than once
}

func TestTick_Integration(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID, X: 50, Y: 50, Alive: true, Direction: North, Cards: [5]Card{Food, None, None, None, None}, Consume: Food}}
	gameMap = [mapWidth][mapHeight]*Tile{} // Reset gameMap
	gameMap[50][50] = &Tile{Terrain: Farm, Zombies: 2, playerIds: []string{playerID}} // Player on Farm
	initialFoodCardCount := getHandSize(playerMap[playerID])

	originalRollDice := rollDice
	rollDice = func() int { return 10 } // Mock dice to high value for combat - no player death
	defer func() { rollDice = originalRollDice }()

	tick(&gameMap, &playerMap)

	assert.Equal(t, 0, gameMap[50][50].Zombies, "Zombies should be cleared after combat (player wins)") // Player wins combat due to mock dice
	assert.Equal(t, 50, playerMap[playerID].X, "Player X should not change within tick")               // X coordinate move tested separately
	assert.Equal(t, 51, playerMap[playerID].Y, "Player Y should have moved North within tick")          // Y coordinate move tested separately
	foodCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], Food)
	assert.Equal(t, -1, foodCardIndex, "Food card should be consumed within tick") // Consume tested within tick
	resourceFoodCardIndex, _ := hasCardWhere(playerMap[playerID].Cards[:], Food)
	if initialFoodCardCount < 5 { // Only check for resource gain if hand was not full initially
		assert.NotEqual(t, -1, resourceFoodCardIndex, "Player should gain resource (Food) within tick") // Resource gain tested within tick
	}
}

func TestRandomizeBots_MovementAndConsume(t *testing.T) {
	botList = []*Player{{ID: "bot1", IsBot: true, Cards: [5]Card{Food, None, None, None, None}}} // Bot with food
	originalRandIntn := r.Intn
	r.Intn = func(n int) int { return 0 } // Mock rand.Intn to always return 0 for predictable direction selection
	defer func() { r.Intn = originalRandIntn }()

	randomizeBots(botList)

	assert.Equal(t, North, botList[0].Direction, "Bot direction should be randomized") // Mocked to North (index 0)
	assert.Equal(t, Dice, botList[0].Play, "Bot Play should be Dice")           // Bot Play is always Dice as per code
	assert.Equal(t, Food, botList[0].Consume, "Bot should consume Food if available")
}

func TestRandomizeBots_Discard(t *testing.T) {
	botList = []*Player{{ID: "bot1", IsBot: true, Cards: [5]Card{Food, Wood, Weapon, Dice, Research}}} // Full hand, needs discard
	originalHasCardWhere := hasCardWhereF // Capture original global function for reset
	hasCardWhereF = func(ar []Card, card Card) (int, bool) { return -1, false }   // Mock hasCardWhereF to always return false for "None" card not found
	defer func() { hasCardWhereF = originalHasCardWhere }()                           // Restore original global function

	randomizeBots(botList)

	assert.NotEqual(t, None, botList[0].Discard, "Bot should discard a card if hand is full")
}

func TestAddPlayer(t *testing.T) {
	playerMap = make(map[string]*Player) // Reset playerMap
	playerID := addPlayer("TestPlayer")

	assert.NotEmpty(t, playerID, "Player ID should be generated")
	player, ok := playerMap[playerID]
	assert.True(t, ok, "Player should be added to playerMap")
	assert.Equal(t, "TestPlayer", player.Name, "Player name should be set correctly")
	assert.True(t, player.Alive, "Player should be initialized as alive")
	assert.False(t, player.IsBot, "Player should be initialized as not a bot")
	assert.Equal(t, defaultDirection, player.Direction, "Player direction should be default")
	assert.Equal(t, Card(None), player.Play, "Player Play should be None by default")
	assert.Equal(t, Card(None), player.Consume, "Player Consume should be None by default")
	assert.Equal(t, Card(None), player.Discard, "Player Discard should be None by default")
	initialCards := [5]Card{Food, Wood, Wood, None, None}
	assert.Equal(t, initialCards, player.Cards, "Player cards should be initialized to default")
}

func TestAddBot(t *testing.T) {
	playerMap = make(map[string]*Player) // Reset playerMap
	botList = []*Player{}              // Reset botList
	bID := 1
	addBot(&playerMap, &botList, &bID)

	botIDStr := "1"
	bot, ok := playerMap[botIDStr]
	assert.True(t, ok, "Bot should be added to playerMap")
	assert.Contains(t, botList, bot, "Bot should be added to botList")
	assert.Equal(t, botIDStr, bot.ID, "Bot ID should be set correctly")
	assert.True(t, bot.IsBot, "Should be initialized as a bot")
	assert.True(t, bot.Alive, "Bot should be initialized as alive")
	assert.Equal(t, Stay, bot.Direction, "Bot direction should be Stay by default")
	assert.Equal(t, Card(None), bot.Play, "Bot Play should be None by default")
	assert.Equal(t, Card(None), bot.Consume, "Bot Consume should be None by default")
	assert.Equal(t, Card(None), bot.Discard, "Bot Discard should be None by default")
	initialCards := [5]Card{Food, Wood, Wood, None, None}
	assert.Equal(t, initialCards, bot.Cards, "Bot cards should be initialized to default")
}

func TestRestockBots_AddsBots(t *testing.T) {
	playerMap = make(map[string]*Player) // Reset playerMap
	botList = []*Player{}              // Reset botList
	botNumber = 3                       // Target bot number
	bID := 0

	restockBots(&playerMap, &botList, &bID)

	assert.Len(t, botList, 3, "Bot list should be restocked to botNumber")
	assert.Equal(t, 3, len(playerMap), "Player map should contain botNumber bots")
	assert.Equal(t, 3, bID, "botID should be incremented correctly")
}

func TestRestockBots_NoBotsNeeded(t *testing.T) {
	playerMap = make(map[string]*Player)          // Reset playerMap
	botList = []*Player{{}, {}, {}}               // Already at botNumber count
	initialBotListLen := len(botList)              // Store initial length
	botNumber = 3                                   // Target bot number is same as current
	bID := 3                                       // Start botID beyond existing bots

	restockBots(&playerMap, &botList, &bID)

	assert.Len(t, botList, initialBotListLen, "Bot list should remain unchanged if no restock needed")
	assert.Len(t, playerMap, initialBotListLen, "Player map should remain unchanged if no restock needed")
	assert.Equal(t, 3, bID, "botID should not change if no bots added") // Important to not increase bID unecessarily

}

func TestHavePlayersWon_WinConditionMet(t *testing.T) {
	playerMap = map[string]*Player{
		"player1": {Cards: [5]Card{Research, Research, None, None, None}}, // Player 1 wins
		"player2": {Cards: [5]Card{Research, None, None, None, None}},        // Player 2 does not win
	}

	won := havePlayersWon(playerMap)
	assert.True(t, won, "havePlayersWon should return true if any player meets win condition")
}

func TestHavePlayersWon_NoWinConditionMet(t *testing.T) {
	playerMap = map[string]*Player{
		"player1": {Cards: [5]Card{Research, None, None, None, None}}, // Player 1 does not win
		"player2": {Cards: [5]Card{Research, None, None, None, None}}, // Player 2 does not win
	}

	won := havePlayersWon(playerMap)
	assert.False(t, won, "havePlayersWon should return false if no player meets win condition")
}

func TestGetPlayerOrNil_PlayerExists(t *testing.T) {
	playerID := "player1"
	playerMap = map[string]*Player{playerID: {ID: playerID}}

	player := getPlayerOrNil(playerID)
	assert.NotNil(t, player, "getPlayerOrNil should return player if exists")
	assert.Equal(t, playerID, player.ID, "getPlayerOrNil should return correct player")
}

func TestGetPlayerOrNil_PlayerDoesNotExist(t *testing.T) {
	playerMap = map[string]*Player{} // Empty player map

	player := getPlayerOrNil("nonexistent")
	assert.Nil(t, player, "getPlayerOrNil should return nil if player does not exist")
}

func TestGetSurroundingsOfPlayer_CenterMap(t *testing.T) {
	playerID := "player1"
	gameMap = [mapWidth][mapHeight]*Tile{} // Reset gameMap
	playerMap = map[string]*Player{playerID: {ID: playerID, X: 50, Y: 50}}
	gameMap[50][50] = &Tile{Terrain: City, Zombies: 1} // Center tile

	surroundings, found := getSurroundingsOfPlayer(playerID)

	assert.True(t, found, "getSurroundingsOfPlayer should return true if player exists")
	assert.Equal(t, "City", surroundings.CE.TileType, "Center tile should be City")
	assert.Equal(t, 1, surroundings.CE.ZombieCount, "Center tile should have 1 zombie")
	assert.NotEmpty(t, surroundings.NW.TileType, "NW tile type should be populated") // Just check one neighbor, structure is similar
}

func TestGetSurroundingsOfPlayer_PlayerNil(t *testing.T) {
	playerMap = map[string]*Player{} // Empty player map

	_, found := getSurroundingsOfPlayer("nonexistent")
	assert.False(t, found, "getSurroundingsOfPlayer should return false if player is nil")
}

// newRandSource creates a deterministic rand.Source for testing.
func newRandSource(seed int64) *rand.Rand {
	source := rand.NewSource(seed)
	return rand.New(source)
}

// hasCardWhereF is a global variable to allow mocking of hasCardWhere function
// in tests where you need to control its behavior, especially when it's used
// in functions that are being unit tested and you want to isolate them.
var hasCardWhereF = hasCardWhere
