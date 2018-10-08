**Thomas Paulmyer's test for Algolia**

-The test has been made with go version go1.10.4 linux/amd64.

-The project must be copied in `$GOPATH/src/tpaulmyer/algolia` in order to work.

-To build the project, you can do a `make` and launch ./bin/api.

-The API has two parameters, `-p [uint]`, that allows you to specify the port the API listens to (default is `8080`) and `-f [string]` to specify the TSV file to read from (default
is `hn_logs.tsv`).

-The API returns an error if the date provided in the URL is invalid or if the
size parameter in the popular request is invalid.

-Don't forget to run the tests with `make test`.

-If you have any remarks, feel free to email me at `tpaulmyer@gmail.com`.
