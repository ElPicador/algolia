package main

import (
	"context"
	"net/http"
	"time"
)

// DateMiddleware is a middleware responsible from parsing the date in the
// URL, fetching the corresponding date layout and inserting the information
// into the context to be used by the handlers.
func (h *Handler) DateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Logger.Printf("request received [%s]", r.RequestURI)
		date, err := GetDateFromURL(r)
		if err != nil {
			out := APIError{Error: err.Error()}
			h.Respond(w, out, http.StatusBadRequest)
			return
		}

		// try to parse date in order to check its validity
		layout, err := GetLayout(date)
		if err != nil {
			out := APIError{Error: err.Error()}
			h.Respond(w, out, http.StatusBadRequest)
			return
		}

		t, err := time.Parse(layout, date)
		if err != nil {
			out := APIError{Error: "Failed to parse date: " + err.Error()}
			h.Respond(w, out, http.StatusBadRequest)
			return
		}

		r = SetDateInContext(DateInfo{Time: t, Layout: layout}, r)
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// DateInfo represents information on a date.
type DateInfo struct {
	time.Time
	Layout string
}

type contextKey string

const date contextKey = "algolia.dateinfo"

// SetDateInContext sets a DateInfo in context to be used by other handlers.
func SetDateInContext(d DateInfo, r *http.Request) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), date, d))
}

// GetDateInContext returns a DateInfo from context.
func GetDateInContext(r *http.Request) DateInfo {
	return r.Context().Value(date).(DateInfo)
}
