// Package docs serves the embedded OpenAPI specification and Swagger UI.
package docs

import (
	_ "embed"
	"net/http"
)

//go:embed swagger.yaml
var spec []byte

//go:embed index.html
var page []byte

// UI serves the standard Swagger UI page.
func UI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(page)
}

// Spec serves the specification used by Swagger UI.
func Spec(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.Write(spec)
}
