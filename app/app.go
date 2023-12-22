package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	handler "github.com/demkowo/goquery/handlers"
	"github.com/demkowo/goquery/repository/postgres"
	service "github.com/demkowo/goquery/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

const portNumber = ":5000"

var (
	router  = gin.Default()
	connStr string
)

func init() {

	// Load the .env file
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	host := os.Getenv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	sslmode := os.Getenv("DB_SSLMODE")
	connStr = fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s", host, dbname, user, password, sslmode)
}

func Start() {

	// Check DB connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	// Add collector
	repo := postgres.New(db)
	service := service.New(repo)
	handler := handler.New(service)
	addRoutes(handler)

	router.Run(portNumber)
}
