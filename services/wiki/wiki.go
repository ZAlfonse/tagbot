package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/zalfonse/lumber"
	"github.com/zalfonse/tagbot/common"
)

var logger *lumber.Logger

// WikiSearchResult ...
type WikiSearchResult struct {
	Title     string `json:"title"`
	Size      int    `json:"size"`
	WordCount int    `json:"wordcount"`
	Snippet   string `json:"snippet"`
	Timestamp string `json:"timestamp"`
}

// WikiSearchQuery ...
type WikiSearchQuery struct {
	SearchResults []WikiSearchResult `json:"search"`
}

// WikiSearchResponse ...
type WikiSearchResponse struct {
	Query WikiSearchQuery `json:"query"`
}

// WikiSearch ...
func WikiSearch(term string) string {
	logger.Info("Searching: " + term)
	if term == "" {
		return "No search term."
	}
	resp, err := http.Get("https://en.wikipedia.org/w/api.php?action=query&list=search&format=json&srsearch=" + url.QueryEscape(term))
	if err != nil {
		return "Error: [" + err.Error() + "]"
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var result WikiSearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "Error: [" + err.Error() + "]"
	}
	t := &url.URL{Path: strings.Replace(result.Query.SearchResults[0].Title, " ", "_", -1)}
	return "https://en.wikipedia.org/wiki/" + t.String()
}

func execute(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var command common.Command
	if err := json.Unmarshal(body, &command); err != nil {
		logger.Error("Error: [" + err.Error() + "]")
		return
	}
	url := WikiSearch(command.Args)
	responseBody, _ := json.Marshal(common.Response{Command: command, Type: "success", Answers: []string{url}})
	fmt.Fprint(w, string(responseBody)) // send data to client side
}

func main() {
	logger = lumber.NewLogger(lumber.TRACE)
	http.HandleFunc("/execute", execute) // set router
	http.ListenAndServe(":80", nil)      // set listen port
}
