package scoresparser

import (
	"fmt"
	"net/http"
)

// GetGames from teh web
func GetGames() {
	resp, _ := http.Get("https://www.bbc.co.uk/sport/football/league-two/scores-fixtures/2018-08?filter=results")
	// bytes, _ := ioutil.ReadAll(resp.Body)

	fixtureMap := ParseBBCScores(resp.Body)
	fmt.Print(fixtureMap)
}
