package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"net/http"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/ping", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/restaurants", func(c *gin.Context) {
		// for now, just return a static list
		c.JSON(http.StatusOK, gin.H{
			"restaurants": []map[string]interface{}{
				{
					"name": "Example Restaurant",
					"latitude": 40.7128,
					"longitude": -74.0060,
				},
				{
					"name": "Pizza Pasta Restaurant",
					"latitude": 11.23,
					"longitude": 12.12,
				},
			},
		})
	})

	router.Run(":8080") // Listen + serve on localhost:8080
}
