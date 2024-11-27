package main

import (
	"encoding/json"
	"log"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"net/http"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/restaurants", getRestaurants)

	router.Run(":8080") // Listen + serve on localhost:8080
}

func getRestaurants(c *gin.Context) {
	latitude := c.Query("latitude")
	longitude := c.Query("longitude")
	chain := c.Query("chain")

	if latitude == "" || longitude == "" || chain == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required query parameters."})
		return
	}

	restaurants, err := searchRestaurants(chain, longitude, latitude)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"restaurants": restaurants})
}

func searchRestaurants(chainName, longitude, latitude string) ([]map[string]interface{}, error) {
	// Build URL for Mapbox API request
	accessToken := "pk.eyJ1IjoibXRvb21leS1kZXYiLCJhIjoiY200MDRiOWxqMTlzODJ3cHlrejE3NjgzNyJ9.AqCgaLym-Ym9jBSrsj6aGg"
	if accessToken == "" {
		log.Fatal("MAPBOX_ACCESS_TOKEN is not set.")
	}

	url := "https://api.mapbox.com/geocoding/v5/mapbox.places/" + chainName + ".json" +
		"?proximity=" + longitude + "," + latitude +
		"&limit=10" + "&types=poi" + "&access_token=" + accessToken

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	// Extract restaurant data from response
	features, ok := data["features"].([]interface{})
	if !ok {
		return nil, nil
	}

	restaurants := make([]map[string]interface{}, 0)
	for _, feature := range features {
		f, ok := feature.(map[string]interface {})
		if !ok {
			// Move to next restaurant in the response (if any)
			continue
		}

		geometry, ok := f["geometry"].(map[string]interface{})
        if !ok {
            continue
        }

        coordinates, ok := geometry["coordinates"].([]interface{})
        if !ok || len(coordinates) < 2 {
            continue
        }

        properties, ok := f["properties"].(map[string]interface{})
        if !ok {
            continue
        }

		// Collect restaurant information
        restaurant := map[string]interface{}{
            "name":      f["text"],
            "address":   properties["address"],
            "latitude":  coordinates[1],
            "longitude": coordinates[0],
        }
		
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}