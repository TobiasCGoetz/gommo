package main

import (
	"testing"
)

func TestInitMap (t *testing.T) {
	var testMap [mapWidth][mapHeight]*Tile
	initMap(&testMap)
	for _, column := range testMap {
		for _, tile := range column {
			if tile == nil {
				t.Errorf("Map tile missing.")
			}
		}
	}
}

func fakeInitMap(gMap *[mapWidth][mapHeight]*Tile, terrain Terrain) {
	for a, column := range gMap {
		for b, _ := range column {
			var tile = Tile{terrain, 0}
			gMap[a][b] = &tile
		}
	}
}

func TestCreateCityList(t *testing.T) {
	var gameMap [mapWidth][mapHeight]*Tile
	fakeInitMap(&gameMap, City)
	var cityList = createCityList(&gameMap)
	var count = 0
	for x := 0; x < mapWidth; x++ {
		for y := 0; y < mapHeight; y++ {
			if cityList[count].X != x || cityList[count].Y != y {
				t.Errorf("We're missing [%d][%d] in the city list", x, y)
			}
			count++
		}
	}
}

//TODO: Figure out if positive Y == North is stupid
func TestMove(t *testing.T) {
	var playerX = 5
	var playerY = 5
	var testPlayer = Player{
		ID:        "test",
		X:         playerX,
		Y:         playerY,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   Wood,
		Cards:     [5]Card{Wood, Wood, Food, None, None},
		Alive:     true,
		IsBot:	   true,
	}
	var testArray = []*Player{&testPlayer}
	move(&testArray)
	if testPlayer.Y != playerY+1 {
		t.Errorf("Move north failed.")
	}
	testPlayer.Direction = South
	move(&testArray)
	if testPlayer.Y != playerY {
		t.Errorf("Move south failed.")
	}
	testPlayer.Direction = East
	move(&testArray)
	if testPlayer.X != playerX+1 {
		t.Errorf("Move east failed.")
	}
	testPlayer.Direction = West
	move(&testArray)
	if testPlayer.X != playerX {
		t.Errorf("Move west failed.")
	}
	testPlayer.Direction = Stay
	testPlayer.X = playerX
	testPlayer.Y = playerY
	move(&testArray)
	if testPlayer.X != playerX || testPlayer.Y != playerY {
		t.Errorf("Staying in place failed.")
	}
	var deadPlayer = Player{
		ID:        "test",
		X:         playerX,
		Y:         playerY,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   Wood,
		Cards:     [5]Card{Wood, Wood, Food, None, None},
		Alive:     false,
		IsBot:     true,
	}
	testArray = []*Player{&deadPlayer}
	move(&testArray)
	if deadPlayer.Y != playerY {
		t.Errorf("Dead player wasn't supposed to move %d tiles north", deadPlayer.Y-5)
	}
}

func TestResources(t *testing.T) {
	var gameMap [mapWidth][mapHeight]*Tile
	fakeInitMap(&gameMap, Farm)
	var testPlayer = Player{
		ID:        "test",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   None,
		Cards:     [5]Card{ None, None, None, None, None },
		Alive:     true,
		IsBot:     true,
	}
	var testArray = []*Player{&testPlayer}
	resources(&testArray, &gameMap)
	if testPlayer.Cards[0] != Food {
		t.Log(testPlayer.Cards)
		t.Errorf("TestPlayer received the wrong card - expected Food but got %s", testPlayer.Cards[0].toString())
	}
	//Dead players sit out
	var deadPlayer = Player{
		ID:        "test",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   None,
		Cards:     [5]Card{None, None, None, None, None},
		Alive:     false,
		IsBot:     true,
	}
	testArray = []*Player{&deadPlayer}
	resources(&testArray, &gameMap)
	if deadPlayer.Cards[0] != None {
		t.Errorf("Dead player did not have to sit out when distributing resources.")
	}
}

