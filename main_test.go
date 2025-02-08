package main

import (
	"math/rand"
	"testing"
)

func TestInitMap(t *testing.T) {
	var testMap [mapWidth][mapHeight]*Tile
	r = rand.New(rand.NewSource(10))
	initMap(*r, &testMap)
	for _, column := range testMap {
		for _, tile := range column {
			if tile == nil {
				t.Errorf("Map tile missing.")
			}
		}
	}
}

func fakeInitMap(terrain Terrain, zombieNr int) [mapWidth][mapHeight]*Tile {
	fakeMap := [mapWidth][mapHeight]*Tile{}
	for a, column := range fakeMap {
		for b := range column {
			var tile = Tile{terrain, zombieNr, []string{}}
			fakeMap[a][b] = &tile
		}
	}
	return fakeMap
}

func TestCreateCityList(t *testing.T) {
	gameMap := fakeInitMap(City, 0)
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

func TestGetMapTile(t *testing.T) {
	var testMap [mapWidth][mapHeight]*Tile
	r = rand.New(rand.NewSource(10))
	initMap(*r, &testMap)
	testMap[13][42] = &Tile{City, 99, []string{}}
	var testTile = getMapTile(13, 42, &testMap)
	if testTile.Terrain == City && testTile.Zombies == 99 {
		return
	} else {
		t.Errorf("GetMapTile picked the wrong tile.")
	}
}

// TODO: Figure out if positive Y == North is stupid
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
		IsBot:     true,
	}
	var testPlayerMap = make(map[string]*Player)
	testPlayerMap[testPlayer.ID] = &testPlayer
	//Test directions
	move(&testPlayerMap)
	if testPlayer.Y != playerY+1 {
		t.Errorf("Move north failed.")
	}
	testPlayer.Direction = South
	testPlayerMap[testPlayer.ID] = &testPlayer
	move(&testPlayerMap)
	if testPlayer.Y != playerY {
		t.Errorf("Move south failed.")
	}
	testPlayer.Direction = East
	move(&testPlayerMap)
	if testPlayer.X != playerX+1 {
		t.Errorf("Move east failed.")
	}
	testPlayer.Direction = West
	move(&testPlayerMap)
	if testPlayer.X != playerX {
		t.Errorf("Move west failed.")
	}
	testPlayer.Direction = Stay
	testPlayer.X = playerX
	testPlayer.Y = playerY
	move(&testPlayerMap)
	if testPlayer.X != playerX || testPlayer.Y != playerY {
		t.Errorf("Staying in place failed.")
	}
	//Test edge behaviour
	testPlayer.X = 2 * mapWidth
	testPlayer.Direction = East
	move(&testPlayerMap)
	if testPlayer.X != mapWidth-1 {
		t.Errorf("Player (%d / %d) is %d out of bounds east", testPlayer.X, testPlayer.Y, testPlayer.X-mapWidth-1)
	}
	testPlayer.X = -10
	testPlayer.Direction = West
	move(&testPlayerMap)
	if testPlayer.X != 0 {
		t.Errorf("Player is %d out of bounds west", testPlayer.X)
	}
	testPlayer.X = 10
	testPlayer.Y = 2 * mapHeight
	testPlayer.Direction = South
	move(&testPlayerMap)
	if testPlayer.Y != mapHeight-1 {
		t.Errorf("Player is %d out of bounds north", testPlayer.Y-mapHeight-1)
	}
	testPlayer.Y = -10
	testPlayer.Direction = North
	move(&testPlayerMap)
	if testPlayer.Y != 0 {
		t.Errorf("Player is %d out of bounds south", testPlayer.Y)
	}
	//Test dead player
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
	testPlayerMap = make(map[string]*Player)
	testPlayerMap[deadPlayer.ID] = &deadPlayer
	move(&testPlayerMap)
	if deadPlayer.Y != playerY {
		t.Errorf("Dead player wasn't supposed to move %d tiles north", deadPlayer.Y-5)
	}
}

