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
			if cityList[count].x != x || cityList[count].y != y {
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
		id:        "test",
		x:         playerX,
		y:         playerY,
		direction: North,
		play:      Dice,
		consume:   Wood,
		discard:   Wood,
		cards:     [5]Card{Wood, Wood, Food, None, None},
		alive:     true,
	}
	var testArray = []*Player{&testPlayer}
	move(&testArray)
	if testPlayer.y != playerY+1 {
		t.Errorf("Move north failed.")
	}
	testPlayer.direction = South
	move(&testArray)
	if testPlayer.y != playerY {
		t.Errorf("Move south failed.")
	}
	testPlayer.direction = East
	move(&testArray)
	if testPlayer.x != playerX+1 {
		t.Errorf("Move east failed.")
	}
	testPlayer.direction = West
	move(&testArray)
	if testPlayer.x != playerX {
		t.Errorf("Move west failed.")
	}
	var deadPlayer = Player{
		id:        "test",
		x:         playerX,
		y:         playerY,
		direction: North,
		play:      Dice,
		consume:   Wood,
		discard:   Wood,
		cards:     [5]Card{Wood, Wood, Food, None, None},
		alive:     false,
	}
	testArray = []*Player{&deadPlayer}
	move(&testArray)
	if deadPlayer.y != playerY {
		t.Errorf("Dead player wasn't supposed to move %d tiles north", deadPlayer.y-5)
	}
}

func TestResources(t *testing.T) {
	var gameMap [mapWidth][mapHeight]*Tile
	fakeInitMap(&gameMap, Farm)
	var testPlayer = Player{
		id:        "test",
		x:         5,
		y:         5,
		direction: North,
		play:      Dice,
		consume:   Wood,
		discard:   None,
		cards:     [5]Card{ None, None, None, None, None },
		alive:     true,
	}
	var testArray = []*Player{&testPlayer}
	resources(&testArray, &gameMap)
	if testPlayer.cards[0] != Food {
		t.Log(testPlayer.cards)
		t.Errorf("TestPlayer received the wrong card - expected Food but got %s", testPlayer.cards[0].toString())
	}
	//Dead players sit out
	var deadPlayer = Player{
		id:        "test",
		x:         5,
		y:         5,
		direction: North,
		play:      Dice,
		consume:   Wood,
		discard:   None,
		cards:     [5]Card{None, None, None, None, None},
		alive:     false,
	}
	testArray = []*Player{&deadPlayer}
	resources(&testArray, &gameMap)
	if deadPlayer.cards[0] != None {
		t.Errorf("Dead player did not have to sit out when distributing resources.")
	}
}

func TestConsume(t *testing.T) {
	var testPlayer = Player{
		id:        "test",
		x:         5,
		y:         5,
		direction: North,
		play:      Dice,
		consume:   Wood,
		discard:   None,
		cards:     [5]Card{None, None, None, None, Wood},
		alive:     true,
	}
	var testArray = []*Player{&testPlayer}
	consume(&testArray)
	if testPlayer.cards[4] != None {
		t.Errorf("Player was not supposed to have resources remaining.")
	}
	var deadPlayer = Player{
		id:        "test",
		x:         5,
		y:         5,
		direction: North,
		play:      Dice,
		consume:   Wood,
		discard:   None,
		//TODO: Test more combinations
		cards:     [5]Card{None, None, None, None, Wood},
		alive:     false,
	}
	testArray = []*Player{&deadPlayer}
	consume(&testArray)
	if deadPlayer.cards[4] != Wood {
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
		id:        "test",
		x:         5,
		y:         5,
		direction: North,
		play:      Dice,
		consume:   Wood,
		discard:   None,
		cards: [5]Card{None, None, None, None, None},
		alive:     false,
	}
	for testNumber, cards := range testCases {
		testPlayer.cards = cards
		if testResults[testNumber] != getHandSize(testPlayer) {
			t.Errorf("[%v] is not %d cards", testPlayer.cards, testResults[testNumber])
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
		id:        "test",
		x:         5,
		y:         5,
		direction: North,
		play:      Dice,
		consume:   Wood,
		discard:   Wood,
		cards: [5]Card{None, None, None, None, None},
		alive:     false,
	}
	var testArray = []*Player{&testPlayer}
	//TODO: Add check if right card was removed
	for testNumber, cards := range testCases {
		testPlayer.cards = cards
		//var pDiscard = testPlayer.discard
		limitCards(&testArray)
		if (getHandSize(testPlayer)) != testResults[testNumber] {
			t.Errorf("[%v] is not %d cards", testPlayer.cards, testResults[testNumber])
		}
	}
}