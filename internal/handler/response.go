package handler

import (
	"encoding/json"
	"net/http"
)

type ResponseWrapper struct {
	Error    *ErrorResponse `json:"error,omitempty"`
	Response interface{}    `json:"response,omitempty"`
	Data     interface{}    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}

func respond(w http.ResponseWriter, code int, data interface{}, resp interface{}, err *ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	res := ResponseWrapper{
		Error:    err,
		Response: resp,
		Data:     data,
	}
	json.NewEncoder(w).Encode(res)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respond(w, code, nil, nil, &ErrorResponse{
		Code: code,
		Text: msg,
	})
}
