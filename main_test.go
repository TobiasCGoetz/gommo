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

func fakeInitMap(gMap *[mapWidth][mapHeight]*Tile) {
	for a, column := range gMap {
		for b, _ := range column {
			var tile = Tile{City, 0}
			gMap[a][b] = &tile
		}
	}
}

func TestCreateCityList(t *testing.T) {
	var gameMap [mapWidth][mapHeight]*Tile
	fakeInitMap(&gameMap)
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
}