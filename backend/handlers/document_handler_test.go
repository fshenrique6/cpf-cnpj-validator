package handlers

import (
	"bytes"
	"cpf-cnpj-validator/database"
	"cpf-cnpj-validator/models"
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
	router.PATCH("/documents/:id/blocklist", UpdateBlocklist)
	router.DELETE("/documents/:id", DeleteDocument)
	router.GET("/documents", GetDocuments)
	return router
}

//testes CreateDocument

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

//testes GetDocumentByID

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

//testes UpdateBlocklist

func TestUpdateBlocklist_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "number", "type", "blocklisted", "created_at", "updated_at"}).
		AddRow(1, "11144477735", "cpf", true, now, now)

	mock.ExpectQuery("UPDATE documents").
		WithArgs(true, 1).
		WillReturnRows(rows)

	router := setupRouter()

	body, _ := json.Marshal(map[string]bool{"blocklisted": true})
	req := httptest.NewRequest(http.MethodPatch, "/documents/1/blocklist", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status esperado 200, recebido %d. corpo: %s", w.Code, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas do mock não cumpridas: %v", err)
	}
}

func TestUpdateBlocklist_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	mock.ExpectQuery("UPDATE documents").
		WithArgs(true, 999).
		WillReturnError(sql.ErrNoRows)

	router := setupRouter()

	body, _ := json.Marshal(map[string]bool{"blocklisted": true})
	req := httptest.NewRequest(http.MethodPatch, "/documents/999/blocklist", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status esperado 404, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}

//testes DeleteDocument

func TestDeleteDocument_Success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	mock.ExpectExec("DELETE FROM documents").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	router := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/documents/1", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status esperado 200, recebido %d. corpo: %s", w.Code, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas do mock não cumpridas: %v", err)
	}
}

func TestDeleteDocument_NotFound(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	mock.ExpectExec("DELETE FROM documents").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0))

	router := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/documents/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status esperado 404, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}

func TestDeleteDocument_InvalidID(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	router := setupRouter()

	req := httptest.NewRequest(http.MethodDelete, "/documents/abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado 400, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}

//testes GetDocument

func TestGetDocuments_NoFilters(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "number", "type", "blocklisted", "created_at", "updated_at"}).
		AddRow(1, "11144477735", "cpf", false, now, now).
		AddRow(2, "11222333000181", "cnpj", true, now, now)

	mock.ExpectQuery("FROM documents WHERE 1=1 ORDER BY created_at DESC").
		WillReturnRows(rows)

	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/documents", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status esperado 200, recebido %d. corpo: %s", w.Code, w.Body.String())
	}

	var docs []models.Document
	if err := json.Unmarshal(w.Body.Bytes(), &docs); err != nil {
		t.Fatalf("resposta não é um JSON válido: %v", err)
	}
	if len(docs) != 2 {
		t.Errorf("esperado 2 documentos, recebido %d", len(docs))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas do mock não cumpridas: %v", err)
	}
}

func TestGetDocuments_FilterByType(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	rows := sqlmock.NewRows([]string{"id", "number", "type", "blocklisted", "created_at", "updated_at"})

	mock.ExpectQuery("WHERE 1=1 AND type = .* ORDER BY created_at DESC").
		WithArgs("cpf").
		WillReturnRows(rows)

	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/documents?type=cpf", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status esperado 200, recebido %d. corpo: %s", w.Code, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas do mock não cumpridas: %v", err)
	}
}

func TestGetDocuments_FilterByBlocklisted(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	rows := sqlmock.NewRows([]string{"id", "number", "type", "blocklisted", "created_at", "updated_at"})

	mock.ExpectQuery("WHERE 1=1 AND blocklisted = .* ORDER BY created_at DESC").
		WithArgs(true).
		WillReturnRows(rows)

	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/documents?blocklisted=true", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status esperado 200, recebido %d. corpo: %s", w.Code, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas do mock não cumpridas: %v", err)
	}
}

func TestGetDocuments_SortAndOrder(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	rows := sqlmock.NewRows([]string{"id", "number", "type", "blocklisted", "created_at", "updated_at"})

	mock.ExpectQuery("ORDER BY number ASC").
		WillReturnRows(rows)

	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/documents?sort_by=number&order=asc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status esperado 200, recebido %d. corpo: %s", w.Code, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("expectativas do mock não cumpridas: %v", err)
	}
}

func TestGetDocuments_InvalidSortField(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/documents?sort_by=nomeinventado", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado 400, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}

func TestGetDocuments_InvalidOrder(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/documents?order=sideways", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status esperado 400, recebido %d. corpo: %s", w.Code, w.Body.String())
	}
}

func TestGetDocuments_EmptyReturnsArray(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("erro ao criar mock: %v", err)
	}
	defer mockDB.Close()

	database.DB = mockDB

	rows := sqlmock.NewRows([]string{"id", "number", "type", "blocklisted", "created_at", "updated_at"})
	mock.ExpectQuery("FROM documents").WillReturnRows(rows)

	router := setupRouter()

	req := httptest.NewRequest(http.MethodGet, "/documents", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Body.String() != "[]" {
		t.Errorf("lista vazia deveria ser [], recebido %s", w.Body.String())
	}
}
