package responsehelper

import (
	"encoding/json"
	"net/http"
)

// CommonResponse is model for common response
type CommonResponse struct {
	Message interface{} `json:"message"`
	Status  interface{} `json:"status"`
}

// RespondAsJSON sends json response
func RespondAsJSON(statusCode int, w http.ResponseWriter, res interface{}) (err error) {
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, err = w.Write(data)
	return err
}

// RespondWithErrorAsJSON sends json response with specified response and status code
func RespondWithErrorAsJSON(statusCode int, w http.ResponseWriter, msg string) (err error) {
	res := map[string]interface{}{
		"error": msg,
	}
	return RespondAsJSON(statusCode, w, res)
}
