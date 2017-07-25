package main

import (
    "net/http"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "net/url"
    "strings"

    "github.com/zalfonse/tagbot/common"
    "github.com/zalfonse/lumber"
)

var logger *lumber.Logger

type WikiSearchResult struct {
  Title string `json:"title"`
  Size int `json:"size"`
  WordCount int `json:"wordcount"`
  Snippet string `json:"snippet"`
  Timestamp string `json:"timestamp"`
}

type WikiSearchQuery struct {
  SearchResults []WikiSearchResult `json:"search"`
}

type WikiSearchResponse struct {
  Query WikiSearchQuery `json:"query"`
}

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
    response_body, _ := json.Marshal(common.Response{command, "success", []string{url}})
    fmt.Fprint(w, string(response_body)) // send data to client side
}

func main() {
    logger = lumber.NewLogger(lumber.TRACE)
    http.HandleFunc("/execute", execute) // set router
    http.ListenAndServe(":80", nil) // set listen port
}
