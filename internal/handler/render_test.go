package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yashG0/papercut/internal/handler"
	"github.com/yashG0/papercut/internal/model"
	"github.com/yashG0/papercut/internal/service"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	svc := service.NewDocumentService()
	h := handler.NewDocumentHandler(svc)

	r.POST("/api/v1/render/pdf", h.RenderPDF)
	return r
}

func TestRenderPDF_Success(t *testing.T) {
	router := setupTestRouter()

	reqBody := model.RenderPDFRequest{
		Title:    "Test_Document",
		Markdown: "# Heading\nThis is a test paragraph.",
		Theme:    "default",
		PageSize: "A4",
	}
	rawJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/render/pdf", bytes.NewBuffer(rawJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d. Body: %s", w.Code, w.Body.String())
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/pdf" {
		t.Errorf("expected Content-Type application/pdf, got %s", contentType)
	}

	// PDF files always begin with the magic header bytes "%PDF-"
	if !bytes.HasPrefix(w.Body.Bytes(), []byte("%PDF-")) {
		t.Errorf("response body does not start with PDF magic bytes (%%PDF-)")
	}
}

func TestRenderPDF_ValidationFailure_EmptyMarkdown(t *testing.T) {
	router := setupTestRouter()

	reqBody := model.RenderPDFRequest{
		Title:    "Missing_Markdown_Doc",
		Markdown: "", // Empty: violates `binding:"required,min=1"`
	}
	rawJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/render/pdf", bytes.NewBuffer(rawJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d", w.Code)
	}
}

func TestRenderPDF_ValidationFailure_InvalidTheme(t *testing.T) {
	router := setupTestRouter()

	reqBody := model.RenderPDFRequest{
		Title:    "Invalid_Theme_Doc",
		Markdown: "Content",
		Theme:    "neon-cyberpunk", // Violates `oneof=default modern invoice`
	}
	rawJSON, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/render/pdf", bytes.NewBuffer(rawJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got %d", w.Code)
	}
}
