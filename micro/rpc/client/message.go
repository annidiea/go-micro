package client

import "net/http"

type Message struct {
	//Header map[string]interface{}
	Header http.Header
	Body   interface{}
}
