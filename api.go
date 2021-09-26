package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func setupAPI() {
	router := gin.Default()
	router.GET("/player")
	router.Run("localhost:8080")
}

func getPlayer(id int, c *gin.Context) {
	c.IndentedJSON(http.StatusOK, playerList[0])
}

