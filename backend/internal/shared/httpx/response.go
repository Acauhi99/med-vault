package httpx

import (
	"encoding/json"
	"net/http"
	"time"
)

type Response struct {
	Data  any           `json:"data"`
	Error *ErrorPayload `json:"error"`
	Meta  Meta          `json:"meta"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details []any  `json:"details"`
}

type Meta struct {
	RequestID string    `json:"request_id"`
	Timestamp time.Time `json:"timestamp"`
	Page      *int      `json:"page,omitempty"`
	PageSize  *int      `json:"page_size,omitempty"`
	Total     *int      `json:"total,omitempty"`
}

func WriteJSON(w http.ResponseWriter, r *http.Request, status int, data any) {
	write(w, r, status, Response{Data: data, Meta: meta(r)})

}

func WriteJSONWithMeta(w http.ResponseWriter, r *http.Request, status int, data any, responseMeta Meta) {
	responseMeta.RequestID = RequestID(r)
	if responseMeta.Timestamp.IsZero() {
		responseMeta.Timestamp = time.Now().UTC()
	}
	write(w, r, status, Response{Data: data, Meta: responseMeta})
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	write(w, r, status, Response{
		Error: &ErrorPayload{Code: code, Message: message, Details: []any{}},
		Meta:  meta(r),
	})
}

func write(w http.ResponseWriter, r *http.Request, status int, response Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}

func meta(r *http.Request) Meta {
	return Meta{RequestID: RequestID(r), Timestamp: time.Now().UTC()}
}

func RequestID(r *http.Request) string {
	if requestID, ok := RequestIDFromContext(r.Context()); ok {
		return requestID
	}
	return ""
}
