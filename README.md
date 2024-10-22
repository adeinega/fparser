This tiny GoLang app has zero deps. Note, the [encoding/csv](https://pkg.go.dev/encoding/csv@go1.23.2) package it uses is a part of [Standard library](https://pkg.go.dev/std).

How to run it
1. clone the repo
2. cd fparser/
3. run `docker run -it --entrypoint sh -v ${PWD}/:/usr/src -w /usr/src golang:1.23.2-alpine3.20`, and finally
4. run `go run .`

The app process two files ([lookup.csv](lookup.csv) and [flow.csv](flow.csv)), and generates a report to stdout. Any errors in lookup.csv force an app to quit. Minor errors in flow.csv are ignored and get printed out to stderr.

#### Improvements
1. instead of defining "static" rules for tags in lookup.csv, the app could relatively easily be extended to use the Berkeley Packet Filter, or an equivalent expression language (EL).
2. the app could be extended to work in two modes for handling errors in flow.csv, for example, _strict_, _relaxed_ and _warn_.