package main

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GetDateFromURL returns the date specified in the URL.
func GetDateFromURL(r *http.Request) (string, error) {
	if strings.Count(r.URL.Path, "/") != 4 {
		return "", errors.New("url badly formatted")
	}

	i := strings.LastIndex(r.URL.Path, "/")
	if i == len(r.URL.Path)-1 {
		return "", errors.New("no date specified in url")
	}

	// unescape the string, in case there are spaces (%20 in an encoded URL)
	return url.PathUnescape(r.URL.Path[i+1:])
}

// GetSizeParameter returns the size parameter for the popularity route.
func GetSizeParameter(r *http.Request) (int, error) {
	s := r.URL.Query().Get("size")
	if s == "" {
		return 0, errors.New("missing size parameter")
	}

	size, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("size parameter invalid: " + err.Error())
	}

	if size < 0 {
		return 0, errors.New("size parameter cannot be negative")
	}

	return size, nil
}
