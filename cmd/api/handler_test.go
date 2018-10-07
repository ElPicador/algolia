package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"tpaulmyer/algolia/datetree"
)

var sampleData = `
2015-08-03 00:00:07	Elixir
2015-08-03 00:00:07	Plop
2015-08-04 12:10:07	will_this_test_succed_?
2015-08-05 23:00:12	Elixir
2015-08-22 00:00:08	yeah
2015-08-22 00:00:08	will_this_test_succed_?
2015-08-22 15:10:10	will_this_test_succed_?
2015-08-23 00:00:09	SoftLayer
2015-09-03 00:05:09	experience
2015-09-03 00:05:11	Elixir
2015-09-03 00:05:13	experience
2015-09-03 00:18:13	will_this_test_succed_?
2015-09-10 11:05:10	hungary
`

func TestHandlerCount(t *testing.T) {
	// create silent handler
	h := Handler{
		Logger:   log.New(ioutil.Discard, "", 0),
		DateTree: getTreeForTests(t),
	}

	t.Run("no result", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/count/2015-08-02%2010:05:07", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Count)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		body := w.Body.String()
		want := `{"count":0}`
		if body != want {
			t.Errorf("wanted %s, got %s", want, body)
		}
	})

	t.Run("year count", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/count/2015", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Count)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		body := w.Body.String()
		want := `{"count":7}`
		if body != want {
			t.Errorf("wanted %s, got %s", want, body)
		}
	})

	t.Run("month count", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/count/2015-08", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Count)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		body := w.Body.String()
		want := `{"count":5}`
		if body != want {
			t.Errorf("wanted %s, got %s", want, body)
		}
	})

	t.Run("day count", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/count/2015-08-22", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Count)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		body := w.Body.String()
		want := `{"count":2}`
		if body != want {
			t.Errorf("wanted %s, got %s", want, body)
		}
	})

	t.Run("hour count", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/count/2015-09-03%2000", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Count)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		body := w.Body.String()
		want := `{"count":3}`
		if body != want {
			t.Errorf("wanted %s, got %s", want, body)
		}
	})

	t.Run("minute count", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/count/2015-09-03%2000:05", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Count)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		body := w.Body.String()
		want := `{"count":2}`
		if body != want {
			t.Errorf("wanted %s, got %s", want, body)
		}
	})

	t.Run("second count", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/count/2015-08-03%2000:00:07", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Count)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		body := w.Body.String()
		want := `{"count":2}`
		if body != want {
			t.Errorf("wanted %s, got %s", want, body)
		}
	})

	t.Run("no date", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/count/", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Count)).ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest {
			t.Error("code should be bad request")
			return
		}
		body := w.Body.String()
		want := `{"error":"no date specified in url"}`
		if body != want {
			t.Errorf("wanted %s, got %s", want, body)
		}
	})

	t.Run("wrong date", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/count/zorglub", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Count)).ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest {
			t.Error("code should be bad request")
			return
		}
		if !strings.Contains(w.Body.String(), "Failed to parse date") {
			t.Errorf("error message is wrong")
		}
	})

	t.Run("unpredictable layout", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/count/2015-08-02%2015:04:05-999999999", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Count)).ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest {
			t.Error("code should be bad request")
			return
		}
		if !strings.Contains(w.Body.String(), "unknown date format") {
			t.Errorf("error message is wrong")
		}
	})
}

