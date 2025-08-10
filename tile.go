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
	// Skip if no players on this tile
	if len(t.playerPtrs) == 0 {
		return
	}

	// Log combat start
	playerIDs := make([]string, 0, len(t.playerPtrs))
	for _, p := range t.playerPtrs {
		playerIDs = append(playerIDs, p.ID)
	}

	eventLogger.LogEvent(EventCombatStart, "", map[string]interface{}{
		"x":              t.XPos,
		"y":              t.YPos,
		"players":        playerIDs,
		"zombies_before": t.Zombies,
	})

	totalPlayerStrength := 0
	playerStrengths := make(map[string]int)

	// Calculate each player's strength
	for _, playerPtr := range t.playerPtrs {
		var player = *playerPtr
		var strength = 0
		weaponIndex, hasCard := hasCardWhere(player.Cards[:], Weapon)
		
		if player.Play == Weapon && hasCard {
			// Player used a weapon card
			strength = weaponStrength
			player.Cards[weaponIndex] = None
			
			eventLogger.LogEvent(EventCardUsed, player.ID, map[string]interface{}{
				"card":      Weapon.String(),
				"card_slot": weaponIndex,
				"x":         t.XPos,
				"y":         t.YPos,
				"strength":  strength,
			})
		} else {
			// Player rolls dice
			strength = rollDice(player.ID)
		}
		
		playerStrengths[player.ID] = strength
		totalPlayerStrength += strength
	}

	// Determine combat outcome
	combatWon := totalPlayerStrength > t.Zombies
	zombiesKilled := 0
	playersKilled := 0

	if combatWon {
		// Players win - kill all zombies
		zombiesKilled = t.Zombies
		t.Zombies = 0
	} else {
		// Zombies win - kill all players
		playersKilled = len(t.playerPtrs)
		for _, playerPtr := range t.playerPtrs {
			playerPtr.Alive = false
			
			// Log player death in combat
			eventLogger.LogEvent(EventPlayerDeath, playerPtr.ID, map[string]interface{}{
				"reason":    "combat",
				"x":         t.XPos,
				"y":         t.YPos,
				"zombies":   t.Zombies,
				"strength":  playerStrengths[playerPtr.ID],
			})
		}
		
		// Zombies multiply from dead players
		t.addZombies(playersKilled)
	}

	// Log combat result
	eventLogger.LogEvent(EventCombatResult, "", map[string]interface{}{
		"x":              t.XPos,
		"y":              t.YPos,
		"players":        playerIDs,
		"player_strength": totalPlayerStrength,
		"zombies_before": t.Zombies + zombiesKilled,
		"zombies_after":  t.Zombies,
		"combat_won":     combatWon,
		"zombies_killed": zombiesKilled,
		"players_killed": playersKilled,
	})
}

func (t Tile) giveResources() {
	for _, playerPtr := range t.playerPtrs {
		cards, amount := t.Terrain.offersResource()
		for i := 0; i < amount; i++ {
			emptyIndex, hasSpace := hasCardWhere(playerPtr.Cards[:], None)
			if !hasSpace {
				continue
			}
			playerPtr.Cards[emptyIndex] = cards
			// Track where research cards are acquired
			if cards == Research {
				playerPtr.ResearchAcquisitionPos[emptyIndex] = [2]int{t.XPos, t.YPos}
			}
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