func TestMoveCycle(t *testing.T) {
	t.Log("Beginning TestMoveCylce")
	var playerX, playerY = 10, 10
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
		IsBot:     true,
	}
	var referencePlayer = testPlayer
	var testPlayerMap = make(map[string]*Player)
	testPlayerMap[testPlayer.ID] = &testPlayer
	var possibleDirections = Directions
	//Test directions
	for _, dir := range possibleDirections {
		t.Logf("%s", dir.toString())
		testPlayer.Direction = dir
		move(&testPlayerMap)
		t.Logf("%d | %d", testPlayer.X, testPlayer.Y)
	}

	if testPlayer.X != referencePlayer.X || testPlayer.Y != referencePlayer.Y {
		t.Errorf("TestPlayer moved away instead of returning home: %d | %d",
			testPlayer.X-referencePlayer.X,
			testPlayer.Y-referencePlayer.Y)
	}
}

func TestGetFirstEmptyHandSlot(t *testing.T) {
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
		Cards:     [5]Card{None, None, None, None, None},
		Alive:     true,
		IsBot:     true,
	}
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
	var testResults = [12]int{
		0, 1, 1, 1, 0, 0, 0, 0, 2, 3, 4, -1,
	}
	for testNumber, cards := range testCases {
		testPlayer.Cards = cards
		if testResults[testNumber] != getFirstEmptyHandSlot(testPlayer.Cards) {
			t.Errorf("[%v] first empty slot is not %d", testPlayer.Cards, testResults[testNumber])
		}
	}
}

func TestFight(t *testing.T) {
	r = rand.New(rand.NewSource(1))
	var loc1 = IntTuple{10, 10}
	var loc2 = IntTuple{13, 13}
	var loc3 = IntTuple{21, 21}
	gameMap := fakeInitMap(Forest, 3)
	gameMap[loc2.X][loc2.Y].Zombies = zombieCutoff
	gameMap[loc3.X][loc3.Y].Zombies = weaponStrength - 1
	var combatGroups = make(map[IntTuple][]*Player)
	var testPlayer = Player{
		ID:        "test",
		X:         10,
		Y:         10,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   Wood,
		Cards:     [5]Card{None, None, None, None, None},
		Alive:     true,
		IsBot:     true,
	}
	combatGroups[loc1] = append(combatGroups[loc1], &testPlayer, &testPlayer, &testPlayer, &testPlayer)
	fight(&gameMap, combatGroups[loc1])
	if gameMap[loc1.X][loc1.Y].Zombies > 0 {
		t.Errorf("%d/%d has %d Z instead of 0", loc1.X, loc1.Y, gameMap[loc1.X][loc1.Y].Zombies)
	}
	testPlayer.X = loc2.X
	testPlayer.Y = loc2.Y
	combatGroups[loc2] = append(combatGroups[loc2], &testPlayer)
	fight(&gameMap, combatGroups[loc2])
	if testPlayer.Alive {
		t.Errorf("Player didn't die correctly")
	}
	testPlayer.Alive = true
	testPlayer.X = loc3.X
	testPlayer.Y = loc3.Y
	testPlayer.Cards = [5]Card{Weapon, None, None, None, None}
	testPlayer.Play = Weapon
	combatGroups[loc3] = append(combatGroups[loc3], &testPlayer)
	fight(&gameMap, combatGroups[loc3])
	if !testPlayer.Alive || gameMap[loc3.X][loc3.Y].Zombies != 0 {
		t.Errorf("Player didn't use weapon correctly")
	}
}

