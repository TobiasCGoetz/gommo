package main

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
