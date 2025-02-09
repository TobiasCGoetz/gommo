package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// TODO: Move functionality&complexity outside of this api file, call suitable functions instead
// Ideally, we wouldn't rely on the data here at all
func setupAPI(gameMap *[mapWidth][mapHeight]*Tile, turnTime *int8, hasWon *bool) {
	router := gin.Default()
	router.Use(cors.Default())
	//Player endpoints
	router.POST("/player/:name", addPlayerHandlerFunc())
	router.GET("/player/:id", getPlayerHandlerFunc())
	router.GET("/player/:id/surroundings", getSurroundingsHandlerFunc())
	router.PUT("/player/:id/direction/:dir", setDirectionHandlerFunc())
	router.PUT("/player/:id/consume/:card", setConsumeHandlerFunc())
	router.PUT("/player/:id/discard/:card", setDiscardHandlerFunc())
	router.PUT("/player/:id/play/:card", setPlayHandlerFunc())
	//Config endpoints
	router.GET("/config/turnTimer", getRemainingTimerHandlerFunc(*turnTime)) //TODO: Call function instead of relying on argument
	router.GET("/config/turnLength", getConfigTurnTimerHandlerFunc())
	router.GET("/config/mapSize", getConfigMapSizeHandlerFunc())
	router.GET("/config/hasWon", getConfigGameStateHandlerFunc(*hasWon)) //TODO: Call function instead of relying on argument
	router.GET("/config", getAllConfigHandlerFunc(*turnTime, *hasWon))   //TODO: Call function instead of relying on argument
	router.Run("0.0.0.0:8080")
}

func getAllConfigHandlerFunc(turnTimer int8, hasWon bool) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var response = ConfigResponse{int(turnTimer), turnLength, hasWon}
		c.JSON(http.StatusOK, response)
	}
	return fn
}

func getConfigGameStateHandlerFunc(hasWon bool) gin.HandlerFunc {
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
		var pID = addPlayer(pName)
		c.JSON(http.StatusOK, pID)
	}
	return fn
}

func getRemainingTimerHandlerFunc(turnTimer int8) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.JSON(http.StatusOK, turnTimer)
	}
	return fn
}

func getSurroundingsHandlerFunc() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		miniMap, success := getSurroundingsOfPlayer(id)
		if !success {
			c.AbortWithStatus(http.StatusBadRequest)
		}
		c.JSON(http.StatusOK, miniMap)
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
