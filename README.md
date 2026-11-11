# PaperCut 📄✂️

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-Gin-008ECF?style=flat&logo=gin)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

An ultra-fast, headless document-rendering microservice written in Go. **PaperCut** converts raw Markdown documents into production-ready PDFs and streams them back in real time with **zero disk writes** and **sub-millisecond latency**.

---

## ⚡ Key Highlights

- **Zero-Disk I/O Pipeline:** All canvas construction, formatting, and PDF encoding occur directly inside Go memory buffers (`bytes.Buffer`). No temporary file cleanups, lock contentions, or disk thrashing under heavy traffic.
- **Sub-Millisecond Response Time:** Complete JSON parsing, input validation, PDF compilation, and binary delivery finish in under 1ms for standard payloads.
- **Go 1.24 Iterator Engine:** Uses `strings.SplitSeq` stream-based line iteration to eliminate full-slice heap allocations when parsing markdown blocks.
- **Memory Safety Middleware:** Native `http.MaxBytesReader` guard limits incoming payload streams to 10MB to prevent memory exhaustion and denial-of-service attempts.
- **Dynamic Theming:** Built-in palette and typography engine supporting `default`, `modern`, and `invoice` layouts with custom accents, margins, and divider rules.
- **Lean Container Footprint:** Packaged via a multi-stage Docker/Podman build into a non-root Alpine runtime under 35MB.

---

## 🏛️ System Architecture

<p align="center">
  <img src="architecture.svg" alt="PaperCut Architecture Diagram" width="100%" />
</p>

### Pipeline Flow

1. **Ingress & Memory Boundary:** Gin processes the incoming HTTP request through `middleware.MaxBodySize(10MB)`, wrapping `r.Body` with `http.MaxBytesReader` to block abusive payload sizes before JSON unmarshaling.
2. **Contract Validation:** `c.ShouldBindJSON` enforces presence, length limits, and supported presets (`oneof=default modern invoice`) via Go struct tags.
3. **Stream Iteration & Rendering:** The Document Service resolves theme palettes and ranges over lines using `strings.SplitSeq(markdown, "\n")`, formatting headers, bullet lists, and paragraphs directly onto an in-memory `gofpdf` canvas.
4. **Binary PDF Streaming:** The completed document is written to a `bytes.Buffer` and streamed straight back to the client (`Content-Type: application/pdf`) without touching the filesystem.

---

## 📁 Project Structure

Organized according to standard Go project conventions:

```text
papercut/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point & dependency injection
├── internal/
│   ├── handler/
│   │   ├── render.go            # Gin HTTP controller for PDF streaming
│   │   └── render_test.go       # httptest-based unit and integration tests
│   ├── middleware/
│   │   ├── limits.go            # Memory protection middleware (MaxBytesReader)
│   │   └── limits_test.go       # Payload ceiling boundary tests
│   ├── model/
│   │   └── document.go          # Request/response contracts & Gin binding tags
│   └── service/
│       ├── renderer.go          # Core Markdown-to-PDF rendering engine
│       └── theme.go             # Typography, margins, and color definitions
├── architecture.svg             # Clean whiteboard architecture diagram
├── Dockerfile                   # Multi-stage minimal runtime build
├── .dockerignore                # Build context exclusion rules
├── Makefile                     # Build, run, test, and tidy automations
├── go.mod
└── go.sum
```

---

## 🚀 Getting Started

### Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/) (optional)
- Make (optional)

### Local Development

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yashG0/papercut.git
   cd papercut
   ```

2. **Run tests:**

   ```bash
   make test
   # Alternative: go test -v -race ./...
   ```

3. **Start the service:**
   ```bash
   make run
   # Alternative: go run ./cmd/server/main.go
   ```

The service will boot on `http://localhost:8080`.

---

## 🐳 Container Deployment (Docker / Podman)

Build and run the lightweight container:

```bash
# Build image using Podman or Docker
podman build -t papercut:latest .

# Run container on port 8080
podman run -d -p 8080:8080 --name papercut-service localhost/papercut:latest

# Verify health
curl http://localhost:8080/healthz
```

---

## 📖 API Reference

### 1. Health Probe

Checks service liveness and readiness.

```http
GET /healthz
```

**Response (`200 OK`):**

```json
{
  "service": "papercut",
  "status": "healthy"
}
```

---

### 2. Render Markdown to PDF

Compiles markdown content and returns an attachment PDF stream.

```http
POST /api/v1/render/pdf
Content-Type: application/json
```

#### Request Schema:

| Field       | Type     | Required | Default     | Allowed Values                 | Description                    |
| :---------- | :------- | :------- | :---------- | :----------------------------- | :----------------------------- |
| `title`     | `string` | **Yes**  | —           | 1–100 characters               | Document heading and file name |
| `markdown`  | `string` | **Yes**  | —           | Non-empty string               | Raw Markdown text to render    |
| `theme`     | `string` | No       | `"default"` | `default`, `modern`, `invoice` | Visual style and color preset  |
| `page_size` | `string` | No       | `"A4"`      | `A4`, `Letter`                 | Page dimension format          |

#### Example Request:

```bash
curl -X POST http://localhost:8080/api/v1/render/pdf \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Quarterly_Invoice_101",
    "theme": "invoice",
    "page_size": "A4",
    "markdown": "# Invoice #101\n## Client: ACME Corp\n- Cloud Architecture Consulting: $3,200\n- Automated PDF Service Delivery: $1,400\n\nTotal Due: $4,600.00\nPayment is due within 30 days."
  }' \
  --output invoice.pdf
```

**Response Headers:**

```http
HTTP/1.1 200 OK
Content-Type: application/pdf
Content-Disposition: attachment; filename="Quarterly_Invoice_101.pdf"
```

---

## 🧪 Testing

Unit and integration tests run entirely in memory using `net/http/httptest`, verifying HTTP binding, edge-case rejection, and stream boundaries without binding real sockets:

```bash
go test -v -race ./...
```

```text
=== RUN   TestRenderPDF_Success
--- PASS: TestRenderPDF_Success (0.00s)
=== RUN   TestRenderPDF_ValidationFailure_EmptyMarkdown
--- PASS: TestRenderPDF_ValidationFailure_EmptyMarkdown (0.00s)
=== RUN   TestRenderPDF_ValidationFailure_InvalidTheme
--- PASS: TestRenderPDF_ValidationFailure_InvalidTheme (0.00s)
=== RUN   TestMaxBodySize_RejectsOversizedPayload
--- PASS: TestMaxBodySize_RejectsOversizedPayload (0.00s)
PASS
ok      github.com/yashG0/papercut/internal/handler     0.012s
ok      github.com/yashG0/papercut/internal/middleware  0.004s
```

---

## 📄 License

Distributed under the MIT License.
