package main

import (
	"time"

	"cpf-cnpj-validator/database"
	"cpf-cnpj-validator/handlers"
	"cpf-cnpj-validator/status"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       12 * time.Hour,
	}))
	router.Use(status.CountRequests())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	router.GET("/status", status.Handler)

	router.POST("/documents", handlers.CreateDocument)
	router.GET("/documents", handlers.GetDocuments)
	router.GET("/documents/:id", handlers.GetDocumentByID)
	router.PUT("/documents/:id", handlers.UpdateDocument)
	router.PATCH("/documents/:id/blocklist", handlers.UpdateBlocklist)
	router.DELETE("/documents/:id", handlers.DeleteDocument)

	router.Run(":8080")
}