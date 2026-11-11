package service

import (
	"strings"
	"testing"

	"github.com/yashG0/papercut/internal/model"
)

var sampleDoc = `# System Architecture Overview
## In-Memory Pipeline Design
- Zero disk overhead with direct buffer writes.
- Stream line tokenization powered by Go 1.24 iterators.
- Memory protection enforced prior to unmarshaling.

### Notes
This payload simulates a typical multi-section technical document rendering cycle with dynamic typography.`

func BenchmarkGeneratePDF_Small(b *testing.B) {
	svc := NewDocumentService()
	req := model.RenderPDFRequest{
		Title:    "BenchmarkSmall",
		Markdown: sampleDoc,
		Theme:    "modern",
		PageSize: "A4",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		buf, err := svc.GeneratePDF(req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		_ = len(buf)
	}
}

func BenchmarkGeneratePDF_Large(b *testing.B) {
	svc := NewDocumentService()
	largeMarkdown := strings.Repeat(sampleDoc+"\n\n", 50)
	req := model.RenderPDFRequest{
		Title:    "BenchmarkLarge",
		Markdown: largeMarkdown,
		Theme:    "invoice",
		PageSize: "A4",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		buf, err := svc.GeneratePDF(req)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		_ = len(buf)
	}
}
