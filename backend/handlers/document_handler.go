package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

var allowedSortFields = map[string]bool{
	"id":         true,
	"number":     true,
	"type":       true,
	"created_at": true,
	"updated_at": true,
}

func GetDocuments(c *gin.Context) {
	query := "SELECT id, number, type, blocklisted, created_at, updated_at FROM documents WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	if docType := c.Query("type"); docType != "" {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, docType)
		argIndex++
	}

	if blocklisted := c.Query("blocklisted"); blocklisted != "" {
		query += fmt.Sprintf(" AND blocklisted = $%d", argIndex)
		args = append(args, blocklisted == "true")
		argIndex++
	}

	sortBy := c.DefaultQuery("sort_by", "created_at")
	if !allowedSortFields[sortBy] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "campo de ordenação inválido"})
		return
	}

	order := strings.ToUpper(c.DefaultQuery("order", "desc"))
	if order != "ASC" && order != "DESC" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ordem inválida, use 'asc' ou 'desc'"})
		return
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, order)

	rows, err := database.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var documents []models.Document

	for rows.Next() {
		var d models.Document
		err := rows.Scan(&d.ID, &d.Number, &d.Type, &d.Blocklisted, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		documents = append(documents, d)
	}

	c.JSON(http.StatusOK, documents)
}

func GetDocumentByID(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var document models.Document
	err = database.DB.QueryRow(`
		SELECT id, number, type, blocklisted, created_at, updated_at
		FROM documents
		WHERE id = $1
	`, idInt).Scan(&document.ID, &document.Number, &document.Type, &document.Blocklisted, &document.CreatedAt, &document.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, document)
}

func UpdateDocument(c *gin.Context) {
	id := c.Param("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input models.Document
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	number := validator.OnlyDigits(input.Number)

	docType, err := validator.DetectType(number)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var updated models.Document
	err = database.DB.QueryRow(`
		UPDATE documents
		SET number = $1, type = $2, updated_at = now()
		WHERE id = $3
		RETURNING id, number, type, blocklisted, created_at, updated_at
	`, number, docType, idInt).Scan(
		&updated.ID, &updated.Number, &updated.Type, &updated.Blocklisted, &updated.CreatedAt, &updated.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}