package routes

import (
	"bookly-api-golang/controllers"
	"github.com/gin-gonic/gin"
	"bookly-api-golang/middlewares"
)

func CategoryRoutes(router *gin.RouterGroup) {
	categoryRoutes := router.Group("/categories", middleware.BasicAuth())
	{
		categoryRoutes.GET("/", controllers.GetAllCategory)
		categoryRoutes.GET("/:id", controllers.GetCategoryByID)
		categoryRoutes.POST("/", controllers.CreateCategory)
		categoryRoutes.PUT("/:id", controllers.UpdateCategory)
		categoryRoutes.DELETE("/:id", controllers.DeleteCategory)
		categoryRoutes.GET("/:id/books", controllers.GetCategoryBooks)
	}
}

func BookRoutes(router *gin.RouterGroup) {
	bookRoutes := router.Group("/books", middleware.BasicAuth())
	{
		bookRoutes.GET("/", controllers.GetAllBook)
		bookRoutes.GET("/:id", controllers.GetBookByID)
		bookRoutes.POST("/", controllers.CreateBook)
		bookRoutes.PUT("/:id", controllers.UpdateBook)
		bookRoutes.DELETE("/:id", controllers.DeleteBook)
	}
}