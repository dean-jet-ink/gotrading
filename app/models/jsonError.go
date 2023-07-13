package models

import (
	"encoding/json"
	"log"
	"net/http"
)

type JSONError struct {
	ErrorMessage string `json:"error"`
	Code         int    `json:"code"`
}

func APIError(w http.ResponseWriter, errorMessage string, code int) {
	jsonError := &JSONError{
		ErrorMessage: errorMessage,
		Code:         code,
	}
	jsonErrorMarshal, err := json.Marshal(jsonError)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(jsonErrorMarshal)
}
