package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func setupAPI(playerList *[]*Player) {
	router := gin.Default()
	router.GET("/player/:id", getPlayerHandlerFunc(playerList))
	router.PUT("player/:id/direction/:dir", setDirectionHandlerFunc(playerList))
	router.PUT("player/:id/consume/:card", setConsumeHandlerFunc(playerList))
	router.PUT("player/:id/discard/:card", setDiscardHandlerFunc(playerList))
	router.PUT("player/:id/play/:card", setPlayHandlerFunc(playerList))
	router.Run("localhost:8080")
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

