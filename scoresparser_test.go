package scoresparser

import (
	"log"
	"os"
	"testing"
)

func TestParseBBCScores(t *testing.T) {
	file, err := os.Open("scores.html")
	if err != nil {
		log.Fatal(err)
	}

	tables := []struct {
		compName string
		numGames int
	}{
		{"Premier League", 7},
		{"Championship", 11},
		{"League One", 12},
		{"League Two", 12},
		{"National League", 12},
	}

	fixtureMap := ParseBBCScores(file)

	for _, gameCount := range tables {
		if len(fixtureMap[gameCount.compName]) != gameCount.numGames {
			t.Errorf("%s game count wrong expected %d got %d",
				gameCount.compName, gameCount.numGames, len(fixtureMap[gameCount.compName]))
		}
	}

	g1 := Fixture{"Premier League", "Tottenham", "Liverpool", 1, 2, "FT"}
	if fixtureMap[g1.Competition][0] != g1 {
		t.Errorf("results do not match expected %v got %v", g1, fixtureMap[g1.Competition][0])
	}
	g1 = Fixture{"Championship", "Bolton", "QPR", 1, 2, "FT"}
	if fixtureMap[g1.Competition][0] != g1 {
		t.Errorf("results do not match expected %v got %v", g1, fixtureMap[g1.Competition][0])
	}
	g1 = Fixture{"League One", "Bradford", "Charlton", 0, 2, "FT"}
	if fixtureMap[g1.Competition][1] != g1 {
		t.Errorf("results do not match expected %v got %v", g1, fixtureMap[g1.Competition][0])
	}
	g1 = Fixture{"League Two", "Colchester", "Cambridge", 3, 0, "FT"}
	if fixtureMap[g1.Competition][2] != g1 {
		t.Errorf("results do not match expected %v got %v", g1, fixtureMap[g1.Competition][0])
	}
	g1 = Fixture{"National League", "Dover", "Solihull Moors", 0, 2, "FT"}
	if fixtureMap[g1.Competition][4] != g1 {
		t.Errorf("results do not match expected %v got %v", g1, fixtureMap[g1.Competition][0])
	}
}
