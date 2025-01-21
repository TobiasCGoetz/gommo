package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func setupAPI(playerList *[]*Player, gameMap *[mapWidth][mapHeight]*Tile, turnTime *uint8) {
	router := gin.Default()
	router.GET("/player/:id", getPlayerHandlerFunc(playerList))
	router.PUT("/player/:id/direction/:dir", setDirectionHandlerFunc(playerList))
	router.PUT("/player/:id/consume/:card", setConsumeHandlerFunc(playerList))
	router.PUT("/player/:id/discard/:card", setDiscardHandlerFunc(playerList))
	router.PUT("/player/:id/play/:card", setPlayHandlerFunc(playerList))
	router.GET("/player/:id/surroundings", getSurroundingsHandlerFunc(playerList, gameMap))
	router.GET("/turnTimer", getRemainingTimerHandlerFunc(turnTime))
	router.Run("localhost:8080")
}

func getRemainingTimerHandlerFunc (turnTimer *uint8) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, *turnTimer)
	}
	return fn
}

//TODO: Add surrounding players info, fix out-of-gamemap access
func getSurroundingsHandlerFunc (playerList *[]*Player, gameMap *[mapWidth][mapHeight]*Tile) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		for _, player := range *playerList {
			if player.ID == id {
				//TODO: Add and check password phrase
				var miniMap = Surroundings{
					NW: *gameMap[player.X-1][player.Y+1],
					NN: *gameMap[player.X][player.Y+1],
					NE: *gameMap[player.X+1][player.Y+1],
					WW: *gameMap[player.X-1][player.Y],
					CE: *gameMap[player.X][player.Y],
					EE: *gameMap[player.X+1][player.Y],
					SW: *gameMap[player.X-1][player.Y-1],
					SS: *gameMap[player.X][player.Y-1],
					SE: *gameMap[player.X+1][player.Y-1],
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

