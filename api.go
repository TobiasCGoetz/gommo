package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func setupAPI(playerList *[]*Player, gameMap *[mapWidth][mapHeight]*Tile, turnTime *uint8) {
	router := gin.Default()
	router.GET("/player/:id", getPlayerHandlerFunc(playerList))
	router.POST("/player/:name", addPlayerHandlerFunc(playerList))
	router.PUT("/player/:id/direction/:dir", setDirectionHandlerFunc(playerList))
	router.PUT("/player/:id/consume/:card", setConsumeHandlerFunc(playerList))
	router.PUT("/player/:id/discard/:card", setDiscardHandlerFunc(playerList))
	router.PUT("/player/:id/play/:card", setPlayHandlerFunc(playerList))
	router.GET("/player/:id/surroundings", getSurroundingsHandlerFunc(playerList, gameMap))
	router.GET("/turnTimer", getRemainingTimerHandlerFunc(turnTime))
	router.Run("localhost:8080")
}

func addPlayerHandlerFunc (playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var pID = addPlayer(playerList, c.Param("name"))
		c.IndentedJSON(http.StatusOK, pID)
	}
	return fn
}

func getRemainingTimerHandlerFunc (turnTimer *uint8) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, *turnTimer)
	}
	return fn
}

//TODO: Add surrounding players info
func getSurroundingsHandlerFunc (playerList *[]*Player, gameMap *[mapWidth][mapHeight]*Tile) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		for _, player := range *playerList {
			if player.ID == id {
				var NW = Tile{Edge, -1, []Player{}}
				var NN = Tile{Edge, -1, []Player{}}
				var NE = Tile{Edge, -1, []Player{}}
				var WW = Tile{Edge, -1, []Player{}}
				var CE = *gameMap[player.X][player.Y]
				var EE = Tile{Edge, -1, []Player{}}
				var SW = Tile{Edge, -1, []Player{}}
				var SS = Tile{Edge, -1, []Player{}}
				var SE = Tile{Edge, -1, []Player{}}

				if player.X > 0 && player.Y < mapWidth-1 {
					NW = *gameMap[player.X-1][player.Y+1]
				}
				if player.X < mapWidth-1 && player.Y < mapHeight-1 {
					NN = *gameMap[player.X][player.Y+1]
					NE = *gameMap[player.X+1][player.Y+1]
					EE = *gameMap[player.X+1][player.Y]
				}
				if player.X < mapWidth-1 && player.Y > 0 {
					SE = *gameMap[player.X+1][player.Y-1]
				}
				if player.X > 0 && player.Y > 0 {
					WW = *gameMap[player.X-1][player.Y]
					SW = *gameMap[player.X-1][player.Y-1]
					SS = *gameMap[player.X][player.Y-1]
				}

				//TODO: Add and check password phrase
				var miniMap = Surroundings{
					NW: NW,
					NN: NN,
					NE: NE,
					WW: WW,
					CE: CE,
					EE: EE,
					SW: SW,
					SS: SS,
					SE: SE,
				}
				c.IndentedJSON(http.StatusOK, miniMap)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Player not found."})
	}
	return fn
}

func setDiscardHandlerFunc (playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card, err = strconv.Atoi(cardStr)
		if err != nil || card >= len(cardTypes)  {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No number recognized"})
		}
		for pNr, player := range *playerList {
			if player.ID == id {
				//TODO: Add and check password phrase
				(*playerList)[pNr].Discard = cardTypes[card]
				c.IndentedJSON(http.StatusOK, *player)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Player not found."})
	}
	return fn
}
func setPlayHandlerFunc (playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card, err = strconv.Atoi(cardStr)
		if err != nil || card >= len(cardTypes)  {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No number recognized"})
		}
		for pNr, player := range *playerList {
			if player.ID == id {
				//TODO: Add and check password phrase
				(*playerList)[pNr].Play = cardTypes[card]
				c.IndentedJSON(http.StatusOK, *player)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Player not found."})
	}
	return fn
}

func setConsumeHandlerFunc (playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card, err = strconv.Atoi(cardStr)
		if err != nil || card >= len(cardTypes)  {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No number recognized"})
		}
		for pNr, player := range *playerList {
			if player.ID == id {
				//TODO: Add and check password phrase
				(*playerList)[pNr].Consume = cardTypes[card]
				c.IndentedJSON(http.StatusOK, *player)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Player not found."})
	}
	return fn
}

func setDirectionHandlerFunc (playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		dirStr := c.Param("dir")
		var dir, err = strconv.Atoi(dirStr)
		if err != nil || dir >= len(Directions)  {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No number recognized"})
		}
		for pNr, player := range *playerList {
			if player.ID == id {
				//TODO: Add and check password phrase
				(*playerList)[pNr].Direction = Directions[dir]
				c.IndentedJSON(http.StatusOK, *player)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Player not found."})
	}
	return fn
}

func getPlayerHandlerFunc (playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		//passwd := c.Param("passwd")
		//TODO: DO NOT SEARCH HERE!
		for _, player := range *playerList {
			if player.ID == id {
				//TODO: Add and check password phrase
				c.IndentedJSON(http.StatusOK, *player)
				return
			}
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Player not found."})
	}
	return fn
}

func sanitizeSurroundingsInfo (surroundings *Surroundings) {
	maskPlayerInfo(&surroundings.NW)
	maskPlayerInfo(&surroundings.NN)
	maskPlayerInfo(&surroundings.NE)
	maskPlayerInfo(&surroundings.WW)
	maskPlayerInfo(&surroundings.CE)
	maskPlayerInfo(&surroundings.EE)
	maskPlayerInfo(&surroundings.SW)
	maskPlayerInfo(&surroundings.SS)
	maskPlayerInfo(&surroundings.SE)
}

func maskPlayerInfo (tile *Tile) {
	for playerNr, player := range tile.Players {
		//Filter the dead
		if player.Alive {
			tile.Players = append(tile.Players[:playerNr], tile.Players[playerNr+1:]...)
			continue
		}
		//Come up with useful playernames
		tile.Players[playerNr].ID = "Spielername"
		//Blank out hidden info
		tile.Players[playerNr].Cards = [5]Card{}
		tile.Players[playerNr].Consume = None
		tile.Players[playerNr].Play = None
		tile.Players[playerNr].IsBot = true
	}
}