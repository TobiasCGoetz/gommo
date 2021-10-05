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
