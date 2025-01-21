package main

import (
	"strings"
	"fmt"
)

type IntTuple struct {
	x int
	y int
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
	terrain Terrain
	zombies int
}

//Player
type Player struct {
	id string
	x, y int
	direction Direction
	play Card
	consume Card
	discard Card
	cards [5]Card
	alive bool
}

func printPlayersList(s []Player) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

func (p Player) toString() string {
	var r strings.Builder
	r.WriteString(p.id)
	r.WriteString(": ")
	r.WriteString(fmt.Sprintf("%d", p.x))
	r.WriteString("|")
	r.WriteString(fmt.Sprintf("%d", p.y))
	r.WriteString(" ")
	r.WriteString(p.cards[0].toString())
	r.WriteString(p.cards[1].toString())
	r.WriteString(p.cards[2].toString())
	r.WriteString(p.cards[3].toString())
	return r.String()
}
