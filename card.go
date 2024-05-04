package main

import (
	"encoding/json"
)

type Card int

const (
	Food Card = iota
	Wood
	Weapon
	Dice
	Research
	None
)

var cardTypes = [6]Card{Food, Wood, Weapon, Dice, Research, None}

func (c Card) toString() string {
	return []string{"Food", "Wood", "Weapon", "Dice", "Research", "None"}[c]
}

func (c Card) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.toString())
}

var cards = map[string]Card{
	"food":     Food,
	"wood":     Wood,
	"weapon":   Weapon,
	"dice":     Dice,
	"research": Research,
	"none":     None,
}
