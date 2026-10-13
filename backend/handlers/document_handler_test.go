package handlers

import (
	"bytes"
	"cpf-cnpj-validator/database"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/documents", CreateDocument)
	router.GET("/documents/:id", GetDocumentByID)
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

func TestCreateDocument_Duplicate(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	mock.ExpectQuery("INSERT INTO documents").
		WithArgs("11144477735", "cpf", false).
		WillReturnError(&pgconn.PgError{Code: "23505"})

	router := setupRouter()

	body, _ := json.Marshal(map[string]string{"number": "111.444.777-35"})
	req := httptest.NewRequest(http.MethodPost, "/documents", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("status esperado 409, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}

func TestCreateDocument_DatabaseError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	mock.ExpectQuery("INSERT INTO documents").
		WithArgs("11144477735", "cpf", false).
		WillReturnError(errors.New("conexão perdida"))

	router := setupRouter()

	body, _ := json.Marshal(map[string]string{"number": "111.444.777-35"})
	req := httptest.NewRequest(http.MethodPost, "/documents", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status esperado 500, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}

func TestGetDocumentByID_Found(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "number", "type", "blocklisted", "created_at", "updated_at"}).
		AddRow(1, "11144477735", "cpf", false, now, now)

	mock.ExpectQuery("SELECT id, number, type, blocklisted, created_at, updated_at").
		WithArgs(1).
		WillReturnRows(rows)

	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/documents/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status esperado 200, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}

func TestGetDocumentByID_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	mock.ExpectQuery("SELECT id, number, type, blocklisted, created_at, updated_at").
		WithArgs(999).
		WillReturnError(sql.ErrNoRows)

	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/documents/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status esperado 404, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}

func TestGetDocumentByID_InvalidID(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/documents/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado 400, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}
