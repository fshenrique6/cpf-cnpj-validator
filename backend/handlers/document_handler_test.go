package handlers

import (
	"bytes"
	"cpf-cnpj-validator/database"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/documents", CreateDocument)
	return router
}

func TestCreateDocument_ValidCPF(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(1, now, now)

	mock.ExpectQuery("INSERT INTO documents").
		WithArgs("11144477735", "cpf", false).
		WillReturnRows(rows)

	router := setupRouter()

	body, _ := json.Marshal(map[string]string{"number": "111.444.777-35"})
	req := httptest.NewRequest(http.MethodPost, "/documents", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status esperado 201, recebido %d. corpo: %s", w.Code, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas do mock não cumpridas: %v", err)
	}
}

func TestCreateDocument_InvalidDocument(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	router := setupRouter()

	body, _ := json.Marshal(map[string]string{"number": "12345678900"})
	req := httptest.NewRequest(http.MethodPost, "/documents", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado 400, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}
