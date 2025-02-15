package main

const mapWidth int = 100
const mapHeight int = 100
const botNumber int = 0
const zombieCutoff int = 13
const weaponStrength int = 6
const playerNameMaxLength int = 20
const playerMinAttack int = 1
const playerMaxAttack int = 6
const turnLength int = 15
const maxTurns int = 500
const victoryNumber int = 2

type terrainReward struct {
	amount    int
	givesCard Card
}

var terrainResources = map[Terrain]terrainReward{
	City:       terrainReward{1, Weapon},
	Forest:     terrainReward{2, Wood},
	Farm:       terrainReward{1, Food},
	Laboratory: terrainReward{1, Research},
}

const defaultDirection = South

var idSalt = "6LIBN8OWPzTKctUvbZtXV2mFn2tCq3qZKjHYbTTnLWtu6oGTU3ow3tuNx9SBTuND"
var hasWon bool = false
