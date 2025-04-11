package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// getting started, you're creating the format (or structure) that resources will follow
// using type is akin to classes in OOP
type album struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Artiste string  `json:"artiste"`
	Year    int     `json:"year"`
	Price   float64 `json:"price"`
}

// this slice is the initial set of resources
var albums []album

func main() {

	// load albums from JSON
	loadAlbums()

	router := gin.Default()

	// associating album to return full list of albums
	router.GET("/albums", getfullAlbums)

	// associating the POST method at the /albums path with the postAlbums function
	router.POST("/albums", addAlbum)

	// associating a new GET method with a function that iterates through all resources
	// returns a single resource by its ID
	router.GET("/albums/:id", getAlbumID)

	// new function to slice albums ommitting an album at ID
	router.DELETE("/albums/:id", removeAlbumID)

	// function to associate update handler
	router.PUT("/albums/:id", albumEdit)

	router.Run("localhost:8888")
}

func loadAlbums() {
	file, err := os.ReadFile("albums.json")
	if err != nil {
		// albums will initiate as an empty list if theres an error reading from albums.json
		albums = []album{}
		return
	}
	json.Unmarshal(file, &albums)
}

// writing a handler to retreive the full list of resources, albums
func getfullAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

// this handler will attempt to parse JSON in a request body and bind it to a var newAlbum
// if it fails, the function will return early, if it succeeds, will append to albums
func addAlbum(c *gin.Context) {
	var newAlbum album
	if err := c.BindJSON(&newAlbum); err != nil {
		return
	}

	// checking
	if newAlbum.Title == "" || newAlbum.Artiste == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "missing atriste or title"})
		return
	}

	// checking
	inputTitle := strings.ToLower((strings.TrimSpace(newAlbum.Title)))
	inputArtiste := strings.ToLower((strings.TrimSpace(newAlbum.Artiste)))

	for _, a := range albums {
		if strings.ToLower(strings.TrimSpace(a.Title)) == inputTitle &&
			strings.ToLower(strings.TrimSpace(a.Artiste)) == inputArtiste {
			occupiedID := a.ID
			message := fmt.Sprintf("album already exists at %s", occupiedID)
			c.IndentedJSON(http.StatusConflict, gin.H{"message": message})
			return
		}
	}

	lastID := albums[len(albums)-1].ID
	idNum, _ := strconv.Atoi(lastID)
	newAlbum.ID = fmt.Sprintf("%03d", (idNum + 1))
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Added '%s' by %s at ID %s", newAlbum.Title, newAlbum.Artiste, newAlbum.ID),
	})
}

func getAlbumID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

// handler to remove an album from resource list, located by id
func removeAlbumID(c *gin.Context) {
	id := c.Param("id")

	for i, a := range albums {
		if a.ID == id {
			deletedTitle := a.Title
			deletedArtiste := a.Artiste
			message := fmt.Sprintf("%s by %s has been deleted.", deletedTitle, deletedArtiste)
			//removing the album at index i
			albums = append(albums[:i], albums[i+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": message})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
}

func albumEdit(c *gin.Context) {
	id := c.Param("id")
	var updatedAlbum album

	if err := c.BindJSON(&updatedAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body."})
		return
	}

	for i, a := range albums {
		if a.ID == id {
			updatedAlbum.ID = id
			albums[i] = updatedAlbum
			updatedTitle := a.Title
			updatedArtiste := a.Artiste
			message := fmt.Sprintf("changes made to %s by %s.", updatedTitle, updatedArtiste)
			c.IndentedJSON(http.StatusOK, gin.H{"message": message})
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found. Cannot PUT"})
}
