package main

import "testing"

func TestGetPlayerOrNil(t *testing.T) {
	var playerX = 5
	var playerY = 5
	var testPlayer = Player{
		ID:        "test",
		X:         playerX,
		Y:         playerY,
		Direction: North,
		Play:      Dice,
		Consume:   Wood,
		Discard:   Wood,
		Cards:     [5]Card{Wood, Wood, Food, None, None},
		Alive:     true,
		IsBot:     true,
	}
	var testMap = make(map[string]*Player)
	testMap[testPlayer.ID] = &testPlayer

	playerFound := getPlayerOrNil(testMap, testPlayer.ID)
	if playerFound == nil {
		t.Errorf("Player in list not found")
	}
	playerNotFound := getPlayerOrNil(testMap, "ThisPlayerDoesntExist")
	if playerNotFound != nil {
		t.Errorf("Non-existant player found")
	}
}
