package main

type Tile struct {
	Terrain   Terrain
	Zombies   int
	playerIds []string
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
