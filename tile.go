package main

import (
	"fmt"
	"sync"
)

type Tile struct {
	Terrain    Terrain
	Zombies    int
	playerPtrs []*Player
	XPos       int
	YPos       int
}

func tileWorker(t *Tile, wg *sync.WaitGroup) {
	defer wg.Done()
	t.resolveCombat()
	//fmt.Println("Worker started with ", t.toString())
}

func (t *Tile) resolveCombat() {
	totalPlayerStrength := 0
	// Iterate through the players on this tile
	for _, playerPtr := range t.playerPtrs {
		var player = *playerPtr
		var strength = 0
		weaponIndex, hasCard := hasCardWhere(player.Cards[:], Weapon)
		if player.Play == Weapon && hasCard { // Check if the played card is a weapon
			strength = weaponStrength
			player.Cards[weaponIndex] = None
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
		for _, playerPtr := range t.playerPtrs {
			playerPtr.Alive = false
			numDeadPlayers++
		}
		t.addZombies(numDeadPlayers)
	}
}

func (t Tile) giveResources() {
	for _, playerPtr := range t.playerPtrs {
		var player = *playerPtr
		cards, amount := t.Terrain.offersResource()
		for i := 0; i < amount; i++ {
			emptyIndex, hasSpace := hasCardWhere(player.Cards[:], None)
			if !hasSpace {
				continue
			}
			player.Cards[emptyIndex] = cards
		}
	}
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

func (t *Tile) addPlayer(playerPtr *Player) {
	t.playerPtrs = append(t.playerPtrs, playerPtr)
}

func (t *Tile) removePlayer(leavingPlayer *Player) {
	index, found := t.findPlayerPtrIndex(leavingPlayer)
	if found {
		t.playerPtrs = append(t.playerPtrs[:index], t.playerPtrs[index+1:]...)
	}
}

func (t Tile) findPlayerPtrIndex(requestedPlayerPtr *Player) (int, bool) {
	for index, playerPtr := range t.playerPtrs {
		if playerPtr == requestedPlayerPtr {
			return index, true
		}
	}
	return -1, false
}

func (t Tile) getMapPiece() MapPiece {
	var planNorth, planEast, planSouth, planWest = 0, 0, 0, 0
	for _, playerPtr := range t.playerPtrs {
		switch playerPtr.Direction {
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
		len(t.playerPtrs),
		planNorth,
		planEast,
		planSouth,
		planWest,
	}
}

func (t *Tile) addZombies(count int) {
	for i := 0; i < count; i++ {
		t.spreadTo()
	}
}

func (t *Tile) removeZombies(count int) bool {
	if count <= t.Zombies {
		t.Zombies -= count
		return true
	} else {
		t.Zombies = 0
		return false
	}
}

func (t Tile) toString() string {
	var r = fmt.Sprintf("%s %d %d", t.Terrain.toString(), t.Zombies, len(t.playerPtrs))
	return r
}
