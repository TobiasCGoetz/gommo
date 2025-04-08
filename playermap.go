package main

import "github.com/google/uuid"

// TODO: Not sure if relying on coordinates in the Player is a good idea
type playerMap struct {
	Players map[string]*Player
}

// TODO: Somehow remove inactive players
func (pm playerMap) addPlayer(playerName string) string {
	var rX = r.Intn(mapWidth - 1)
	var rY = r.Intn(mapHeight - 1)
	playerID, _ := uuid.NewV7()
	idString := playerID.String()
	var player = Player{
		ID:        idString,
		Name:      playerName,
		X:         rX,
		Y:         rY,
		Direction: defaultDirection,
		Play:      None,
		Consume:   None,
		Discard:   None,
		Cards:     [5]Card{Food, Wood, Wood, None, None},
		Alive:     true,
		IsBot:     false,
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

		//Perform move
		switch player.Direction {
		case North:
			player.Y += 1
		case East:
			player.X += 1
		case South:
			player.Y -= 1
		case West:
			player.X -= 1
		case Stay:
			return
		}

		//Write new coordinates
		pm.Players[player.ID] = player

		//Prevent out-of-map moves
		if mapWidth <= player.X {
			player.X = mapWidth - 1
		}
		if player.X < 0 {
			player.X = 0
		}
		if mapHeight <= player.Y {
			player.Y = mapHeight - 1
		}
		if player.Y < 0 {
			player.Y = 0
		}
		//Reset move direction
		player.Direction = defaultDirection
		//Write new state
		pm.Players[player.ID] = player
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
