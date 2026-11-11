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

	// 2. Select theme
	theme := GetTheme(req.Theme)

	// 3. Initialize PDF canvas in memory
	pdf := gofpdf.New("P", "mm", pageSize, "")
	pdf.SetMargins(theme.MarginLeft, theme.MarginTop, theme.MarginRight)
	pdf.AddPage()

	// 4. Render Header / Title
	pdf.SetFont(theme.FontFamily, "B", 20)
	pdf.SetTextColor(theme.HeaderColor.R, theme.HeaderColor.G, theme.HeaderColor.B)
	pdf.CellFormat(0, 15, req.Title, "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Draw divider rule using theme rule color
	pdf.SetDrawColor(theme.RuleColor.R, theme.RuleColor.G, theme.RuleColor.B)
	pdf.SetLineWidth(0.5)
	pdf.Line(theme.MarginLeft, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(6)

	// 5. Parse & Render Markdown Lines using iterators
	for rawLine := range strings.SplitSeq(req.Markdown, "\n") {
		line := strings.TrimRight(rawLine, "\r")

		// Empty line -> vertical space
		if strings.TrimSpace(line) == "" {
			pdf.Ln(4)
			continue
		}

		// Heading 3: ### Sub-section (must be evaluated before # and ##)
		if strings.HasPrefix(line, "### ") {
			pdf.SetFont(theme.FontFamily, "B", 11)
			pdf.SetTextColor(theme.AccentColor.R, theme.AccentColor.G, theme.AccentColor.B)
			pdf.CellFormat(0, 7, strings.TrimPrefix(line, "### "), "", 1, "L", false, 0, "")
			pdf.Ln(1.5)
			continue
		}

		// Heading 2: ## Subtitle
		if strings.HasPrefix(line, "## ") {
			pdf.SetFont(theme.FontFamily, "B", 13)
			pdf.SetTextColor(theme.AccentColor.R, theme.AccentColor.G, theme.AccentColor.B)
			pdf.CellFormat(0, 8, strings.TrimPrefix(line, "## "), "", 1, "L", false, 0, "")
			pdf.Ln(2)
			continue
		}

		// Heading 1: # Title
		if strings.HasPrefix(line, "# ") {
			pdf.SetFont(theme.FontFamily, "B", 16)
			pdf.SetTextColor(theme.HeaderColor.R, theme.HeaderColor.G, theme.HeaderColor.B)
			pdf.CellFormat(0, 10, strings.TrimPrefix(line, "# "), "", 1, "L", false, 0, "")
			pdf.Ln(2)
			continue
		}

		// Bullet point: - or *
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			cleanText := strings.TrimPrefix(strings.TrimPrefix(line, "- "), "* ")
			pdf.SetFont(theme.FontFamily, "", 10)
			pdf.SetTextColor(theme.AccentColor.R, theme.AccentColor.G, theme.AccentColor.B)
			pdf.CellFormat(6, 6, chr(149), "", 0, "R", false, 0, "") // Bullet point symbol
			pdf.SetTextColor(theme.BodyColor.R, theme.BodyColor.G, theme.BodyColor.B)
			pdf.CellFormat(0, 6, " "+cleanText, "", 1, "L", false, 0, "")
			continue
		}

		// Regular body paragraph
		pdf.SetFont(theme.FontFamily, "", 10)
		pdf.SetTextColor(theme.BodyColor.R, theme.BodyColor.G, theme.BodyColor.B)
		pdf.MultiCell(0, 6, line, "", "L", false)
	}

	// 6. Output PDF bytes directly to an in-memory buffer
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
