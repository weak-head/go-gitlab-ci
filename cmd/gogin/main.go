package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.GET("/info", infoEndpoint)
		v1.GET("/user/:name", getUserNameEndpoint)
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.GET("/info", infoEndpoint)
		v2.GET("/user/:name", getUserNameEndpoint)
		v2.GET("/user/:name/*action", getUserNameActionEndpoint)
		v2.GET("/usr/group", getUserGroupEndpoint)
	}

	router.Run(":8080")
}

func infoEndpoint(c *gin.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"app":      "gogit",
		"hostname": hostname,
		"version": gin.H{
			"commit_id":  "s26q6Qo3QG",
			"build_date": "2022-01-11",
		},
	})
}

func getUserNameEndpoint(c *gin.Context) {
	name := c.Param("name")

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Hello %s", name),
		"name":    name,
	})
}

func getUserNameActionEndpoint(c *gin.Context) {
	name := c.Param("name")
	action := c.Param("action")

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Hello %s, should we %s?", name, action),
		"name":    name,
		"action":  action,
	})
}

func getUserGroupEndpoint(c *gin.Context) {
	c.String(http.StatusOK, "The available groups are [...]")
}
