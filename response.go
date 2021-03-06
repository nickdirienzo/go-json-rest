package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type ErrorBlock struct {
	Code    int    `json:"statusCode"`
	Message string `json:"message"`
}

type ResponseError struct {
	ApiVersion string     `json:"apiVersion"`
	Method     string     `json:"method"`
	Error      ErrorBlock `json:"error"`
}

// Inherit from an object implementing the http.ResponseWriter interface,
// and provide additional methods.
type ResponseWriter struct {
	http.ResponseWriter
	isIndented bool
}

// Encode the object in JSON, set the content-type header,
// and call Write.
func (self *ResponseWriter) WriteJson(v interface{}, statusCode int) error {
	self.Header().Set("content-type", "application/json")
	code := strconv.Itoa(statusCode)
	self.Header().Set("StatusCode", code)
	var b []byte
	var err error
	if self.isIndented {
		b, err = json.MarshalIndent(v, "", "  ")
	} else {
		b, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}
	self.Write(b)
	return nil
}

// Produce an error response in JSON with the following structure, '{"Error":"My error message"}'
// The standard plain text net/http Error helper can still be called like this:
// http.Error(w, "error message", code)
func Error(w *ResponseWriter, error string, code int, method string) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(code)
	resp := ResponseError{
		ApiVersion: "0",
		Error: ErrorBlock{
			Code:    code,
			Message: error,
		},
		Method: method,
	}
	err := w.WriteJson(&resp, code)
	if err != nil {
		panic(err)
	}
}

// Produce a 404 response with the following JSON, '{"Error":"Resource not found"}'
// The standard plain text net/http NotFound helper can still be called like this:
// http.NotFound(w, r.Request)
func NotFound(w *ResponseWriter, r *Request) {
	Error(w, "Resource not found", http.StatusNotFound, "generic")
}
