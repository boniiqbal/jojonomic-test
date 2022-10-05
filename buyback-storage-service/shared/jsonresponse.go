package shared

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Response struct {
	Error  bool        `json:"error"`
	ReffID string      `json:"reff_id,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func RespondWithError(w http.ResponseWriter, code int, err bool, reffID string, message string) {
	RespondWithJSON(w, code, map[string]string{"error": strconv.FormatBool(err), "reff_id": reffID, "message": message})
}

func RespondSuccess(w http.ResponseWriter, code int, err bool, reffID string, data interface{}) {
	response := &Response{
		Error:  err,
		ReffID: reffID,
		Data:   data,
	}
	RespondWithJSON(w, code, response)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
