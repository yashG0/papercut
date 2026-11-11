package service

import (
	"bytes"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/yashG0/papercut/internal/model"
)

// DocumentService defines the rendering behavior
type DocumentService struct{}

// NewDocumentService creates a new service instance
func NewDocumentService() *DocumentService {
	return &DocumentService{}
}

// GeneratePDF transforms Markdown content into a PDF binary stored in memory
func (s *DocumentService) GeneratePDF(req model.RenderPDFRequest) ([]byte, error) {
	// 1. Resolve default page size (fallback to A4)
	pageSize := req.PageSize
	if pageSize == "" {
		pageSize = "A4"
	}

	// 2. Initialize PDF canvas in memory
	pdf := gofpdf.New("P", "mm", pageSize, "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// 3. Render Header / Title
	pdf.SetFont("Helvetica", "B", 20)
	pdf.CellFormat(0, 15, req.Title, "", 1, "L", false, 0, "")
	pdf.Ln(4)

	// Draw a light divider rule below the title
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(6)

	// 4. Parse & Render Markdown Lines
	for rawLine := range strings.SplitSeq(req.Markdown, "\n") {
		line := strings.TrimRight(rawLine, "\r")

		// Empty line -> vertical space
		if strings.TrimSpace(line) == "" {
			pdf.Ln(4)
			continue
		}

		// Heading 1: # Title
		if strings.HasPrefix(line, "# ") {
			pdf.SetFont("Helvetica", "B", 16)
			pdf.SetTextColor(30, 30, 30)
			pdf.CellFormat(0, 10, strings.TrimPrefix(line, "# "), "", 1, "L", false, 0, "")
			pdf.Ln(2)
			continue
		}

		// Heading 2: ## Subtitle
		if strings.HasPrefix(line, "## ") {
			pdf.SetFont("Helvetica", "B", 13)
			pdf.SetTextColor(60, 60, 60)
			pdf.CellFormat(0, 8, strings.TrimPrefix(line, "## "), "", 1, "L", false, 0, "")
			pdf.Ln(2)
			continue
		}

		// Bullet point: - or *
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			cleanText := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")
			pdf.SetFont("Helvetica", "", 11)
			pdf.SetTextColor(50, 50, 50)
			pdf.CellFormat(6, 6, chr(149), "", 0, "R", false, 0, "") // 149 is bullet point symbol
			pdf.CellFormat(0, 6, " "+cleanText, "", 1, "L", false, 0, "")
			continue
		}

		// Regular body paragraph
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(50, 50, 50)
		pdf.MultiCell(0, 6, line, "", "L", false)
	}

	// 5. Output PDF bytes directly to an in-memory buffer
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Helper to convert character codes for PDF symbols
func chr(c int) string {
	return string([]byte{byte(c)})
}
