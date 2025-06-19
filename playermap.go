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
		ID:          idString,
		Name:        playerName,
		CurrentTile: entryTile,
		Direction:   defaultDirection,
		Play:        None,
		Consume:     None,
		Discard:     None,
		Cards:       [5]Card{Food, Wood, Wood, None, None},
		Alive:       true,
		IsBot:       false,
	}
	pm.Players[idString] = &player
	return idString
}

func (pm playerMap) move() {
	//Set new coordinates per player from move
	for _, player := range pm.Players {
		if !player.Alive {
			continue
		}
		//Fetch current player state
		var player = pm.Players[player.ID]
		var playerX = player.CurrentTile.XPos
		var playerY = player.CurrentTile.YPos
		var targetX = playerX
		var targetY = playerY

		//Perform move
		switch player.Direction {
		case North:
			targetY += 1
		case East:
			targetX += 1
		case South:
			targetY -= 1
		case West:
			targetX -= 1
		case Stay:
			return
		}

		//Prevent out-of-map moves
		if targetX >= mapWidth {
			targetX = mapWidth - 1
		}
		if targetX < 0 {
			targetX = 0
		}
		if targetY >= mapHeight {
			targetY = mapHeight - 1
		}
		if targetY < 0 {
			targetY = 0
		}

		// Remove player from current position
		// Add player to new position
		// Set player.CurrentTile
		var oldTile = player.CurrentTile
		var newTile = gMap.getTileFromPos(targetX, targetY)
		oldTile.removePlayer(player)
		newTile.addPlayer(player)
		player.CurrentTile = newTile

		//Reset move direction
		player.Direction = defaultDirection
	}
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
			} else {
				player.Cards[4] = None
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

func NewPlayerMap() playerMap {
	return playerMap{make(map[string]*Player)}
}

//TODO: Move bot stuff here?
