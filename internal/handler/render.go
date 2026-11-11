package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yashG0/papercut/internal/model"
	"github.com/yashG0/papercut/internal/service"
)

// DocumentHandler handles document rendering endpoints
type DocumentHandler struct {
	service *service.DocumentService
}

// NewDocumentHandler injects the service dependency
func NewDocumentHandler(s *service.DocumentService) *DocumentHandler {
	return &DocumentHandler{service: s}
}

// RenderPDF processes the Markdown input and streams the generated PDF back
func (h *DocumentHandler) RenderPDF(c *gin.Context) {
	var req model.RenderPDFRequest

	// 1. Bind and validate JSON body using struct tags
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Error:   "validation failed",
			Details: err.Error(),
		})
		return
	}

	// 2. Delegate rendering logic to the service layer
	pdfBytes, err := h.service.GeneratePDF(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Error:   "failed to generate PDF",
			Details: err.Error(),
		})
		return
	}

	// 3. Set download/attachment headers
	fileName := fmt.Sprintf("%s.pdf", req.Title)
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	// 4. Stream binary PDF bytes directly to the client response
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}
