package model

type RenderPDFRequest struct {
	Title    string `json:"title" binding:"required,min=1,max=100"`
	Markdown string `json:"markdown" binding:"required,min=1"`
	Theme    string `json:"theme" binding:"omitempty,oneof=default modern invoice"`
	PageSize string `json:"page_size" binding:"omitempty,oneof=A4 Letter"`
}

type HealthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
