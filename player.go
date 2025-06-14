package main

import (
	"fmt"
	"strings"
)

type Player struct {
	ID          string
	Name        string
	X, Y        int
	CurrentTile *Tile
	Direction   Direction
	Play        Card
	Consume     Card
	Discard     Card
	Cards       [5]Card
	Alive       bool
	IsBot       bool
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

func hasCardWhere(ar []Card, card Card) (int, bool) {
	for a, c := range ar {
		if c == card {
			return a, true
		}
	}
	return -1, false
}

func (p Player) firstIndexOfCardType(target Card) int {
	for i, card := range p.Cards {
		if card == target {
			return i
		}
	}
	return -1
}

func (p Player) getHandSize() int { //TODO: Move to method
	var count = 0
	for _, card := range p.Cards {
		if card != None {
			count++
		}
	}
	return count
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
