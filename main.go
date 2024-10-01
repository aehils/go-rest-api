package main

import (
	"net/http"

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
var albums = []album{
	{ID: "001", Title: "My Beautiful Dark Twisted Fantasy", Artiste: "Kanye West", Year: 2010, Price: 59.99},
	{ID: "002", Title: "Moon Music", Artiste: "Coldplay", Year: 2024, Price: 34.99},
	{ID: "003", Title: "Modus Vivendi", Artiste: "070 Shake", Year: 2020, Price: 32.59},
	{ID: "004", Title: "HARDSTONE PSYCHO", Artiste: "Don Toliver", Year: 2024, Price: 42.00},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getfullAlbums)
	// associating the POST method at the /albums path with the postAlbums function
	router.POST("/albums", addAlbum)
	// associating a new GET method with a function that iterates through all resources
	// returns a single resource by its ID
	router.GET("/albums/:id", getAlbumID)
	router.Run("localhost:8888")
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
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)
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
