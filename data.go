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

func (t Tile) toString() string {
	return []string{"Forest", "Farm", "City", "Laboratory"}[t]
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
	cards [4]Card
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
