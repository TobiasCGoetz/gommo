package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupAPI() {
	router := gin.Default()
	router.Use(cors.Default())
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
		c.JSON(http.StatusOK, registry.Dispatch(GetConfigEvent{}))
	}
	return fn
}

func getConfigGameStateHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(http.StatusOK, hasWon)
	}
	return fn
}

func getConfigTurnTimerHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(http.StatusOK, turnLength)
	}
	return fn
}

func getConfigMapSizeHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(http.StatusOK, IntTuple{mapWidth, mapHeight})
	}
	return fn
}

func addPlayerHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var pName = filterPlayerName(c.Param("name"))
		c.JSON(http.StatusOK, registry.Dispatch(CreateUserEvent{BaseEvent{}, pName}))
	}
	return fn
}

func getSurroundingsHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		surroundings := registry.Dispatch(NewGetSurroundingsEvent(id))
		c.JSON(http.StatusOK, surroundings)
		return
	}
	return fn
}

func setDiscardHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card, err = strconv.Atoi(cardStr)
		if err != nil || card >= len(cardTypes) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		playerPtr := getPlayerOrNil(id)
		if playerPtr != nil {
			(playerPtr).Discard = cardTypes[card]
			c.Status(http.StatusOK)
			return
		} else {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
	return fn
}

func setPlayHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card, err = strconv.Atoi(cardStr)
		if err != nil || card >= len(cardTypes) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		playerPtr := getPlayerOrNil(id)
		if playerPtr != nil {
			(*playerPtr).Play = cardTypes[card]
			c.Status(http.StatusOK)
			return
		} else {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
	}
	return fn
}

func setConsumeHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card = cards[strings.ToLower(cardStr)]
		playerPtr := getPlayerOrNil(id)
		if playerPtr != nil {
			(*playerPtr).Consume = card
			c.Status(http.StatusOK)
			return
		} else {
			c.AbortWithStatus(http.StatusForbidden)
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
