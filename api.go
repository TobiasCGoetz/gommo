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
func setupAPI(playerMap *map[string]*Player, gameMap *[mapWidth][mapHeight]*Tile, turnTime *int8, hasWon *bool) {
	router := gin.Default()
	router.Use(cors.Default())
	//GET endpoints only receive call-by-value arguments
	//POST/PUT endpoints receive a pointer to enable writes
	//Player endpoints
	router.POST("/player/:name", addPlayerHandlerFunc(playerMap))
	router.GET("/player/:id", getPlayerHandlerFunc(*playerMap))
	router.GET("/player/:id/surroundings", getSurroundingsHandlerFunc(*playerMap, *gameMap))
	router.PUT("/player/:id/direction/:dir", setDirectionHandlerFunc(playerMap))
	router.PUT("/player/:id/consume/:card", setConsumeHandlerFunc(playerMap))
	router.PUT("/player/:id/discard/:card", setDiscardHandlerFunc(playerMap))
	router.PUT("/player/:id/play/:card", setPlayHandlerFunc(playerMap))
	//Config endpoints
	router.GET("/config/turnTimer", getRemainingTimerHandlerFunc(*turnTime))
	router.GET("/config/turnLength", getConfigTurnTimerHandlerFunc())
	router.GET("/config/mapSize", getConfigMapSizeHandlerFunc())
	router.GET("/config/hasWon", getConfigGameStateHandlerFunc(*hasWon))
	router.GET("/config", getAllConfigHandlerFunc(*turnTime, *hasWon))
	router.Run("0.0.0.0:8080")
}

func getPlayerOrNil(playerMap map[string]*Player, id string) *Player {
	var player = playerMap[id]
	return player
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

func addPlayerHandlerFunc(playerMap *map[string]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var pName = filterPlayerName(c.Param("name"))
		var pID = addPlayer(playerMap, pName)
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

func tileToMapPiece(tile *Tile) MapPiece {
	if tile == nil {
		return MapPiece{Edge.toString(), 0, 0, 0, 0, 0, 0}
	}
	var planNorth, planEast, planSouth, planWest = 0, 0, 0, 0
	for _, player := range tile.Players {
		switch player.Direction {
		case North:
			planNorth++
		case East:
			planEast++
		case South:
			planSouth++
		case West:
			planWest++
		}
	}
	return MapPiece{
		tile.Terrain.toString(),
		tile.Zombies,
		len(tile.Players),
		planNorth,
		planEast,
		planSouth,
		planWest,
	}
}

func getSurroundingsHandlerFunc(playerMap map[string]*Player, gameMap [mapWidth][mapHeight]*Tile) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		player := getPlayerOrNil(playerMap, id)
		if player == nil { //TODO: If nil else function or invert? Make them all identical!
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else {
			var NW = tileToMapPiece(gameMap[player.X-1][player.Y-1])
			var NN = tileToMapPiece(gameMap[player.X][player.Y-1])
			var NE = tileToMapPiece(gameMap[player.X+1][player.Y-1])
			var WW = tileToMapPiece(gameMap[player.X-1][player.Y])
			var CE = tileToMapPiece(gameMap[player.X][player.Y])
			var EE = tileToMapPiece(gameMap[player.X+1][player.Y])
			var SW = tileToMapPiece(gameMap[player.X-1][player.Y+1])
			var SS = tileToMapPiece(gameMap[player.X][player.Y+1])
			var SE = tileToMapPiece(gameMap[player.X+1][player.Y+1])

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
			c.JSON(http.StatusOK, miniMap)
			return
		}
	}
	return fn
}

func setDiscardHandlerFunc(playerMap *map[string]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card, err = strconv.Atoi(cardStr)
		if err != nil || card >= len(cardTypes) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		playerPtr := getPlayerOrNil(*playerMap, id)
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

func setPlayHandlerFunc(playerMap *map[string]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card, err = strconv.Atoi(cardStr)
		if err != nil || card >= len(cardTypes) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		playerPtr := getPlayerOrNil(*playerMap, id)
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

func setConsumeHandlerFunc(playerMap *map[string]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card = cards[strings.ToLower(cardStr)]
		playerPtr := getPlayerOrNil(*playerMap, id)
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

func setDirectionHandlerFunc(playerMap *map[string]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		dirStr := c.Param("dir")
		var dir = directions[strings.ToLower(dirStr)]
		playerPtr := getPlayerOrNil(*playerMap, id)
		if playerPtr != nil {
			(*playerPtr).Direction = dir
			c.Status(http.StatusOK)
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
	return fn
}

func getPlayerHandlerFunc(playerMap map[string]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		playerPtr := getPlayerOrNil(playerMap, id)
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

func maskPlayerInfo(tile *Tile) {
	for playerNr, player := range tile.Players {
		//Filter the dead
		if player.Alive {
			tile.Players = append(tile.Players[:playerNr], tile.Players[playerNr+1:]...)
			continue
		}
		tile.Players[playerNr].ID = ""
		//Blank out hidden info
		tile.Players[playerNr].Cards = [5]Card{}
		tile.Players[playerNr].Consume = None
		tile.Players[playerNr].Play = None
		tile.Players[playerNr].IsBot = true
	}
}