func TestHandlerPopularity(t *testing.T) {
	// create silent handler
	h := Handler{
		Logger:   log.New(ioutil.Discard, "", 0),
		DateTree: getTreeForTests(t),
	}

	getBody := func(t *testing.T, w *httptest.ResponseRecorder) PopularResult {
		var ret PopularResult
		err := json.Unmarshal(w.Body.Bytes(), &ret)
		if err != nil {
			t.Errorf("unmarshal of response failed: %+v", err)
		}

		return ret
	}

	t.Run("no result", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015-08-02%2010:05:07?size=8", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		ret := getBody(t, w)
		if len(ret.Queries) != 0 {
			t.Error("payload should be empty")
		}
	})

	t.Run("no size", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest {
			t.Error("code should be ok")
			return
		}
		body := w.Body.String()
		want := `{"error":"missing size parameter"}`
		if body != want {
			t.Errorf("wanted %s, got %s", want, body)
		}
	})

	t.Run("year popularity", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015?size=10", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		res := getBody(t, w)
		if len(res.Queries) != 7 {
			t.Error("wrong size for result array")
		}
		expected := PopularResult{Queries: []Query{
			{Query: "will_this_test_succed_?", Count: 4},
			{Query: "Elixir", Count: 3},
			{Query: "experience", Count: 2},
			{Query: "Plop", Count: 1},
			{Query: "SoftLayer", Count: 1},
			{Query: "hungary", Count: 1},
			{Query: "yeah", Count: 1},
		}}
		for i := range res.Queries {
			if res.Queries[i] != expected.Queries[i] {
				t.Error("arrays should be equal")
			}
		}
	})

	t.Run("size smaller than result", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015?size=3", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		res := getBody(t, w)
		if len(res.Queries) != 3 {
			t.Error("wrong size for result array")
		}
		expected := PopularResult{Queries: []Query{
			{Query: "will_this_test_succed_?", Count: 4},
			{Query: "Elixir", Count: 3},
			{Query: "experience", Count: 2},
		}}
		for i := range res.Queries {
			if res.Queries[i] != expected.Queries[i] {
				t.Error("arrays should be equal")
			}
		}
	})

	t.Run("negative size", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015?size=-1", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest {
			t.Error("code should be bad request")
			return
		}
		if !strings.Contains(w.Body.String(), "size parameter cannot be negative") {
			t.Error("wrong error message")
			return
		}
	})

	t.Run("size is not a number", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015?size=lol", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusBadRequest {
			t.Error("code should be bad request")
			return
		}
		if !strings.Contains(w.Body.String(), "size parameter invalid") {
			t.Error("wrong error message")
			return
		}
	})

	t.Run("month popularity", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015-08?size=10", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		res := getBody(t, w)
		if len(res.Queries) != 5 {
			t.Error("wrong size for result array")
		}
		expected := PopularResult{Queries: []Query{
			{Query: "will_this_test_succed_?", Count: 3},
			{Query: "Elixir", Count: 2},
			{Query: "Plop", Count: 1},
			{Query: "SoftLayer", Count: 1},
			{Query: "yeah", Count: 1},
		}}
		for i := range res.Queries {
			if res.Queries[i] != expected.Queries[i] {
				t.Error("arrays should be equal")
			}
		}
	})

	t.Run("day popularity", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015-08-22?size=5", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		res := getBody(t, w)
		if len(res.Queries) != 2 {
			t.Error("wrong size for result array")
		}
		expected := PopularResult{Queries: []Query{
			{Query: "will_this_test_succed_?", Count: 2},
			{Query: "yeah", Count: 1},
		}}
		for i := range res.Queries {
			if res.Queries[i] != expected.Queries[i] {
				t.Error("arrays should be equal")
			}
		}
	})

	t.Run("hour popularity", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015-09-03%2000?size=5", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		res := getBody(t, w)
		if len(res.Queries) != 3 {
			t.Error("wrong size for result array")
		}
		expected := PopularResult{Queries: []Query{
			{Query: "experience", Count: 2},
			{Query: "Elixir", Count: 1},
			{Query: "will_this_test_succed_?", Count: 1},
		}}
		for i := range res.Queries {
			if res.Queries[i] != expected.Queries[i] {
				t.Error("arrays should be equal")
			}
		}
	})

	t.Run("minute popularity", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015-09-03%2000:05?size=10", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		res := getBody(t, w)
		if len(res.Queries) != 2 {
			t.Error("wrong size for result array")
		}
		expected := PopularResult{Queries: []Query{
			{Query: "experience", Count: 2},
			{Query: "Elixir", Count: 1},
		}}
		for i := range res.Queries {
			if res.Queries[i] != expected.Queries[i] {
				t.Error("arrays should be equal")
			}
		}
	})

	t.Run("second popularity", func(t *testing.T) {
		r := httptest.NewRequest("GET", "/v1/queries/popularity/2015-08-03%2000:00:07?size=10", nil)
		w := httptest.NewRecorder()

		h.DateMiddleware(http.HandlerFunc(h.Popular)).ServeHTTP(w, r)
		if w.Code != http.StatusOK {
			t.Error("code should be ok")
			return
		}
		res := getBody(t, w)
		if len(res.Queries) != 2 {
			t.Error("wrong size for result array")
		}
		expected := PopularResult{Queries: []Query{
			{Query: "Elixir", Count: 1},
			{Query: "Plop", Count: 1},
		}}
		for i := range res.Queries {
			if res.Queries[i] != expected.Queries[i] {
				t.Error("arrays should be equal")
			}
		}
	})
}

func getTreeForTests(t *testing.T) *datetree.Tree {
	r := strings.NewReader(sampleData)
	tsvr := NewTSVReader(r)
	go func() {
		for err := range tsvr.Error() {
			t.Error(err)
		}
	}()

	// create new date tree and insert every line in it
	tree := datetree.NewTree()
	for l := range tsvr.Lines() {
		if len(l) != 2 {
			continue
		}

		ti, err := time.Parse(Second, l[0])
		if err != nil {
			t.Error(err)
			continue
		}

		tree.Insert(l[1], ti)
	}

	tree.IndexPopularity()
	return tree
}
