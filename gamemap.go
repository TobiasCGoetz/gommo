package main

import (
	"math/rand"
	"sync"
)

type gameMap struct {
	gMap [mapWidth][mapHeight]*Tile
}

func NewGameMap() gameMap {
	instance := gameMap{}
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

func (g gameMap) getTileFromPos(xPos int, yPos int) *Tile {
	return g.gMap[xPos][yPos]
}

func (g gameMap) getSurroundingsFromPos(xPos int, yPos int) Surroundings {
	return Surroundings{
		NW: g.gMap[xPos-1][yPos-1].getMapPiece(),
		NN: g.gMap[xPos][yPos-1].getMapPiece(),
		NE: g.gMap[xPos+1][yPos-1].getMapPiece(),
		WW: g.gMap[xPos-1][yPos].getMapPiece(),
		CE: g.gMap[xPos][yPos].getMapPiece(),
		EE: g.gMap[xPos+1][yPos].getMapPiece(),
		SW: g.gMap[xPos-1][yPos+1].getMapPiece(),
		SS: g.gMap[xPos][yPos+1].getMapPiece(),
		SE: g.gMap[xPos+1][yPos+1].getMapPiece(),
	}
}

func (g gameMap) spreadFromSpreader(xCoord int, yCoord int) {
	// TODO: decide if spread is 4 or 8 directions
	var xOffsets = []int{0, -1, 0, 1, 0}
	var yOffsets = []int{-1, 0, 0, 0, 1} //TODO: Check y-axis direction again!
	for neighbor := 0; neighbor < len(xOffsets); neighbor++ {
		var xTarget = xCoord + xOffsets[neighbor]
		var yTarget = yCoord + yOffsets[neighbor]
		if xTarget < 0 || xTarget >= mapWidth || yTarget < 0 || yTarget >= mapHeight {
			continue
		}
		g.gMap[xTarget][yTarget].spreadTo()
	}
}

func (g gameMap) removeZombiesFromTile(xPos int, yPos int, count int) bool {
	return g.gMap[xPos][yPos].removeZombies(count)
}

func (g gameMap) addZombiesToTile(xPos int, yPos int, count int) {
	g.gMap[xPos][yPos].addZombies(count)
}

func (g gameMap) consume(playerMap *map[string]*Player) { //TODO: Simplify this a lot
	for _, playerPtr := range *playerMap {
		//Fetch current player state
		if !playerPtr.Alive {
			continue
		}
		var playerX = playerPtr.CurrentTile.XPos
		var playerY = playerPtr.CurrentTile.YPos
		//We don't allow death by indecision
		if playerPtr.Consume == None {
			_, hasCard := hasCardWhere(playerPtr.Cards[:], Food)
			if hasCard {
				playerPtr.Consume = Food
			} else {
				playerPtr.Consume = Wood
			}
		}

		//Now remove that card or kill the player
		cardPos, hasCard := hasCardWhere(playerPtr.Cards[:], playerPtr.Consume)
		if hasCard {
			if playerPtr.Consume == Wood {

				var zombiesAttracted = 0

				//Remove zombies from surrounding tiles
				for i := -1; i <= 1; i++ {
					for j := -1; j <= 1; j++ {
						if g.removeZombiesFromTile(playerX+i, playerY+j, 1) {
							zombiesAttracted++
						}
					}
				}
				//Add to players tile
				g.addZombiesToTile(playerX, playerY, zombiesAttracted)
			}
			playerPtr.Cards[cardPos] = None //Remove card from hand
		} else {
			playerPtr.Alive = false //Card not in hand, kill the player
		}
	}
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
	var rX = r.Intn(mapWidth - 1)
	var rY = r.Intn(mapHeight - 1)
	return g.gMap[rX][rY]
}
