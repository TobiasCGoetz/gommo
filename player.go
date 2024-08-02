package main

import (
	"fmt"
	"strings"
)

type Player struct {
	ID        string
	Name      string
	X, Y      int
	Direction Direction
	Play      Card
	Consume   Card
	Discard   Card
	Cards     [5]Card
	Alive     bool
	IsBot     bool
}

func (p Player) hasWinCondition() bool {
	var numberOfResearchs = 0
	for _, card := range p.Cards {
		if card == Research {
			numberOfResearchs++
		}
	}
	if numberOfResearchs < victoryNumber {
		return false
	} else {
		fmt.Println("Player has won")
		fmt.Println(p.toString())
		return true
	}
}

func (p Player) toString() string {
	var r strings.Builder
	r.WriteString(p.ID)
	r.WriteString(p.Name)
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

func printPlayersList(s []Player) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
