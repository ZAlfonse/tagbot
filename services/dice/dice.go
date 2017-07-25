package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zalfonse/lumber"
	"github.com/zalfonse/tagbot/common"
)

var logger *lumber.Logger

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func roll(count, max int) []string {
	rand.Seed(time.Now().Unix())
	var results []string
	for i := 0; i < count && i < 20; i++ {
		results = append(results, strconv.Itoa(random(1, max+1)))
	}
	return results
}

func execute(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	var command common.Command
	if err := json.Unmarshal(body, &command); err != nil {
		logger.Error("Error: [" + err.Error() + "]")
		fmt.Fprintf(w, "Couldn't roll the dice.")
		return
	}

	args := strings.Split(command.Args, "d")
	var results []string
	if len(args) == 2 {
		count, err := strconv.Atoi(args[0])
		if err != nil {
			results = append(results, "Invalid number of rolls.")
		}
		sides, err := strconv.Atoi(args[1])
		if err != nil {
			results = append(results, "Invalid number of dice sides.")
		}
		if len(results) > 0 {
			response_body, _ := json.Marshal(common.Response{command, "failure", results})
			fmt.Fprintf(w, string(response_body))
			return
		}
		results = roll(count, sides)
	} else {
		results = roll(1, 20)
	}
	response_body, _ := json.Marshal(common.Response{command, "success", results})
	fmt.Fprintf(w, string(response_body)) // send data to client side
}

func main() {
	logger = lumber.NewLogger(lumber.TRACE)
	http.HandleFunc("/execute", execute) // set router
	http.ListenAndServe(":80", nil)      // set listen port
}
