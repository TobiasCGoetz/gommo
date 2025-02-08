package main

type Tile struct {
	Terrain   Terrain
	Zombies   int
	playerIds []string
}

func (t Tile) isCity() bool {
	if t.Terrain == City {
		return true
	}
	return false
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

type MapPiece struct {
	TileType             string
	ZombieCount          int
	PlayerCount          int
	PlayersPlanMoveNorth int
	PlayersPlanMoveEast  int
	PlayersPlanMoveSouth int
	PlayersPlanMoveWest  int
}

type Surroundings struct {
	NW MapPiece
	NN MapPiece
	NE MapPiece
	WW MapPiece
	CE MapPiece
	EE MapPiece
	SW MapPiece
	SS MapPiece
	SE MapPiece
}
