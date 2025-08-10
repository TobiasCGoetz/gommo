package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupAPI() {
	router := gin.Default()

	// Configure CORS to allow all origins for local development
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	// Add middleware for error handling and logging
	router.Use(errorHandlingMiddleware())

	router.GET("/player/:id", getPlayerHandler)
	router.GET("/player/:id/surroundings", getSurroundingsHandler)
	router.GET("/config", getAllConfigHandler)
	router.POST("/player/:name", addPlayerHandler)
	router.PUT("/player/:id/direction/:dir", setDirectionHandler)
	router.PUT("/player/:id/play/:cardType", setPlayHandler)

	// Event log endpoints
	router.GET("/player/:id/events", getPlayerEventsHandler)
	router.GET("/player/:id/events/type/:eventType", getPlayerEventsByTypeHandler)

	router.Run("0.0.0.0:8080")
}

func getAllConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gState)
}

func addPlayerHandler(c *gin.Context) {
	var pId = pMap.addPlayer(
		filterPlayerName(c.Param("name")),
		gMap.getNewPlayerEntryTile())
	c.JSON(http.StatusOK, pId)
}

func getSurroundingsHandler(c *gin.Context) {
	id := c.Param("id")
	var player = pMap.getPlayer(id)

	// Special case: if player is dead, return laboratory surrounded by edges
	if !player.Alive {
		deadPlayerSurroundings := Surroundings{
			NW: MapPiece{Edge.toString(), 0, 0, 0, 0, 0, 0},
			NN: MapPiece{Edge.toString(), 0, 0, 0, 0, 0, 0},
			NE: MapPiece{Edge.toString(), 0, 0, 0, 0, 0, 0},
			WW: MapPiece{Edge.toString(), 0, 0, 0, 0, 0, 0},
			CE: MapPiece{Laboratory.toString(), 0, 0, 0, 0, 0, 0},
			EE: MapPiece{Edge.toString(), 0, 0, 0, 0, 0, 0},
			SW: MapPiece{Edge.toString(), 0, 0, 0, 0, 0, 0},
			SS: MapPiece{Edge.toString(), 0, 0, 0, 0, 0, 0},
			SE: MapPiece{Edge.toString(), 0, 0, 0, 0, 0, 0},
		}
		c.JSON(http.StatusOK, deadPlayerSurroundings)
		return
	}

	var xPos = player.CurrentTile.XPos
	var yPos = player.CurrentTile.YPos
	c.JSON(http.StatusOK, gMap.getSurroundingsFromPos(xPos, yPos))
	return
}

func setPlayHandler(c *gin.Context) {
	id := c.Param("id")
	cardStr := c.Param("cardType")
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

func setDirectionHandler(c *gin.Context) {
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

func getPlayerHandler(c *gin.Context) {
	id := c.Param("id")
	playerPtr := getPlayerOrNil(id)
	if playerPtr != nil {
		c.JSON(http.StatusOK, (*playerPtr))
		return
	} else {
		c.AbortWithStatus(http.StatusForbidden)
	}
}

func filterPlayerName(name string) string {
	if len(name) > gameConfig.Player.NameMaxLength {
		return name[0:gameConfig.Player.NameMaxLength]
	}
	//TODO: Filter bad words
	return name
}

// getPlayerEventsHandler returns a handler for getting recent events for a player
func getPlayerEventsHandler(c *gin.Context) {
	playerID := c.Param("id")

	// Verify player exists
	if playerID != "" && pMap.getPlayerPtr(playerID) == nil {
		sendErrorResponse(c, http.StatusNotFound, "player_not_found", "Player not found")
		return
	}

	// Parse query parameters
	lastTurns := 5 // Default to last 5 turns
	if turnsStr := c.Query("turns"); turnsStr != "" {
		if turns, err := strconv.Atoi(turnsStr); err == nil && turns > 0 {
			lastTurns = turns
		}
	}

	// Get events with filters
	filters := EventFilters{
		LastTurns: lastTurns,
	}

	events := eventLogger.GetPlayerEvents(playerID, filters)

	sendSuccessResponse(c, http.StatusOK, gin.H{
		"events": events,
		"count":  len(events),
	})
}

// getPlayerEventsByTypeHandler returns a handler for getting filtered events for a player
func getPlayerEventsByTypeHandler(c *gin.Context) {
	playerID := c.Param("id")
	eventType := EventType(c.Param("eventType"))

	// Verify player exists
	if playerID != "" && pMap.getPlayerPtr(playerID) == nil {
		sendErrorResponse(c, http.StatusNotFound, "player_not_found", "Player not found")
		return
	}

	// Validate event type
	validType := false
	for _, et := range EventTypeList() {
		if et == eventType {
			validType = true
			break
		}
	}

	if !validType {
		sendErrorResponse(c, http.StatusBadRequest, "invalid_event_type", "Invalid event type")
		return
	}

	// Parse query parameters
	lastTurns := 5 // Default to last 5 turns
	if turnsStr := c.Query("turns"); turnsStr != "" {
		if turns, err := strconv.Atoi(turnsStr); err == nil && turns > 0 {
			lastTurns = turns
		}
	}

	// Get filtered events
	events := eventLogger.GetPlayerEvents(playerID, EventFilters{
		EventType: eventType,
		LastTurns: lastTurns,
	})

	sendSuccessResponse(c, http.StatusOK, gin.H{
		"events": events,
		"count":  len(events),
		"type":   eventType,
	})
}

// errorHandlingMiddleware provides centralized error handling and logging
func errorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered in API handler: %v", r)
				sendErrorResponse(c, http.StatusInternalServerError, "internal_server_error", "An unexpected error occurred")
			}
		}()

		c.Next()

		// Handle other HTTP errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			sendErrorResponse(c, c.Writer.Status(), "api_error", err.Error())
		}
	}
}

// sendErrorResponse sends a standardized error response
func sendErrorResponse(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}

// sendSuccessResponse sends a standardized success response
func sendSuccessResponse(c *gin.Context, status int, data gin.H) {
	if data == nil {
		data = gin.H{}
	}

	// Add success flag to response
	data["success"] = true

	c.JSON(status, data)
}
