package server

import (
	"encoding/json"
	"net/http"
)

type Message struct {
	//Header map[string]interface{}
	Header http.Header `json:"Header"`

	Body json.RawMessage
}
