package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const fieldDstPort = 6
const fieldProtocol = 7

type tag struct {
	dstPort  int
	protocol int
	name     string
	count    int
}

// https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
var ianaProtocolNumbers = map[string]int{
	"icmp": 1,
	"tcp":  6,
	"udp":  17,
	"sctp": 132,
}

func main() {
	lookupFile, err := os.Open("lookup.csv") // Might need to have additional restrictions for prod use.
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer lookupFile.Close()

	tags, err := buildLookupTable(lookupFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	flowFile, err := os.Open("flow.csv") // May need to have additional restrictions for prod use.
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer flowFile.Close()

	m, untagged := stat(flowFile, tags)

	fmt.Fprintln(os.Stdout, "Tag | Count")
	for name, val := range m {
		fmt.Fprintf(os.Stdout, "%s %d\n", name, val)
	}
	fmt.Fprintf(os.Stdout, "untagged = %d\n", untagged)

	fmt.Fprintln(os.Stdout, "Port | Protocol | Count")
	for _, tag := range tags {
		if tag.count > 0 {
			fmt.Fprintf(os.Stdout, "%d %d %d\n", tag.dstPort, tag.protocol, tag.count)
		}
	}
}

func buildLookupTable(r io.Reader) ([]tag, error) {
	parser := csv.NewReader(r)
	parser.TrimLeadingSpace = true
	parser.FieldsPerRecord = 3
	parser.Comment = '#'

	records, err := parser.ReadAll()
	if err != nil {
		return nil, err
	}

	a := make([]tag, len(records))
	for i, record := range records {
		port, err := strconv.Atoi(record[0])
		if err != nil {
			return nil, err
		}

		if port < 0 || port > 65535 {
			return nil, fmt.Errorf("invalid port %d", port)
		}

		if v, exists := ianaProtocolNumbers[strings.ToLower(record[1])]; !exists {
			return nil, fmt.Errorf("invalid protocol number: %s", record[1])
		} else {
			a[i].protocol = v
		}

		a[i].dstPort = port
		a[i].name = strings.ToLower(record[2]) // Tag names are case-insensitive.
	}
	return a, nil
}

func stat(r io.Reader, tags []tag) (map[string]int, int) {
	parser := csv.NewReader(r)
	parser.Comma = ' ' // flow.csv doesn't use the comma as the field delimiter.
	parser.TrimLeadingSpace = true

	m := map[string]int{}
	untagged := 0

	for {
		record, err := parser.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue // Should be OK to be a bit more resilient (TBD).
		}

		port, err := strconv.Atoi(record[fieldDstPort])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue //  Should be OK to be a bit more resilient (TBD).
		}

		protocol, err := strconv.Atoi(record[fieldProtocol])
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue //  Should be OK to be a bit more resilient (TBD).
		}

		marked := false
		for i, t := range tags {
			if t.dstPort == port && t.protocol == protocol {
				// For Name | Count
				if _, exists := m[t.name]; !exists {
					m[t.name] = 1
				} else {
					m[t.name]++
				}

				// For Port | Protocol | Count
				pTag := &tags[i]
				(*pTag).count++

				marked = true
			}
		}

		if !marked {
			untagged++
		}
	}
	return m, untagged
}
