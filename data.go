package main

import (
	"strings"
	"fmt"
)

//Tile
type Tile int

const (
	Forest Tile = iota
	Farm
	City
	Laboratory
)

var tileTypes = [4]Tile{ Forest, Farm, City, Laboratory }

func (t Tile) toString() string {
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

func (d Direction) toString() string {
	return []string{"North", "East", "South", "West", "Stay"}[d]
}

//Cards
type Card int

const (
	Food Card = iota
	Wood
	Weapon
	None
)


func (c Card) toString() string {
	return []string{"Food", "Wood", "Weapon", "None"}[c]
}

//Players

type Player struct {
	id string
	x, y int
	dir Direction
	play Card
	cards [5]Card
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
