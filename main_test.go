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

