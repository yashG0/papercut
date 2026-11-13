# PaperCut 📄✂️

[![Go Version](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=flat&logo=go)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-Gin-008ECF?style=flat&logo=gin)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

A headless document-rendering microservice written in Go. **PaperCut** converts a subset of Markdown into PDFs and streams them back in the HTTP response, rendered entirely in memory with **no disk writes**.

---

## ⚡ Key Highlights

- **No-Disk Pipeline:** Canvas construction, formatting, and PDF encoding all happen in Go memory buffers (`bytes.Buffer`). There are no temporary files to create or clean up.
- **Measured Performance:** About 0.17 ms to lay out and encode a small document, measured with Go's native benchmark harness (hardware and caveats are listed under [Benchmarks](#-benchmarks)).
- **Allocation-Aware Line Parsing:** Uses `strings.SplitSeq` (Go 1.24) to iterate over lines without first allocating a slice of all lines.
- **Request Size Limit:** `http.MaxBytesReader` rejects request bodies over 10 MB before JSON decoding.
- **Dynamic Theming:** Built-in palette and typography engine supporting `default`, `modern`, and `invoice` layouts with custom accents, margins, and divider rules.
- **Small Container Image:** Multi-stage Docker/Podman build into a non-root Alpine runtime, under 35 MB.
- **Tested:** `httptest`-based handler and middleware tests, run with the race detector.

---

## 🏛️ System Architecture

<p align="center">
  <img src="asserts/architecture/arch.png" alt="PaperCut Architecture Diagram" width="100%" />
</p>

### Pipeline Flow

1. **Ingress & Size Limit:** Gin passes the request through `middleware.MaxBodySize(10MB)`, which wraps `r.Body` with `http.MaxBytesReader` to reject oversized bodies before JSON decoding.
2. **Contract Validation:** `c.ShouldBindJSON` enforces presence, length limits, and supported presets (`oneof=default modern invoice`) through Go struct tags.
3. **Line Iteration & Rendering:** The document service resolves theme palettes and ranges over lines with `strings.SplitSeq(markdown, "\n")`, formatting headers, bullet lists, and paragraphs onto an in-memory `gofpdf` canvas.
4. **Binary PDF Streaming:** The finished document is written to a `bytes.Buffer` and streamed back to the client (`Content-Type: application/pdf`) without touching the filesystem.

---

## 📝 Supported Markdown

PaperCut renders a deliberately small subset of Markdown:

| Element | Supported |
| :------ | :-------- |
| Headers | Yes |
| Bullet lists | Yes |
| Paragraphs | Yes |
| Tables, code blocks, links, numbered lists, images, inline bold/italic | No |

---

## 📊 Benchmarks

Microbenchmarks run against the rendering pipeline (`internal/service`) using Go's built-in benchmark harness (`go test -bench=. -benchmem`):

```text
goos: linux
goarch: amd64
pkg: github.com/yashG0/papercut/internal/service
cpu: AMD Ryzen 5 5600GT with Radeon Graphics
BenchmarkGeneratePDF_Small-12        7723        171498 ns/op      866748 B/op          430 allocs/op
BenchmarkGeneratePDF_Large-12         489       2277970 ns/op    11142393 B/op         9774 allocs/op
PASS
ok      github.com/yashG0/papercut/internal/service    2.442s
```

### Profile Breakdown

| Benchmark   | Input                                              | Latency (`ns/op`) | Approx. Time | Allocated per op | Allocs per op   |
| :---------- | :------------------------------------------------- | :---------------- | :----------- | :--------------- | :-------------- |
| **`Small`** | Standard 1-page document (~`<SMALL_KB>` KB)        | `171,498 ns`      | **~0.17 ms** | `846 KB`         | `430 allocs/op` |
| **`Large`** | 50 repeated multi-line blocks (~`<LARGE_KB>` KB)   | `2,277,970 ns`    | **~2.28 ms** | `10.6 MB`        | `9,774 allocs/op` |

**Caveats:** these benchmarks cover the rendering pipeline only, not the HTTP layer or JSON decoding, and were run on a single machine.

To reproduce locally:

```bash
go test -bench=. -benchmem ./internal/service/
```

---

## 📁 Project Structure

Organized according to standard Go project conventions:

```text
papercut/
├── assets/
│   └── architecture/
│       └── arch.png            # Architecture diagram
├── cmd/
│   └── server/
│       └── main.go             # Application entry point & dependency injection
├── internal/
│   ├── handler/
│   │   ├── render.go           # Gin HTTP controller for PDF streaming
│   │   └── render_test.go      # httptest-based unit & integration tests
│   ├── middleware/
│   │   ├── limits.go           # Request size limit (MaxBytesReader)
│   │   └── limits_test.go      # Payload ceiling boundary tests
│   ├── model/
│   │   └── document.go         # Request/response contracts & Gin binding tags
│   └── service/
│       ├── renderer.go         # Core Markdown-to-PDF rendering engine
│       ├── renderer_bench_test.go # Execution latency & allocation benchmarks
│       └── theme.go            # Typography, margins, and color definitions
├── Dockerfile                  # Multi-stage minimal runtime build
├── makefile                    # Build, run, test, and tidy automations
├── LICENSE
├── go.mod
├── go.sum
└── README.md
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
   go test -v -race ./...
   ```

3. **Start the service:**

   ```bash
   go run ./cmd/server/main.go
   ```

The service starts on `http://localhost:8080`.

---

## 🐳 Container Deployment (Docker / Podman)

Build and run the container:

```bash
# Build the image with Podman or Docker
podman build -t papercut:latest .

# Run the container on port 8080
podman run -d -p 8080:8080 --name papercut-service localhost/papercut:latest

# Verify health
curl http://localhost:8080/healthz
```

---

## 📖 API Reference

### 1. Health Probe

Checks service liveness.

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

Compiles Markdown content and returns the PDF as an attachment.

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

Handler and middleware tests run in memory with `net/http/httptest`, covering request binding, rejection of invalid input, and the request size limit, without opening real sockets:

```bash
go test -v -race ./...

# with coverage
go test -cover ./...
```

---

## ⚠️ Known Limitations

- **No concurrency limit.** Each render buffers the whole PDF in memory, so memory use grows with the number of parallel requests. The 10 MB cap limits the size of a single request, not total memory.
- **No authentication or rate limiting.**
- **Markdown subset only:** headers, bullet lists, and paragraphs (see [Supported Markdown](#-supported-markdown)).
- **Benchmarks cover the rendering pipeline only,** not the HTTP layer, and were run on one machine.

---

## 🔮 Next Steps

- Limit concurrent renders with a semaphore middleware
- Benchmark allocation across input sizes (1 KB to 10 MB)
- CI with GitHub Actions (`go test -race ./...`)
- Sanitize the `title` field before using it in `Content-Disposition`

---

## 📄 License

Distributed under the MIT License.
