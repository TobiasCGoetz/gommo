package main

import (
	"math/rand"
	"sync"
)

type gameMap struct {
	gMap   [][]*Tile
	width  int
	height int
}

func NewGameMap() gameMap {
	width := gameConfig.Map.Width
	height := gameConfig.Map.Height
	
	instance := gameMap{
		width:  width,
		height: height,
		gMap:   make([][]*Tile, width),
	}
	
	for i := range instance.gMap {
		instance.gMap[i] = make([]*Tile, height)
	}
	
	instance.init()
	return instance
}

func (g *gameMap) init() {
	for a, column := range g.gMap {
		for b := range column {
			choice := rand.Intn(len(terrainTypes) - 1)
			g.gMap[a][b] = &Tile{terrainTypes[choice], 0, []*Player{}, a, b}
		}
	}
}

func (g gameMap) handleCombat() {
	var wg = sync.WaitGroup{}
	for x, _ := range g.gMap {
		for y, _ := range g.gMap[x] {
			wg.Add(1)
			go tileWorker(g.gMap[x][y], &wg)
		}
	}
	wg.Wait()
}

func (g gameMap) resources() {
	for x, _ := range g.gMap {
		for _, tile := range g.gMap[x] {
			tile.giveResources()
		}
	}
}

// getTile safely retrieves a tile, returning an Edge tile for out-of-bounds requests.
func (g gameMap) getTile(x, y int) *Tile {
	if x < 0 || x >= g.width || y < 0 || y >= g.height {
		return &Tile{Terrain: Edge, XPos: x, YPos: y} // Return a default edge tile
	}
	return g.gMap[x][y]
}

func (g gameMap) getTileFromPos(xPos int, yPos int) *Tile {
	return g.getTile(xPos, yPos)
}

func (g gameMap) getSurroundingsFromPos(xPos int, yPos int) Surroundings {
	return Surroundings{
		NW: g.getTile(xPos-1, yPos-1).getMapPiece(),
		NN: g.getTile(xPos, yPos-1).getMapPiece(),
		NE: g.getTile(xPos+1, yPos-1).getMapPiece(),
		WW: g.getTile(xPos-1, yPos).getMapPiece(),
		CE: g.getTile(xPos, yPos).getMapPiece(),
		EE: g.getTile(xPos+1, yPos).getMapPiece(),
		SW: g.getTile(xPos-1, yPos+1).getMapPiece(),
		SS: g.getTile(xPos, yPos+1).getMapPiece(),
		SE: g.getTile(xPos+1, yPos+1).getMapPiece(),
	}
}

func (g gameMap) spreadFromSpreader(xCoord int, yCoord int) {
	// TODO: decide if spread is 4 or 8 directions
	var xOffsets = []int{0, -1, 0, 1, 0}
	var yOffsets = []int{-1, 0, 0, 0, 1} //TODO: Check y-axis direction again!
	for neighbor := 0; neighbor < len(xOffsets); neighbor++ {
		var xTarget = xCoord + xOffsets[neighbor]
		var yTarget = yCoord + yOffsets[neighbor]
		if xTarget < 0 || xTarget >= g.width || yTarget < 0 || yTarget >= g.height {
			continue
		}
		g.gMap[xTarget][yTarget].spreadTo()
	}
}

func (g gameMap) fireAttractingTo(xPos int, yPos int) {
	var zombiesMoved = 0
	for x := -1; x <= 1; x++ {
		for y := -1; y <= 1; y++ {
			tile := g.getTile(xPos+x, yPos+y)
			if tile.Zombies > 0 {
				tile.removeZombies(1)
				zombiesMoved++
			}
		}
	}
	g.getTile(xPos, yPos).addZombies(zombiesMoved)
}

func (g gameMap) removeZombiesFromTile(xPos int, yPos int, count int) bool {
	return g.getTile(xPos, yPos).removeZombies(count)
}

func (g gameMap) addZombiesToTile(xPos int, yPos int, count int) {
	g.getTile(xPos, yPos).addZombies(count)
}

func (g *gameMap) spread() {
	for x, _ := range g.gMap {
		for y, tile := range g.gMap[x] {
			if tile.isSpreader() {
				g.spreadFromSpreader(x, y)
			}
		}
	}
}

func (g gameMap) getNewPlayerEntryTile() *Tile {
	var rX = r.Intn(g.width - 1)
	var rY = r.Intn(g.height - 1)
	return g.gMap[rX][rY]
}
