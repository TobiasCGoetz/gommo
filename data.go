package main

import (
	"fmt"
	"strings"
)

type IntTuple struct {
	x int `json:"x"`
	y int `json:"y"`
}

//Terrain
type Terrain int

const (
	Forest Terrain = iota
	Farm
	City
	Laboratory
)

var terrainTypes = [4]Terrain{ Forest, Farm, City, Laboratory }

func (t Terrain) toString() string {
	return []string{"Forest", "Farm", "City", "Laboratory"}[t]
}

type Direction int

const (
	North Direction = iota
	East
	South
	West
	Stay
)

var directions = [4]Direction{ North, East, South, West }

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
	terrain Terrain `json:"terrain"`
	zombies int `json:"zombies"`
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
