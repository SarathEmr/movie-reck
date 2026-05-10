package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"movie-reck/config"
	"movie-reck/controller"
	"movie-reck/handler"
	"movie-reck/migration"
	"movie-reck/repo"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func connectDB(dsn string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	maxRetries := 10
	retryDelay := 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err == nil {
			if err := db.Ping(); err == nil {
				log.Printf("Successfully connected to database after %d attempts", i+1)
				return db, nil
			}
		}
		if i < maxRetries-1 {
			log.Printf("Failed to connect to database (attempt %d/%d), retrying in %s...", i+1, maxRetries, retryDelay)
			time.Sleep(retryDelay)
		}
	}
	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}

func main() {
	config.InitConfig()

	// Database connection
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")

	kkk := viper.GetString("GEMINI_API_KEY")
	log.Println("--- GEMINI_API_KEY is ", kkk)

	log.Printf("Connecting to database: host=%s port=%s user=%s dbname=%s", dbHost, dbPort, dbUser, dbName)

	// data source name (DSN) format for PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := connectDB(dsn)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	defer db.Close()

	// Run migrations
	migrationsPath, err := filepath.Abs("./migrations")
	if err != nil {
		log.Fatalf("failed to resolve migrations path: %v", err)
	}

	migrator := migration.NewMigrator(db)
	if err := migrator.Migrate(migrationsPath); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize layers
	userRepo := repo.NewUserRepo(db)
	loginController := controller.NewLoginController(userRepo)
	loginHandler := handler.NewLoginHandler(loginController)

	movieRepo := repo.NewMovieRepo(db)
	reckController := controller.NewReckController(movieRepo)
	reckHandler := handler.NewReckHandler(reckController)

	// Routes
	r.POST("/login", loginHandler.Login)
	r.GET("/movie/recommendation", reckHandler.GetRecommendations)

	r.Run(":8080")
}
