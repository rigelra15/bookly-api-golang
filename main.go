package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"bookly-api-golang/database"
	"bookly-api-golang/controllers"
	"bookly-api-golang/middlewares"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	"github.com/gin-contrib/cors"
	"time"
	_ "bookly-api-golang/docs"
	"os"

	_ "github.com/lib/pq"
)

var (
	DB *sql.DB
	err error
)

// @title Bookly API
// @version 1.0
// @description API untuk mengelola kategori dan buku di Bookly dengan menggunakan Golang dan PostgreSQL.
// @description Author: Rigel Ramadhani W. - Sanbercode Bootcamp Golang Batch 63
// @contact.name Rigel Ramadhani W.
// @contact.url https://github.com/rigelra15

// @host bookly-api-golang-production.up.railway.app
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @security BearerAuth

func main() {
	err = godotenv.Load("config/.env")
	if err != nil {
		panic("Error loading .env file")
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"),
	)

	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer DB.Close()
	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	database.DBMigrate(DB)

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		categoryRoutes := api.Group("/categories", middlewares.JWTAuthMiddleware())
		{
			categoryRoutes.GET("/", controllers.GetAllCategory)
			categoryRoutes.GET("/:id", controllers.GetCategoryByID)
			categoryRoutes.POST("/", controllers.CreateCategory)
			categoryRoutes.PUT("/:id", controllers.UpdateCategory)
			categoryRoutes.DELETE("/:id", controllers.DeleteCategory)
			categoryRoutes.GET("/:id/books", controllers.GetCategoryBooks)
		}

		bookRoutes := api.Group("/books", middlewares.JWTAuthMiddleware())
		{
			bookRoutes.GET("/", controllers.GetAllBook)
			bookRoutes.GET("/:id", controllers.GetBookByID)
			bookRoutes.POST("/", controllers.CreateBook)
			bookRoutes.PUT("/:id", controllers.UpdateBook)
			bookRoutes.DELETE("/:id", controllers.DeleteBook)
		}
		userRoutes := api.Group("/users")
		{
			userRoutes.POST("/login", controllers.Login)
			userRoutes.GET("/", controllers.GetAllUsers)
			userRoutes.GET("/:id", controllers.GetUserByID)
			userRoutes.POST("/", controllers.CreateUser)
			userRoutes.PUT("/:id", controllers.UpdateUser)
			userRoutes.DELETE("/:id", controllers.DeleteUser)
		}
	}

	router.Run(":" + os.Getenv("PORT"))
}