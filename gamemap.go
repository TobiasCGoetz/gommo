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
			g.gMap[a][b] = &Tile{terrainTypes[choice], 0, []string{}}
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

func (g gameMap) getSurroundingsOfPlayer(id string) (Surroundings, bool) {
	player := getPlayerOrNil(id)
	if player == nil { //TODO: If nil else function or invert? Make them all identical!
		return Surroundings{}, false
	} else {
		var NW = g.gMap[player.X-1][player.Y-1].getMapPiece()
		var NN = g.gMap[player.X][player.Y-1].getMapPiece()
		var NE = g.gMap[player.X+1][player.Y-1].getMapPiece()
		var WW = g.gMap[player.X-1][player.Y].getMapPiece()
		var CE = g.gMap[player.X][player.Y].getMapPiece()
		var EE = g.gMap[player.X+1][player.Y].getMapPiece()
		var SW = g.gMap[player.X-1][player.Y+1].getMapPiece()
		var SS = g.gMap[player.X][player.Y+1].getMapPiece()
		var SE = g.gMap[player.X+1][player.Y+1].getMapPiece()

		var miniMap = Surroundings{
			NW: NW,
			NN: NN,
			NE: NE,
			WW: WW,
			CE: CE,
			EE: EE,
			SW: SW,
			SS: SS,
			SE: SE,
		}
		return miniMap, true
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

func (g gameMap) consume(playerMap *map[string]*Player) {
	for playerID := range *playerMap {
		//Fetch current player state
		var player = (*playerMap)[playerID]
		if !player.Alive {
			continue
		}

		//We don't allow death by indecision
		if player.Consume == None {
			_, hasCard := hasCardWhere(player.Cards[:], Food)
			if hasCard {
				player.Consume = Food
			} else {
				player.Consume = Wood
			}
		}

		//Now remove that card or kill the player
		cardPos, hasCard := hasCardWhere(player.Cards[:], player.Consume)
		if hasCard {
			if player.Consume == Wood {

				var zombiesAttracted = 0

				var tileNW = getMapTile(player.X-1, player.Y+1, &g.gMap)
				var tileNN = getMapTile(player.X, player.Y+1, &g.gMap)
				var tileNE = getMapTile(player.X+1, player.Y+1, &g.gMap)
				var tileWW = getMapTile(player.X-1, player.Y, &g.gMap)
				var tileEE = getMapTile(player.X+1, player.Y, &g.gMap)
				var tileSW = getMapTile(player.X-1, player.Y-1, &g.gMap)
				var tileSS = getMapTile(player.X, player.Y-1, &g.gMap)
				var tileSE = getMapTile(player.X+1, player.Y-1, &g.gMap)

				var tileArray = []*Tile{tileNW, tileNN, tileNE, tileWW, tileEE, tileSW, tileSS, tileSE}

				//Remove zombies from surrounding tiles
				for _, nextTile := range tileArray {
					if nextTile.Zombies > 0 {
						zombiesAttracted++
						nextTile.Zombies--
					}
				}

				//Add to players tile
				getMapTile(player.X, player.Y, &g.gMap).Zombies += zombiesAttracted
			}
			player.Cards[cardPos] = None //Remove card from hand
			(*playerMap)[player.ID] = player
		} else {
			player.Alive = false //Card not in hand, kill the player
			(*playerMap)[player.ID] = player
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
