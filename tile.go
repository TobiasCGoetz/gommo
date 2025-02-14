package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type Tile struct {
	Terrain   Terrain
	Zombies   int
	playerIds []string
}

func (t *Tile) resolveCombat() {
	totalPlayerStrength := 0
	// Iterate through the players on this tile
	for _, id := range t.playerIds {
		var player = getPlayerOrNil(id)
		if player == nil || !player.Alive {
			continue
		}
		var strength = 0
		if player.Play == Weapon { // Check if the played card is a weapon
			strength = weaponStrength
		} else {
			strength = rollDice() // Alternatively, roll a dice
		}
		totalPlayerStrength += strength
	}

	if totalPlayerStrength > t.Zombies {
		t.Zombies = 0
	} else {
		// Kill all players on the tile
		numDeadPlayers := 0
		for _, id := range t.playerIds {
			var player = getPlayerOrNil(id)
			if player == nil {
				continue
			}
			if player.Alive {
				player.Alive = false
				numDeadPlayers++
			}
		}
		t.Zombies += numDeadPlayers // Add killed player count to zombies
	}
}

func rollDice() int {
	return rand.Intn(6) + 1 // rand.Intn(6) generates 0-5, so we add 1
}

func (t Tile) isSpreader() bool {
	return t.Terrain.isCity() || t.Zombies >= zombieCutoff
}

func (t *Tile) spreadTo() {
	if t.Zombies < zombieCutoff {
		t.Zombies++
	}
}

func (t *Tile) spreadToUnbound() {
	t.Zombies++
}

func (t *Tile) addPlayer(incomingPlayer string) {
	t.playerIds = append(t.playerIds, incomingPlayer)
}

func (t *Tile) removePlayer(leavingPlayer string) {
	index, found := t.findPlayerIdIndex(leavingPlayer)
	if found {
		t.playerIds = append(t.playerIds[:index], t.playerIds[index+1])
	}
}

func (t Tile) findPlayerIdIndex(playerId string) (int, bool) {
	for index, pId := range t.playerIds {
		if pId == playerId {
			return index, true
		}
	}
	return -1, false
}

func (t Tile) getMapPiece() MapPiece {
	var planNorth, planEast, planSouth, planWest = 0, 0, 0, 0
	for _, pId := range t.playerIds {
		switch getPlayerOrNil(pId).Direction {
		case North:
			planNorth++
		case East:
			planEast++
		case South:
			planSouth++
		case West:
			planWest++
		}
	}
	return MapPiece{
		t.Terrain.toString(),
		t.Zombies,
		len(t.playerIds),
		planNorth,
		planEast,
		planSouth,
		planWest,
	}
}

func (t Tile) toString() string {
	var r = fmt.Sprintf("%s %d %s", t.Terrain.toString(), t.Zombies, strings.Join(t.playerIds, ","))
	return r
}
