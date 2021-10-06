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
