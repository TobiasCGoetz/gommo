package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
	router.Run("0.0.0.0:8080")
}

// getPlayerOrNil returns a pointer to the referenced player or nil
//
// This will perform a lookup given a playerID and return a pointer to the Player or nil.
// Parameters:
//
//	playerList: PlayerList to be searched
//	id: PlayerID that will be searched for
//
// Returns: *Player or nil
func getPlayerOrNil(playerList *[]*Player, id string) *Player {
	for pNr, player := range *playerList {
		if player.ID == id {
			return (*playerList)[pNr]
		}
	}
	return nil
}

func parseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, errors.New("JWT didn't parse.")
		}
		return token, nil
	})
	return token, err
}

func authPlayer(tokenString string, userName string) bool {
	token, err := parseJWT(tokenString)
	if err != nil {
		return verifyJWT(*token, userName)
	}
	return true
}

func addPlayerHandlerFunc(playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		var pName = filterPlayerName(c.Param("name"))
		var pID = addPlayer(playerList, pName)
		c.IndentedJSON(http.StatusOK, pID)
	}
	return fn
}

func getRemainingTimerHandlerFunc(turnTimer *uint8) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, *turnTimer)
	}
	return fn
}

// TODO: Add surrounding players info
func getSurroundingsHandlerFunc(playerList *[]*Player, gameMap *[mapWidth][mapHeight]*Tile) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		playerPtr := getPlayerOrNil(playerList, id)
		if playerPtr == nil { //TODO: If nil else function or invert? Make them all identical!
			c.AbortWithStatus(http.StatusForbidden)
			return
		} else {
			player := *playerPtr

			//Construct empty minimap
			var NW = Tile{Edge, -1, []Player{}}
			var NN = Tile{Edge, -1, []Player{}}
			var NE = Tile{Edge, -1, []Player{}}
			var WW = Tile{Edge, -1, []Player{}}
			var CE = *gameMap[player.X][player.Y]
			var EE = Tile{Edge, -1, []Player{}}
			var SW = Tile{Edge, -1, []Player{}}
			var SS = Tile{Edge, -1, []Player{}}
			var SE = Tile{Edge, -1, []Player{}}

			//Fill minimap
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
			sanitizeSurroundingsInfo(&miniMap)
			c.IndentedJSON(http.StatusOK, miniMap)
			return
		}
	}
	return fn
}

func setDiscardHandlerFunc(playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card, err = strconv.Atoi(cardStr)
		if err != nil || card >= len(cardTypes) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		playerPtr := getPlayerOrNil(playerList, id)
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

func setPlayHandlerFunc(playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card, err = strconv.Atoi(cardStr)
		if err != nil || card >= len(cardTypes) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		playerPtr := getPlayerOrNil(playerList, id)
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

func setConsumeHandlerFunc(playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		cardStr := c.Param("card")
		var card, err = strconv.Atoi(cardStr)
		if err != nil || card >= len(cardTypes) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		playerPtr := getPlayerOrNil(playerList, id)
		if playerPtr != nil {
			(*playerPtr).Consume = cardTypes[card]
			c.Status(http.StatusOK)
			return
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
	return fn
}

func setDirectionHandlerFunc(playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		dirStr := c.Param("dir")
		var dir, err = strconv.Atoi(dirStr)
		if err != nil || dir >= len(Directions) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		playerPtr := getPlayerOrNil(playerList, id)
		if playerPtr != nil {
			(*playerPtr).Direction = Directions[dir]
			c.Status(http.StatusOK)
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
	return fn
}

func getPlayerHandlerFunc(playerList *[]*Player) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		id := c.Param("id")
		//passwd := c.Param("passwd")
		//TODO: USE A MAP HERE
		playerPtr := getPlayerOrNil(playerList, id)
		if playerPtr != nil {
			c.IndentedJSON(http.StatusOK, (*playerPtr))
			return
		} else {
			c.AbortWithStatus(http.StatusForbidden)
		}
	}
	return fn
}

func sanitizeSurroundingsInfo(surroundings *Surroundings) {
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
