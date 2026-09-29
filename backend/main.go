package main

import (
	"cpf-cnpj-validator/database"
	"cpf-cnpj-validator/handlers"
	"cpf-cnpj-validator/status"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

	router := gin.Default()
	router.Use(status.CountRequests())

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.POST("/documents", handlers.CreateDocument)
	router.GET("/documents", handlers.GetDocuments)
	router.GET("/documents/:id", handlers.GetDocumentByID)
	router.PUT("/documents/:id", handlers.UpdateDocument)
	router.PATCH("/documents/:id/blocklist", handlers.UpdateBlocklist)
	router.DELETE("/documents/:id", handlers.DeleteDocument)

	router.GET("/status", status.Handler)

	router.Run(":8080")
}