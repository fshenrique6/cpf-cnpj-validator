package handlers

import (
	"net/http"

	"cpf-cnpj-validator/database"
	"cpf-cnpj-validator/models"
	"cpf-cnpj-validator/validator"

	"github.com/gin-gonic/gin"
)

func CreateDocument(c *gin.Context) {
	var newDocument models.Document

	if err := c.ShouldBindJSON(&newDocument); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	number := validator.OnlyDigits(newDocument.Number)

	docType, err := validator.DetectType(number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
		INSERT INTO documents (number, type, blocklisted)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err = database.DB.QueryRow(query, number, docType, newDocument.Blocklisted).
		Scan(&newDocument.ID, &newDocument.CreatedAt, &newDocument.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "documento já cadastrado ou erro ao criar"})
		return
	}

	newDocument.Number = number
	newDocument.Type = docType

	c.JSON(http.StatusCreated, newDocument)
}