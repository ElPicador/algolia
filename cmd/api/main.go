package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"tpaulmyer/algolia/datetree"
)

func main() {
	var port uint
	var file string
	flag.UintVar(&port, "p", 8080, "port the http server listen to")
	flag.StringVar(&file, "f", "hn_logs.tsv", "TSV file to read from")
	flag.Parse()

	logger := log.New(os.Stdout, "api", log.LstdFlags)
	f, err := os.Open(file)
	if err != nil {
		logger.Fatalln("failed to open file:", err.Error())
	}

	logger.Println("reading file", file)
	tsvr := NewTSVReader(f)
	go func() {
		for err := range tsvr.Error() {
			logger.Println("error while reading tsv:", err.Error())
		}
	}()

	// Create new date tree and insert every line in it.
	tree := datetree.NewTree()
	var t time.Time
	for l := range tsvr.Lines() {
		if len(l) != 2 {
			continue
		}

		t, err = time.Parse(Second, l[0])
		if err != nil {
			logger.Printf("wrong date [%s] encountered in file: %s\n", l[0], err.Error())
			continue
		}

		tree.Insert(l[1], t)
	}

	logger.Println("indexing popularity")
	tree.IndexPopularity()
	logger.Printf("%d queries processed", tree.TotalCount)

	// create handler
	var h Handler
	h.Logger = logger
	h.DateTree = tree

	// create server mux
	mux := http.NewServeMux()
	mux.Handle("/1/queries/count/", h.DateMiddleware(http.HandlerFunc(h.Count)))
	mux.Handle("/1/queries/popular/", h.DateMiddleware(http.HandlerFunc(h.Popular)))
	h.Logger.Println("server listening on port", port)
	err = http.ListenAndServe(":"+strconv.FormatUint(uint64(port), 10), mux)
	h.Logger.Fatal(err)
}
