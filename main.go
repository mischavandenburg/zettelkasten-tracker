package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	db                *sql.DB
	markdownFileCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "zettelkasten_markdown_file_count",
		Help: "Current number of markdown files in Zettelkasten",
	})
)

func init() {
	// Use environment variables for database connection details
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Construct the connection string
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}

	// Initialize Prometheus metric
	prometheus.MustRegister(markdownFileCount)
}

func postCount(c *gin.Context) {
	countStr := c.PostForm("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count"})
		return
	}

	// Update the database
	_, err = db.Exec("INSERT INTO markdown_counts (count, timestamp) VALUES ($1, $2)", count, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not insert into database"})
		return
	}

	// Update Prometheus gauge
	markdownFileCount.Set(float64(count))

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func main() {
	router := gin.Default()

	// Route to accept counts
	router.POST("/count", postCount)

	// Prometheus metrics endpoint
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Start serving the application
	log.Fatal(router.Run(":3009"))
}
