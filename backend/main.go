package main

import (
	"cpf-cnpj-validator/database"
	"cpf-cnpj-validator/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.POST("/documents", handlers.CreateDocument)
	router.GET("/documents", handlers.GetDocuments)
	router.GET("/documents/:id", handlers.GetDocumentByID)
	router.PUT("/documents/:id", handlers.UpdateDocument)

	router.Run(":8080")
}