func TestConsume(t *testing.T) {
	var testPlayer = Player{
		ID:        "test",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   None,
		Cards:     [5]Card{None, None, None, None, Wood},
		Alive:     true,
		IsBot:     true,
	}
	var testArray = []*Player{&testPlayer}
	consume(&testArray)
	if testPlayer.Cards[4] != None {
		t.Errorf("Player was not supposed to have resources remaining.")
	}
	var deadPlayer = Player{
		ID:        "test",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   None,
		//TODO: Test more combinations
		Cards:     [5]Card{None, None, None, None, Wood},
		Alive:     false,
		IsBot:     true,
	}
	testArray = []*Player{&deadPlayer}
	consume(&testArray)
	if deadPlayer.Cards[4] != Wood {
		t.Errorf("Dead player's not supposed to consume anything.")
	}
}

func TestHandSize (t *testing.T) {
	var testCases = [12][5]Card{
		{None, None, None, None, None},
		{Wood, None, None, None, None},
		{Food, None, None, None, None},
		{Weapon, None, None, None, None},
		{None, Wood, None, None, None},
		{None, None, Wood, None, None},
		{None, None, None, Wood, None},
		{None, None, None, None, Wood},
		{Wood, Wood, None, None, None},
		{Wood, Wood, Wood, None, None},
		{Wood, Wood, Wood, Wood, None},
		{Wood, Wood, Wood, Wood, Wood},
	}
	var testResults = [12]int {
		0, 1, 1, 1, 1, 1, 1, 1, 2, 3, 4, 5,
	}
	var testPlayer = Player{
		ID:        "test",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   None,
		Cards: [5]Card{None, None, None, None, None},
		Alive:     false,
		IsBot:     true,
	}
	for testNumber, cards := range testCases {
		testPlayer.Cards = cards
		if testResults[testNumber] != getHandSize(testPlayer) {
			t.Errorf("[%v] is not %d cards", testPlayer.Cards, testResults[testNumber])
		}
	}
}

func TestLimitCards (t *testing.T) {
	var testCases = [6][5]Card {
		{None, None, None, None, None},
		{Wood, None, None, None, None},
		{Wood, Wood, None, None, None},
		{Wood, Wood, Wood, None, None},
		{Wood, Wood, Wood, Wood, None},
		{Wood, Wood, Wood, Wood, Food},
	}
	var testResults = [6]int {
		0, 1, 2, 3, 4, 4,
	}
	var testPlayer = Player{
		ID:        "test",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   Wood,
		Cards: [5]Card{None, None, None, None, None},
		Alive:     false,
		IsBot:     true,
	}
	var testArray = []*Player{&testPlayer}
	//TODO: Add check if right card was removed
	for testNumber, cards := range testCases {
		testPlayer.Cards = cards
		//var pDiscard = testPlayer.discard
		limitCards(&testArray)
		if (getHandSize(testPlayer)) != testResults[testNumber] {
			t.Errorf("[%v] is not %d cards", testPlayer.Cards, testResults[testNumber])
		}
	}
}

func TestRestockBots (t *testing.T) {
	var testPlayerList []*Player
	var testBotList []*Player
	var testCases [5]int = [5]int{0, 30, 49, 50, 51}
	var botID = 0
	var testBot = Player{
		ID:        "test",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   Wood,
		Cards: [5]Card{None, None, None, None, None},
		Alive:     false,
		IsBot:     true,
	}
	for i := 0; i < len(testCases); i++ {
		testPlayerList = nil
		testBotList = nil
		//Add bots
		for j := 0; j < testCases[i]; j++ {
			testPlayerList = append(testPlayerList, &testBot)
			testBotList = append(testBotList, &testBot)
		}
		//Call restocking routine
		restockBots(&testPlayerList, &testBotList, &botID)
		//Check correct amount restocked
		if len(testBotList) < 50 {
			t.Errorf("Bot restocking off by %d.", 50-len(testBotList))
		}
	}
}