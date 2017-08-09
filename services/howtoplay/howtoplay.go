package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/zalfonse/lumber"
	"github.com/zalfonse/tagbot/common"
)

var logger *lumber.Logger

func parseSkills(s *goquery.Selection) string {
	var skillOrder []string
	s.Find(".Cell.ListCell .Item .ExtraString").Each(func(i int, sub *goquery.Selection) {
		skillOrder = append(skillOrder, sub.Text())
	})
	return strings.Join(skillOrder, " > ")
}

func parseStarter(s *goquery.Selection) string {
	var startItems []string
	s.Find(".Cell.ListCell .Item img").Each(func(i int, sub *goquery.Selection) {
		itemName, _ := sub.Attr("alt")
		startItems = append(startItems, itemName)
	})
	return strings.Join(startItems, " > ")
}

func parseCore(s *goquery.Selection) string {
	var coreItems []string
	s.Find(".Cell.ListCell .Item img").Each(func(i int, sub *goquery.Selection) {
		itemName, _ := sub.Attr("alt")
		coreItems = append(coreItems, itemName)
	})
	return strings.Join(coreItems, " > ")
}

func parseBoots(s *goquery.Selection) string {
	boots, _ := s.Find(".Cell.Single img").Attr("alt")
	return boots
}

func parseKeystone(s *goquery.Selection) string {
	itemName, _ := s.Find(".Cell.Single img").Attr("alt")
	return itemName
}

func parseRunes(s *goquery.Selection) string {
	var runes []string
	s.Find(".Cell.Single .RuneItemList").Each(func(i int, sub *goquery.Selection) {
		number := sub.Find(".Value").Text()
		label := sub.Find(".Label").Text()
		runeName := fmt.Sprintf("(%s)%s", label, number)

		runes = append(runes, runeName)
	})
	return strings.Join(runes, " | ")
}

func lookupopgg(role, champ string) []string {
	url := fmt.Sprintf("https://na.op.gg/champion/%s/statistics/%s/overview", champ, role)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return []string{"Error: [" + err.Error() + "]"}
	}
	var results []string
	selections := doc.Find(".ChampionStatsSummaryTable .Content .Row.First")
	results = append(results, "Skill Order: "+parseSkills(selections.Eq(0)))
	results = append(results, "Starter Items: "+parseStarter(selections.Eq(2)))
	results = append(results, "Core Items: "+parseCore(selections.Eq(3)))
	results = append(results, "Boots: "+parseBoots(selections.Eq(4)))
	results = append(results, "Keystone: "+parseKeystone(selections.Eq(5)))
	results = append(results, "Runes: "+parseRunes(selections.Eq(7)))

	return results
}

func execute(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var command common.Command
	var answers []string
	if err := json.Unmarshal(body, &command); err != nil {
		logger.Error("Error: [" + err.Error() + "]")
		responseBody, _ := json.Marshal(common.Response{Command: command, Type: "failure", Answers: []string{err.Error()}})
		fmt.Fprintf(w, string(responseBody))
		return
	}

	args := strings.Split(command.Args, " ")
	if len(args) < 2 {
		logger.Error("Not enough arguments: " + command.Args)
		responseBody, _ := json.Marshal(common.Response{Command: command, Type: "failure", Answers: []string{"Not enough arguments"}})
		fmt.Fprintf(w, string(responseBody))
		return
	}

	answers = lookupopgg(args[0], args[1])
	responseBody, _ := json.Marshal(common.Response{Command: command, Type: "success", Answers: answers})
	fmt.Fprintf(w, string(responseBody)) // send data to client side
}

func main() {
	logger = lumber.NewLogger(lumber.TRACE)
	http.HandleFunc("/execute", execute) // set router
	http.ListenAndServe(":80", nil)      // set listen port
}
