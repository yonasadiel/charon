package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Request interface of Charon Http Request Wrapper
type Request interface {
	GetRequestData() map[string]string
	GetSessionData(key string) interface{}
	SetSessionData(key string, value interface{})
	SaveSession()
	SendJSON(output interface{}, code int)
}

// HTTPHandler receive Charon wrapped request and ressponse
type HTTPHandler func(Request)

// HTTPRequest wrapper of Charon Http Request
type HTTPRequest struct {
	r *http.Request
	w http.ResponseWriter
	s *sessions.Session
}

// GetRequestData return the data of request
func (req *HTTPRequest) GetRequestData() map[string]string {
	return mux.Vars(req.r)
}

// GetSessionData return the data of session with known key
func (req *HTTPRequest) GetSessionData(key string) interface{} {
	return req.s.Values[key]
}

// SetSessionData set the data of session
func (req *HTTPRequest) SetSessionData(key string, value interface{}) {
	req.s.Values[key] = value
}

// SaveSession saves the session to the cookie
func (req *HTTPRequest) SaveSession() {
	req.s.Save(req.r, req.w)
}

// SendJSON write json as http response
func (req *HTTPRequest) SendJSON(output interface{}, code int) {
	response, _ := json.Marshal(output)

	req.w.Header().Set("Content-Type", "application/json")
	req.w.WriteHeader(code)
	req.w.Write(response)
}

// MockRequest is Request object that is mocked for testing purposes
type MockRequest struct {
	RequestData  map[string]string
	SessionData  map[string]interface{}
	JSONResponse []byte
	StatusCode   int
}

// GetRequestData return the data of request
func (req *MockRequest) GetRequestData() map[string]string {
	return req.RequestData
}

// GetSessionData return the data of session with known key
func (req *MockRequest) GetSessionData(key string) interface{} {
	return req.SessionData[key]
}

// SetSessionData set the data of session
func (req *MockRequest) SetSessionData(key string, value interface{}) {
	req.SessionData[key] = value
}

// SaveSession do nothing because the session is already saved
func (req *MockRequest) SaveSession() {
	//
}

// SendJSON write json as http response
func (req *MockRequest) SendJSON(output interface{}, code int) {
	var err error
	req.JSONResponse, err = json.Marshal(output)
	if err != nil {
		req.StatusCode = http.StatusInternalServerError
	} else {
		req.StatusCode = code
	}
}
