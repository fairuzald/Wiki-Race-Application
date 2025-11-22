package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"backend/models"
)

// AutoCompleteHandler handles autocomplete requests to fetch data from the Wikipedia API
func AutoCompleteHandler(c *gin.Context) {
	// Get search term and limit from query parameter
	searchTerm := c.Query("search")
	limit := c.Query("limit")

	// Create query parameters for the Wikipedia API request
	queryParams := url.Values{}
	queryParams.Set("action", "query")
	queryParams.Set("format", "json")
	queryParams.Set("gpssearch", searchTerm)
	queryParams.Set("generator", "prefixsearch")
	queryParams.Set("prop", "pageprops|pageimages|pageterms|info")
	queryParams.Set("inprop", "url")
	queryParams.Set("redirects", "")
	queryParams.Set("ppprop", "displaytitle")
	queryParams.Set("piprop", "thumbnail")
	queryParams.Set("pilimit", "max")
	// Set the thumbnail maximum size to 160px
	queryParams.Set("pithumbsize", "160")
	queryParams.Set("wbptterms", "description")
	queryParams.Set("gpsnamespace", "0")

	// Set gpslimit only if limit is provided and is not null
	if limit != "" {
		queryParams.Set("gpslimit", limit)
	}

	queryParams.Set("origin", "*")

	// Create the URL for the Wikipedia API request
	apiURL := fmt.Sprintf("https://en.wikipedia.org/w/api.php?%s", queryParams.Encode())

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create request"})
		return
	}

	req.Header.Set("User-Agent", "GoGoPowerRangers/1.0 (https://github.com/yourusername/GoGoPowerRangers; contact@example.com)")

	// Send the GET request to Wikipedia API
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch data from Wikipedia"})
		return
	}
	defer resp.Body.Close()

	// Check if response status is OK
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Wikipedia API returned status %d: %s", resp.StatusCode, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Wikipedia API returned status %d", resp.StatusCode),
		})
		return
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to read response body"})
		return
	}

	// Decode the JSON response from Wikipedia API into the AutoComplete struct
	var result models.AutoComplete
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		log.Printf("Failed to decode JSON response: %v\nResponse body: %s", err, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to decode JSON response",
			"error":   err.Error(),
		})
		return
	}

	// Handle the case when there are no results
	if len(result.Query.Pages) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": nil, "message": "No results found"})
		return
	}

	// Format and append the results to allResults
	var formattedResults []gin.H

	for _, page := range result.Query.Pages {
		description := ""
		if len(page.Terms.Description) > 0 {
			description = page.Terms.Description[0]
		}

		formattedResult := gin.H{
			"pageid":      page.PageID,
			"title":       page.Title,
			"description": description,
			"image":       page.Thumbnail.Source,
			"url":         page.FullURL,
		}
		formattedResults = append(formattedResults, formattedResult)
	}

	// Return the combined formatted results as JSON
	if len(formattedResults) == 0 {
		c.JSON(http.StatusOK, gin.H{"data": nil, "message": "No results found"})

	} else {
		c.JSON(http.StatusOK, gin.H{"data": formattedResults, "message": "Results retrieved successfully"})
	}
}
