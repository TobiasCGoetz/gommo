package main

import (
	"fmt"
	"strings"
)

type IntTuple struct {
	X int
	Y int
}

//Terrain
type Terrain int

const (
	Forest Terrain = iota
	Farm
	City
	Laboratory
	Edge
)

var terrainTypes = [5]Terrain{ Forest, Farm, City, Laboratory, Edge }

func (t Terrain) toString() string {
	return []string{"Forest", "Farm", "City", "Laboratory", "Edge"}[t]
}

type Direction int

const (
	North Direction = iota
	East
	South
	West
	Stay
)

var Directions = [4]Direction{North, East, South, West }

func (d Direction) toString() string {
	return []string{"North", "East", "South", "West", "Stay"}[d]
}

//Cards
type Card int

const (
	Food Card = iota
	Wood
	Weapon
	Dice
	Research
	None
)

var cardTypes = [6]Card{ Food, Wood, Weapon, Dice, Research, None }

func (c Card) toString() string {
	return []string{"Food", "Wood", "Weapon", "Dice", "Research", "None"}[c]
}

//Tile
type Tile struct {
	Terrain Terrain
	Zombies int
}

type Surroundings struct {
	NW Tile
	NN Tile
	NE Tile
	WW Tile
	CE Tile
	EE Tile
	SW Tile
	SS Tile
	SE Tile
}

//Player
type Player struct {
	ID string
	X, Y int
	Direction Direction
	Play Card
	Consume Card
	Discard Card
	Cards [5]Card
	Alive bool
	IsBot bool
}

func printPlayersList(s []Player) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func (p Player) toString() string {
	var r strings.Builder
	r.WriteString(p.ID)
	r.WriteString(": ")
	r.WriteString(fmt.Sprintf("%d", p.X))
	r.WriteString("|")
	r.WriteString(fmt.Sprintf("%d", p.Y))
	r.WriteString(" ")
	r.WriteString(p.Cards[0].toString())
	r.WriteString(p.Cards[1].toString())
	r.WriteString(p.Cards[2].toString())
	r.WriteString(p.Cards[3].toString())
	return r.String()
}
