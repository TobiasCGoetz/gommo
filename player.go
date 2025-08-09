package main

import (
	"fmt"
	"strings"
)

type Player struct {
	ID                    string
	Name                  string
	CurrentTile           *Tile
	Direction             Direction
	Play                  Card
	Consume               Card
	Discard               Card
	Cards                 [5]Card
	ResearchAcquisitionPos [5][2]int // Track x,y coordinates where each research card was acquired
	Alive                 bool
	IsBot                 bool
}

func (p *Player) consume() {
	if !p.Alive {
		return
	}
	var playerX = p.CurrentTile.XPos
	var playerY = p.CurrentTile.YPos
	//We don't allow death by indecision
	if p.Consume == None {
		_, hasCard := hasCardWhere(p.Cards[:], Food)
		if hasCard {
			p.Consume = Food
		} else {
			p.Consume = Wood
		}
	}

	//Now remove that card or kill the player
	cardPos, hasCard := hasCardWhere(p.Cards[:], p.Consume)
	if hasCard {
		if p.Consume == Wood {
			gMap.fireAttractingTo(playerX, playerY)
		}
		p.Cards[cardPos] = None //Remove card from hand
		// Clear research acquisition position when card is removed
		p.ResearchAcquisitionPos[cardPos] = [2]int{-1, -1}
	} else {
		p.Alive = false //Card not in hand, kill the player
	}
}

func (p *Player) cardInput(inputCard string) {
	if inputCard == Weapon.toString() {
		p.Play = Weapon
	} else {
		p.Consume = cards[inputCard]
	}
}

func (p Player) hasWinCondition() bool {
	// Must be at a laboratory to win
	if p.CurrentTile.Terrain != Laboratory {
		return false
	}
	
	var numberOfResearchs = 0
	currentX := p.CurrentTile.XPos
	currentY := p.CurrentTile.YPos
	
	for i, card := range p.Cards {
		if card == Research {
			// Check if this research card was acquired at the current laboratory
			acquisitionX := p.ResearchAcquisitionPos[i][0]
			acquisitionY := p.ResearchAcquisitionPos[i][1]
			
			// If research was acquired at current location, it doesn't count for victory
			if acquisitionX == currentX && acquisitionY == currentY {
				continue
			}
			
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
	//r.WriteString(fmt.Sprintf("%d", p.X))
	r.WriteString("|")
	//r.WriteString(fmt.Sprintf("%d", p.Y))
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
