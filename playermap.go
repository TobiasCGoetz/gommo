package main

import "github.com/google/uuid"

// TODO: Not sure if relying on coordinates in the Player is a good idea
type playerMap struct {
	Players map[string]*Player
}

// TODO: Somehow remove inactive players
func (pm playerMap) addPlayer(playerName string, entryTile *Tile) string {
	playerID, _ := uuid.NewV7()
	idString := playerID.String()
	var player = Player{
		ID:                     idString,
		Name:                   playerName,
		CurrentTile:            entryTile,
		Direction:              defaultDirection,
		Play:                   None,
		Consume:                None,
		Discard:                None,
		Cards:                  [5]Card{Food, Wood, Wood, None, None},
		ResearchAcquisitionPos: [5][2]int{{-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}, {-1, -1}}, // Initialize with invalid positions
		Alive:                  true,
		IsBot:                  false,
	}
	pm.Players[idString] = &player
	entryTile.addPlayer(&player) // Actually add the player to the tile
	
	// Log player join event
	eventLogger.LogEvent(EventPlayerJoin, idString, map[string]interface{}{
		"name": playerName,
		"x":    entryTile.XPos,
		"y":    entryTile.YPos,
	})
	
	return idString
}

func (pm playerMap) move() {
	for _, player := range pm.Players {
		if !player.Alive {
			continue
		}

		oldTile := player.CurrentTile
		oldX, oldY := oldTile.XPos, oldTile.YPos

		if player.Direction == Stay {
			eventLogger.LogEvent(EventPlayerMove, player.ID, map[string]interface{}{
				"from_x": oldX,
				"from_y": oldY,
				"to_x":   oldX,
				"to_y":   oldY,
				"reason": "stayed",
			})
			continue
		}

		targetX, targetY := calculateNewPosition(oldX, oldY, player.Direction)
		clampedX, clampedY := clampToMapBoundaries(targetX, targetY)

		if clampedX == oldX && clampedY == oldY {
			eventLogger.LogEvent(EventPlayerMove, player.ID, map[string]interface{}{
				"from_x": oldX,
				"from_y": oldY,
				"to_x":   clampedX,
				"to_y":   clampedY,
				"reason": "blocked_by_boundary",
			})
			continue
		}

		handlePlayerMovement(player, oldTile, clampedX, clampedY)
	}
}

// calculateNewPosition computes the target coordinates based on direction.
func calculateNewPosition(x, y int, direction Direction) (int, int) {
	switch direction {
	case North:
		y++
	case East:
		x++
	case South:
		y--
	case West:
		x--
	}
	return x, y
}

// clampToMapBoundaries ensures coordinates are within the map limits.
func clampToMapBoundaries(x, y int) (int, int) {
	if x >= mapWidth {
		x = mapWidth - 1
	}
	if x < 0 {
		x = 0
	}
	if y >= mapHeight {
		y = mapHeight - 1
	}
	if y < 0 {
		y = 0
	}
	return x, y
}

// handlePlayerMovement updates the player's position and logs the event.
func handlePlayerMovement(player *Player, oldTile *Tile, newX, newY int) {
	newTile := gMap.getTileFromPos(newX, newY)
	oldTile.removePlayer(player)
	newTile.addPlayer(player)
	player.CurrentTile = newTile

	eventLogger.LogEvent(EventPlayerMove, player.ID, map[string]interface{}{
		"from_x": oldTile.XPos,
		"from_y": oldTile.YPos,
		"to_x":   newX,
		"to_y":   newY,
	})

	player.Direction = defaultDirection
}

func (p playerMap) playersConsume() {
	for _, playerPtr := range p.Players {
		playerPtr.consume()
	}
}

func (pm playerMap) limitCards() {
	for mapKey := range pm.Players {
		var player = pm.Players[mapKey]
		if player.getHandSize() > 4 { //TODO: Make configurable
			var cardPos, hasCard = hasCardWhere(player.Cards[:], player.Discard)
			if hasCard && player.Discard != None && cardPos > -1 { //Better safe...
				player.Cards[cardPos] = None
				player.ResearchAcquisitionPos[cardPos] = [2]int{-1, -1} // Clear research position
			} else {
				player.Cards[4] = None
				player.ResearchAcquisitionPos[4] = [2]int{-1, -1} // Clear research position
			}
		}
		player.Discard = None
		pm.Players[mapKey] = player
	}
}

func (pm playerMap) havePlayersWon() bool {
	for _, player := range pm.Players {
		if player.hasWinCondition() {
			return true
		}
	}
	return false
}

func (pm playerMap) getPlayer(id string) Player {
	return *pm.Players[id]
}

func (pm playerMap) getPlayerPtr(id string) *Player {
	return pm.Players[id]
}

func NewPlayerMap() playerMap {
	return playerMap{make(map[string]*Player)}
}

//TODO: Move bot stuff here?
