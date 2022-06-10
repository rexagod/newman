//go:build generate
// +build generate

package main

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
	"sort"
	"strings"
)

const ENDPOINT_DOWN = "ENDPOINT_DOWN"

//go:generate cp ../artefacts/quotes.json ./quotes.json
//go:embed quotes.json
var marshalledQuotes []byte

func fetchAndSanitizeQuotes() ([]string, error) {
	var err error
	var quotes []string
	type endpoint struct {
		url        string
		entrypoint string
		allowed    map[string]string
		lessFn     func([]string) func(int, int) bool
		separator  string
	}
	endpoints := []endpoint{
		{
			url:        "https://seinfeld-quotes.herokuapp.com/quotes",
			entrypoint: "quotes",
			allowed: map[string]string{
				"author": "",
				"quote":  "",
			},
			separator: ": ",
			lessFn: func(t []string) func(int, int) bool {
				return func(i, j int) bool {
					return len(t[i]) > len(t[j])
				}
			},
		},
	}
	for _, endpoint := range endpoints {
		var response *http.Response
		response, err = http.Get(endpoint.url)
		if err != nil || response.StatusCode != 200 {
			panic(errors.New(ENDPOINT_DOWN + ": " + endpoint.url))
		}
		var body []byte
		body, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %v", err)
		}
		err = response.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to close response body: %v", err)
		}
		var unmarshalledQuotes map[string][]map[string]string
		err = json.Unmarshal(body, &unmarshalledQuotes)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
		}
		v := unmarshalledQuotes[endpoint.entrypoint]
		/*
			the size of the array below is the number of quotes received
			the size of every map instance in that array is the number of fields in the response
			the former will always be exponentially larger than the latter
			the complexity of this loop is linear: O(mn) =~ O(n); where m << n, but m >= 1.
		*/
		// traverse the array of maps
		for _, vv := range v {
			// traverse the map
			var t []string
			for k, vvv := range vv {
				if _, ok := endpoint.allowed[k]; ok {
					endpoint.allowed[k] = vvv
					t = append(t, k)
				}
			}
			sort.Slice(t, endpoint.lessFn(t))
			// traverse endpoints.allowed (negligible complexity)
			var quote string
			for _, vvvv := range t {
				quote += endpoint.allowed[vvvv] + endpoint.separator
			}
			quotes = append(quotes, quote[:len(quote)-len(endpoint.separator)])
		}
	}
	return quotes, nil
}

func compareQuotes(currentQuotes, fetchedQuotes []string) ([]string, error) {
	if len(fetchedQuotes) <= len(currentQuotes) {
		return nil, nil
	}
	m := make(map[string]bool)
	for _, v := range fetchedQuotes {
		m[v] = true
	}
	for _, v := range currentQuotes {
		m[v] = true
	}
	var newQuotes []string
	for k := range m {
		newQuotes = append(newQuotes, k)
	}
	return newQuotes, nil
}

func updateQuotes() error {
	var err error
	type quotes struct {
		Quotes []string `json:"quotes"`
	}
	var q quotes
	err = json.Unmarshal(marshalledQuotes, &q)
	if err != nil {
		return fmt.Errorf("failed to unmarshal quotes: %v", err)
	}
	defer func() {
		r := recover()
		switch r.(type) {
		case string:
			if strings.HasPrefix(r.(string), ENDPOINT_DOWN) {
				// try again
				err := updateQuotes()
				if err != nil {
					klog.Fatalf("failed to update quotes: %v", err)
				}
			}
		default:
			if r != nil {
				panic(r)
			}
		}
	}()
	var fetchedQuotes []string
	fetchedQuotes, err = fetchAndSanitizeQuotes()
	if err != nil {
		return fmt.Errorf("failed to fetch quotes: %v", err)
	}
	var totalQuotes []string
	totalQuotes, err = compareQuotes(q.Quotes, fetchedQuotes)
	if err != nil {
		return fmt.Errorf("failed to compare quotes: %v", err)
	}
	q.Quotes = totalQuotes
	var marshalledQ []byte
	marshalledQ, err = json.Marshal(q)
	if err != nil {
		return fmt.Errorf("failed to marshal quotes: %v", err)
	}
	var ownerReadWritePermission fs.FileMode = 0600
	err = ioutil.WriteFile("artefacts/quotes.json", marshalledQ, ownerReadWritePermission)
	if err != nil {
		return fmt.Errorf("failed to write quotes: %v", err)
	}
	return nil
}

func main() {
	err := updateQuotes()
	if err != nil {
		klog.Fatalf("failed to update quotes: %v", err)
	}
}
