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
	


	router.GET("/player/:id", getPlayerHandlerFunc())
	router.GET("/player/:id/surroundings", getSurroundingsHandlerFunc())
	router.GET("/config", getAllConfigHandlerFunc())
	router.POST("/player/:name", addPlayerHandlerFunc())
	router.PUT("/player/:id/direction/:dir", setDirectionHandlerFunc())
	router.PUT("/player/:id/play/:cardType", setPlayHandlerFunc())
	
	// Event log endpoints
	router.GET("/player/:id/events", getPlayerEventsHandlerFunc())
	router.GET("/player/:id/events/type/:eventType", getPlayerEventsByTypeHandlerFunc())
	
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

// getPlayerEventsHandlerFunc returns a handler for getting recent events for a player
func getPlayerEventsHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
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
}

// getPlayerEventsByTypeHandlerFunc returns a handler for getting filtered events for a player
func getPlayerEventsByTypeHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
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
