package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type album struct {
	ID     string  `json:"-" db:"id"`
	Title  string  `json:"title" db:"title"`
	Artist string  `json:"artist" db:"artist"`
	Price  float64 `json:"price" db:"price"`
}

var db *sqlx.DB
var err error

func main() {

	// Database connection string
	dsn := "root:root@(localhost:3306)/albumsdb"

	// Initialize a mysql database connection
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	// Verify the connection to the database is still alive
	err = db.Ping()
	if err != nil {
		panic("Failed to ping the database: " + err.Error())
	}

	router := gin.Default()

	albums := router.Group("/album")
	{
		albums.PUT("/:id", updateAlbumById)
		albums.GET("/", getAlbums)
		albums.GET("/:id", getAlbumById)
		albums.POST("/", createAlbums)

	}

	router.Run(":8080")
}

func getAlbums(c *gin.Context) {
	albums := []album{}

	err = db.Select(&albums, "Select * from albums")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, albums)
}

func createAlbums(c *gin.Context) {
	var newAlbum album

	if err := c.BindJSON(&newAlbum); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	response, err := db.NamedExec("INSERT INTO albums (title, artist, price) VALUES (:title, :artist, :price)", &newAlbum)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	// Get the last inserted id
	lastId, err := response.LastInsertId()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	err = db.Get(&newAlbum, "SELECT id, title, artist, price FROM albums where id = ?", lastId)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusCreated, newAlbum)
}

func getAlbumById(c *gin.Context) {
	id := c.Param("id")
	var album album

	err = db.Get(&album, "SELECT id, title, artist, price FROM albums where id = ?", id)
	if err != nil {

		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
			return
		}
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, album)
}

func updateAlbumById(c *gin.Context) {
	id := c.Param("id")
	var alb album

	err = db.Get(&alb, "SELECT id, title, artist, price FROM albums where id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
			return
		}
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	if err := c.BindJSON(&alb); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	_, err = db.Query("UPDATE albums SET title = ?, artist = ?, price = ? WHERE id = ?", alb.Title, alb.Artist, alb.Price, id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	var updatedAlbum album
	err = db.Get(&updatedAlbum, "SELECT id, title, artist, price FROM albums where id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Album not found"})
			return
		}
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, updatedAlbum)

}
