package main

import (
	"encoding/json"
	"log"
	"net/http"
	"tpaulmyer/algolia/datetree"
	"tpaulmyer/algolia/null"
)

// Handler is the structure that handles calls to the API.
type Handler struct {
	Logger   *log.Logger
	DateTree *datetree.Tree
}

// APIError is used to return a JSON error to the user.
type APIError struct {
	Error string `json:"error"`
}

// CountResult is used to return the result of a count query to the API user.
type CountResult struct {
	Count int `json:"count"`
}

// Count is the handler responsible for the /1/queries/count/<DATE_PREFIX> route.
func (h *Handler) Count(w http.ResponseWriter, r *http.Request) {
	t := GetDateInContext(r)
	s := datetree.Search{
		Year:   t.Year(),
		Month:  null.Int{Valid: len(t.Layout) >= len(Month), Int: int(t.Month())},
		Day:    null.Int{Valid: len(t.Layout) >= len(Day), Int: t.Day()},
		Hour:   null.Int{Valid: len(t.Layout) >= len(Hour), Int: t.Hour()},
		Minute: null.Int{Valid: len(t.Layout) >= len(Minute), Int: t.Minute()},
		Second: null.Int{Valid: len(t.Layout) >= len(Second), Int: t.Second()},
	}

	count := h.DateTree.Count(s)
	out := CountResult{Count: count}
	h.Respond(w, out, http.StatusOK)
}

// Query is the API representation of a query.
type Query struct {
	Query string `json:"query"`
	Count int    `json:"count"`
}

// PopularResult ils the result returned to the API user from the Popular route.
type PopularResult struct {
	Queries []Query `json:"queries"`
}

// Popular is the handler responsible for the /1/queries/popular/<DATE_PREFIX> route.
func (h *Handler) Popular(w http.ResponseWriter, r *http.Request) {
	t := GetDateInContext(r)
	size, err := GetSizeParameter(r)
	if err != nil {
		out := APIError{Error: err.Error()}
		h.Respond(w, out, http.StatusBadRequest)
		return
	}

	s := datetree.Search{
		Year:       t.Year(),
		Month:      null.Int{Valid: len(t.Layout) >= len(Month), Int: int(t.Month())},
		Day:        null.Int{Valid: len(t.Layout) >= len(Day), Int: t.Day()},
		Hour:       null.Int{Valid: len(t.Layout) >= len(Hour), Int: t.Hour()},
		Minute:     null.Int{Valid: len(t.Layout) >= len(Minute), Int: t.Minute()},
		Second:     null.Int{Valid: len(t.Layout) >= len(Second), Int: t.Second()},
		Popularity: size,
	}
	pop := h.DateTree.Popular(s)

	out := PopularResult{
		Queries: make([]Query, len(pop)),
	}
	for i, v := range pop {
		out.Queries[i] = Query(v)
	}

	h.Respond(w, out, http.StatusOK)
}

// Respond returns a payload and a statuscode to the user.
func (h *Handler) Respond(w http.ResponseWriter, out interface{}, statusCode int) {
	d, err := json.Marshal(out)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"an internal error occurred"}`))
		h.Logger.Println("got internal error while marshalling response:", err.Error())
		return
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(d)
	if err != nil {
		h.Logger.Println("failed to return response to user:", err)
		return
	}

	h.Logger.Println("response sent with status", statusCode)
}
