package app

import (
	"net/http"
)

// Middleware wrap Charon HTTPHandler to default http.Handler
func Middleware(f HTTPHandler) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req HTTPRequest
		req.r = r
		req.s = Charon.getSession(r)
		req.w = w

		f(&req)
	})
}