func TestHandleCombat(t *testing.T) {
	r = rand.New(rand.NewSource(1))
	var loc = IntTuple{10, 10}
	gameMap := fakeInitMap(Forest, 0)
	gameMap[loc.X][loc.Y] = &Tile{City, 99, []string{}}
	var testPlayer1 = Player{
		ID:        "test1",
		X:         loc.X,
		Y:         loc.Y,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   Wood,
		Cards:     [5]Card{None, None, None, None, None},
		Alive:     true,
		IsBot:     true,
	}
	var testPlayer2 = Player{
		ID:        "test2",
		X:         loc.X,
		Y:         loc.Y,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   Wood,
		Cards:     [5]Card{None, None, None, None, None},
		Alive:     true,
		IsBot:     true,
	}
	var testPlayer3 = Player{
		ID:        "test3",
		X:         loc.X,
		Y:         loc.Y,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   Wood,
		Cards:     [5]Card{None, None, None, None, None},
		Alive:     true,
		IsBot:     true,
	}
	var testPlayerMap = make(map[string]*Player)
	testPlayerMap[testPlayer1.ID] = &testPlayer1
	testPlayerMap[testPlayer2.ID] = &testPlayer2
	testPlayerMap[testPlayer3.ID] = &testPlayer3
	handleCombat(&gameMap, &testPlayerMap)
	if testPlayer1.Alive || testPlayer2.Alive || testPlayer3.Alive {
		t.Errorf("Not all players died.")
	}
}

func TestSpread(t *testing.T) {
	gameMap := fakeInitMap(Forest, 1)
	gameMap[10][10] = &Tile{City, zombieCutoff, []string{}}
	gameMap[12][10] = &Tile{City, zombieCutoff, []string{}}
	gameMap[99][99] = &Tile{City, 4, []string{}}
	var cityList = createCityList(&gameMap)
	spread(&gameMap, &cityList)
	var testCases = []IntTuple{
		//{9, 11},
		{10, 11},
		//{11, 11},
		{12, 11},
		//{13, 11},
		{9, 10},
		{11, 10},
		{13, 10},
		//{9, 9},
		{10, 9},
		//{11, 9},
		{12, 9},
		//{13, 9},
		{99, 99},
	}
	var testResults = []int{
		2, 2, //2, 2, 3, 2, 2,
		2, 3, 2, //2, 3, 2,
		2, 2, //2 ,2, 3, 2, 2,
		5,
	}
	for testNumber, coords := range testCases {
		if testResults[testNumber] != gameMap[coords.X][coords.Y].Zombies {
			t.Errorf("%d/%d %d zombies instead of %d", coords.X, coords.Y, gameMap[coords.X][coords.Y].Zombies, testResults[testNumber])
		}
	}
}

// TODO: Test player starves
func TestConsumeWoodAttracts(t *testing.T) {
	gameMap := fakeInitMap(Forest, 1)
	var playerX = 5
	var playerY = 5
	var testPlayer = Player{
		ID:        "testPlayer",
		X:         playerX,
		Y:         playerY,
		Direction: Stay,
		Play:      Dice,
		Consume:   Wood,
		Discard:   None,
		Cards:     [5]Card{Wood, None, None, None, None},
		Alive:     true,
		IsBot:     true,
	}
	var testPlayerMap = make(map[string]*Player)
	testPlayerMap[testPlayer.ID] = &testPlayer
	consume(&testPlayerMap, &gameMap)
	if gameMap[playerX][playerY].Zombies != 9 {
		t.Errorf("Tile %d/%d was meant to have 9 zombies not %d", playerX, playerY, gameMap[playerX][playerY].Zombies)
	}
}

// TODO: Test other tiles but Farm
func TestResources(t *testing.T) {
	gameMap := fakeInitMap(Farm, 0)
	var testPlayer = Player{
		ID:        "test",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   None,
		Cards:     [5]Card{None, None, None, None, None},
		Alive:     true,
		IsBot:     true,
	}
	var testPlayerMap = make(map[string]*Player)
	testPlayerMap[testPlayer.ID] = &testPlayer
	resources(&testPlayerMap, gameMap)
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
	testPlayerMap = make(map[string]*Player)
	testPlayerMap[deadPlayer.ID] = &deadPlayer
	resources(&testPlayerMap, gameMap)
	if deadPlayer.Cards[0] != None {
		t.Errorf("Dead player did not have to sit out when distributing resources.")
	}
}

func TestConsume(t *testing.T) {
	gameMap := fakeInitMap(City, 0)
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
	var testPlayerMap = make(map[string]*Player)
	testPlayerMap[testPlayer.ID] = &testPlayer
	consume(&testPlayerMap, &gameMap)
	if testPlayerMap[testPlayer.ID].Cards[4] != None {
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
		Cards: [5]Card{None, None, None, None, Wood},
		Alive: false,
		IsBot: true,
	}
	testPlayerMap = make(map[string]*Player)
	testPlayerMap[deadPlayer.ID] = &deadPlayer
	consume(&testPlayerMap, &gameMap)
	if deadPlayer.Cards[4] != Wood {
		t.Errorf("Dead player's not supposed to consume anything.")
	}
}

