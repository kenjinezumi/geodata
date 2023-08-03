package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Replace the API URL with the specific GIBS API endpoint you want to query
	apiURL := "https://gibs.earthdata.nasa.gov/wmts/epsg4326/best/VIIRS_SNPP_CorrectedReflectance_TrueColor/default/2021-09-15/250m/{z}/{y}/{x}.jpg"

	// Define the coordinates of the Paris area
	minLat, maxLat := 48.8156, 48.9022 // Latitude range (south to north)
	minLon, maxLon := 2.2241, 2.4699   // Longitude range (west to east)

	// Define the zoom level and tile range to cover the Paris area
	zoom := 8
	minXTile, minYTile := getWebMercatorTileCoordinates(minLat, minLon, zoom)
	maxXTile, maxYTile := getWebMercatorTileCoordinates(maxLat, maxLon, zoom)

	// Loop through the tile range and fetch each image
	for xTile := minXTile; xTile <= maxXTile; xTile++ {
		for yTile := minYTile; yTile <= maxYTile; yTile++ {
			// Replace the placeholders in the API URL with the current tile coordinates
			apiURL = strings.ReplaceAll(apiURL, "{z}", fmt.Sprintf("%d", zoom))
			apiURL = strings.ReplaceAll(apiURL, "{x}", fmt.Sprintf("%d", xTile))
			apiURL = strings.ReplaceAll(apiURL, "{y}", fmt.Sprintf("%d", yTile))

			// Fetch the image using the API
			response, err := http.Get(apiURL)
			if err != nil {
				fmt.Println("Error fetching data:", err)
				return
			}
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				fmt.Println("API returned status:", response.Status)
				return
			}

			// Create an image file to save the retrieved image
			file, err := os.Create(fmt.Sprintf("satellite_image_%d_%d.jpg", xTile, yTile))
			if err != nil {
				fmt.Println("Error creating file:", err)
				return
			}
			defer file.Close()

			// Copy the response body (image data) to the file
			_, err = io.Copy(file, response.Body)
			if err != nil {
				fmt.Println("Error saving image:", err)
				return
			}

			fmt.Printf("Saved satellite image %d_%d.jpg\n", xTile, yTile)
		}
	}
}

// Function to calculate Web Mercator tile coordinates (x, y) from latitude and longitude
func getWebMercatorTileCoordinates(lat, lon float64, zoom int) (int, int) {
	n := 1 << uint(zoom)
	xTile := int((lon + 180.0) / 360.0 * float64(n))
	yTile := int((1.0 - (math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0)) / math.Pi)) / 2.0 * float64(n))
	return xTile, yTile
}
