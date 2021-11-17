package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func setupAPI(playerList *[]*Player) {
	router := gin.Default()
	router.GET("/player/:id", getPlayerHandlerFunc(playerList))
	router.Run("localhost:8080")
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

