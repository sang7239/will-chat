package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// openGraphPrefix is a prefix used for og meta properties
const openGraphPrefix = "og:"

// openGraphProps is a map of open graph property names and values
type openGraphProps map[string]string

// fetchHtml fetches html body from the passed in url
func fetchHTML(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if resp.StatusCode >= 400 {
		return nil, errors.New(resp.Status)
	}
	ctype := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ctype, "text/html") {
		return nil, err
	}
	return resp.Body, nil
}

func fetchOgProps(body io.ReadCloser) (openGraphProps, error) {
	ogp := make(openGraphProps)
	tokenizer := html.NewTokenizer(body)
	for {
		tokenType := tokenizer.Next()
		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if "head" == token.Data {
				break
			}
		}
		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()
			if "meta" == token.Data {
				var prop, content string
				for _, a := range token.Attr {
					switch a.Key {
					case "property":
						prop = a.Val
					case "content":
						content = a.Val
					}
				}
				if prop != "" && content != "" {
					if strings.HasPrefix(prop, openGraphPrefix) {
						prop = strings.TrimPrefix(prop, openGraphPrefix)
					}
					ogp[prop] = content
				}
			}
		}
	}
	return ogp, nil
}

// getPageSummary fetches a webpage and returns it's open graph properties summary
func getPageSummary(url string) (openGraphProps, error) {

	//Get the URL
	//If there was an error, return it

	//ensure that the response body stream is closed eventually

	//if the response StatusCode is >= 400
	//return an error, using the response's .Status
	//property as the error message

	//if the response's Content-Type header does not
	//start with "text/html", return an error noting
	//what the content type was and that you were
	//expecting HTML
	body, err := fetchHTML(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// tokenize and fetch the open graph properties from the url body
	ogProps, err := fetchOgProps(body)
	if err != nil {
		return nil, err
	}
	//HINTS: https://info344-s17.github.io/tutorials/tokenizing/
	//https://godoc.org/golang.org/x/net/html
	return ogProps, nil
}

// SummaryHandler fetches the URL in the `url` query string paramter, extracts
// summary information about the returned page and sends the summary properties
// to client as a JSON-encoded object.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	//get the `url` query string parameter
	//if you use r.FormValue() it will also handle cases where
	//the client did POST with `url` as a form field
	//HINT: https://golang.org/pkg/net/http/#Request.FormValue
	URL := r.FormValue("url")

	//if no `url` parameter was provided, respond with
	//an http.StatusBadRequest error and return
	//HINT: https://golang.org/pkg/net/http/#Error

	if len(URL) == 0 {
		http.Error(w, "no url parameter", http.StatusBadRequest)
		return
	}

	//call getPageSummary() passing the requested URL
	//and holding on to the returned openGraphProps map
	//(see type definition above)
	//if you get back an error, respond to the client
	//with that error and an http.StatusBadRequest code
	props, err := getPageSummary(URL)
	if err != nil {
		http.Error(w, "request error: "+err.Error(), http.StatusBadRequest)
		return
	}

	//otherwise, respond by writing the openGrahProps
	//map as a JSON-encoded object
	//add the following headers to the response before
	//you write the JSON-encoded object:
	//   Content-Type: application/json; charset=utf-8
	//this tells the client that you are sending it JSON

	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(props); err != nil {
		http.Error(w, "json encoding error: "+err.Error(), http.StatusInternalServerError)
	}
}
