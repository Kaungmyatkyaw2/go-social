package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {

	maxBytes := 1_048_578

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)

	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func (app *application) writeJSONError(w http.ResponseWriter, status int, message string) error {
	return app.jsonResponse(w, status, message)

}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data    any  `json:"data,omitempty"`
		Success bool `json:"success"`
		Error   any  `json:"error,omitempty"`
	}

	response := &envelope{}

	response.Success = status >= 200 && status < 300

	if response.Success {
		response.Data = data
	} else {
		response.Error = data
	}

	return writeJSON(w, status, response)
}