func TestHandSize(t *testing.T) {
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
	var testResults = [12]int{
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
		Cards:     [5]Card{None, None, None, None, None},
		Alive:     false,
		IsBot:     true,
	}
	for testNumber, cards := range testCases {
		testPlayer.Cards = cards
		if testResults[testNumber] != getHandSize(&testPlayer) {
			t.Errorf("[%v] is not %d cards", testPlayer.Cards, testResults[testNumber])
		}
	}
}

func TestLimitCards(t *testing.T) {
	var testCases = [6][5]Card{
		{None, None, None, None, None},
		{Wood, None, None, None, None},
		{Wood, Wood, None, None, None},
		{Wood, Wood, Wood, None, None},
		{Wood, Wood, Wood, Wood, None},
		{Wood, Wood, Wood, Wood, Food},
	}
	var testResults = [6]int{
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
		Cards:     [5]Card{None, None, None, None, None},
		Alive:     false,
		IsBot:     true,
	}
	//TODO: Add check if right card was removed
	for testNumber, cards := range testCases {
		var testPlayerMap = make(map[string]*Player)
		testPlayer.Cards = cards
		testPlayer.Discard = Wood
		testPlayerMap[testPlayer.ID] = &testPlayer
		limitCards(&testPlayerMap)
		if (getHandSize(testPlayerMap[testPlayer.ID])) != testResults[testNumber] {
			t.Errorf("[%v] is not %d cards", testPlayer.Cards, testResults[testNumber])
		}
	}
}

func TestPlayerConsumeFallback(t *testing.T) {
	gameMap := fakeInitMap(City, 0)
	var testPlayer = Player{
		ID:        "testPlayer",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   None,
		Discard:   None,
		Cards:     [5]Card{Food, None, None, None, None},
		Alive:     true,
		IsBot:     true,
	}
	var testPlayerMap = make(map[string]*Player)
	testPlayerMap[testPlayer.ID] = &testPlayer
	consume(&testPlayerMap, &gameMap)
	var _, hasCard = playerHasCard(&testPlayer, Wood)
	if hasCard {
		t.Errorf("Consume fallback to Food failed")
	}

	//Test for Wood fallback
	testPlayer = Player{
		ID:        "testPlayer",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   None,
		Discard:   None,
		Cards:     [5]Card{Wood, None, None, None, None},
		Alive:     true,
		IsBot:     true,
	}
	testPlayerMap = make(map[string]*Player)
	testPlayerMap[testPlayer.ID] = &testPlayer
	consume(&testPlayerMap, &gameMap)
	_, hasCard = playerHasCard(&testPlayer, Wood)
	if hasCard {
		t.Errorf("Consume fallback to Wood failed")
	}

}

func TestRestockBots(t *testing.T) {
	r = rand.New(rand.NewSource(10))
	var testPlayerMap map[string]*Player
	var testBotList []*Player
	var testCases = [5]int{0, 30, 49, 50, 51}
	var botID = 0
	var testBot = Player{
		ID:        "test",
		X:         5,
		Y:         5,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   Wood,
		Cards:     [5]Card{None, None, None, None, None},
		Alive:     false,
		IsBot:     true,
	}
	for i := 0; i < len(testCases); i++ {
		testPlayerMap = make(map[string]*Player)
		testBotList = nil
		//Add bots
		for j := 0; j < testCases[i]; j++ {
			testPlayerMap[testBot.ID] = &testBot
			testBotList = append(testBotList, &testBot)
		}
		//Call restocking routine
		restockBots(&testPlayerMap, &testBotList, &botID)
		//Check correct amount restocked
		if len(testBotList) < botNumber {
			t.Errorf("Bot restocking off by %d.", botNumber-len(testBotList))
		}
	}
}
