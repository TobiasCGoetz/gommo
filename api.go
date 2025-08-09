package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupAPI() {
	router := gin.Default()
	router.Use(cors.Default())

	// Add middleware for error handling and logging
	router.Use(errorHandlingMiddleware())

	router.GET("/player/:id", getPlayerHandlerFunc())
	router.GET("/player/:id/surroundings", getSurroundingsHandlerFunc())
	router.GET("/config", getAllConfigHandlerFunc())
	router.POST("/player/:name", addPlayerHandlerFunc())
	router.PUT("/player/:id/direction/:dir", setDirectionHandlerFunc())
	router.PUT("/player/:id/play/:cardType", setPlayHandlerFunc())
	router.Run("0.0.0.0:8080")
}

func getAllConfigHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(http.StatusOK, gState)
	}
	return fn
}

func addPlayerHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var pId = pMap.addPlayer(
			filterPlayerName(c.Param("name")),
			gMap.getNewPlayerEntryTile())
		c.JSON(http.StatusOK, pId)
	}
	return fn
}

func getSurroundingsHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		var player = pMap.getPlayer(id)
		var xPos = player.CurrentTile.XPos
		var yPos = player.CurrentTile.YPos
		c.JSON(http.StatusOK, gMap.getSurroundingsFromPos(xPos, yPos))
		return
	}
	return fn
}

func setPlayHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		playerPtr := getPlayerOrNil(id)
		if playerPtr != nil {
			playerPtr.cardInput(cardStr)
			c.Status(http.StatusOK)
			return
		} else {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
	return fn
}

func setDirectionHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		dirStr := c.Param("dir")
		var dir = directions[strings.ToLower(dirStr)]
		playerPtr := getPlayerOrNil(id)
		if playerPtr != nil {
			(*playerPtr).Direction = dir
			c.Status(http.StatusOK)
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
	return fn
}

func getPlayerHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		playerPtr := getPlayerOrNil(id)
		if playerPtr != nil {
			c.JSON(http.StatusOK, (*playerPtr))
			return
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
	return fn
}

func filterPlayerName(name string) string {
	if len(name) > playerNameMaxLength {
		return name[0:playerNameMaxLength]
	}
	//TODO: Filter bad words
	return name
}

// errorHandlingMiddleware provides centralized error handling and logging
func errorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in API handler: %v", r)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
			}
		}()

		c.Next()
	}
}
