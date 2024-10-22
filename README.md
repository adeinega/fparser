This tiny GoLang app has zero deps. Note, package [encoding/csv](https://pkg.go.dev/encoding/csv@go1.23.2) is a part of [Standard library](https://pkg.go.dev/std).

How to run it (this isn't a quite typical way to run GoLang apps)

1. `clone the repo`
2. run
`docker run -it --entrypoint sh -v ${PWD}/:/usr/src -w /usr/src golang:1.23.2-alpine3.20`
3. and then run `go run .`

The app takes two CSV files ([lookup.csv](lookup.csv) and [flow.csv](flow.csv)) and generates a report to stdout.

Any errors in lookup.csv force an app to quit.
Minor errors in flow.csv are ignored and get printed out to stdout (they go to stderr).

