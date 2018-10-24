package scoresparser

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// Fixture game result
type Fixture struct {
	Competition   string
	HomeTeam      string
	AwayTeam      string
	HomeTeamScore int
	AwayTeamScore int
	Aside         string
}

// CompetitionNames interesting competitions
var CompetitionNames = map[string]bool{
	"Premier League":  true,
	"Championship":    true,
	"League One":      true,
	"League Two":      true,
	"National League": true,
}

func panicClass(node *html.Node, className string) *html.Node {
	foundClassAttr := false
	for _, a := range node.Attr {
		if a.Key == "class" && strings.Contains(a.Val, className) {
			foundClassAttr = true
			break
		}
	}

	if !foundClassAttr {
		panic("node does not have class: " + className)
	}

	return node
}

func panicTag(node *html.Node, tag string) *html.Node {
	if node.Type != html.ElementNode || node.Data != tag {
		panic("node [" + node.Data + "] is not an element of tag " + tag)
	}
	return node
}

func parseTeamScore(node *html.Node) (string, int) {
	teamName := panicTag(
		panicClass(node.FirstChild, "sp-c-fixture__team-name").FirstChild.FirstChild,
		"abbr").FirstChild.Data
	score, err := strconv.Atoi(panicClass(
		panicClass(node.FirstChild.NextSibling, "sp-c-fixture__block").FirstChild,
		"sp-c-fixture__number").FirstChild.Data)

	if err != nil {
		panic("error parsing score for [" + teamName + "] " + err.Error())
	}

	return teamName, score
}

func parseFixture(competition string, fixtureNode *html.Node, asideNode *html.Node) Fixture {

	var fixture Fixture
	fixture.Competition = competition

	comment := asideNode.FirstChild.FirstChild.NextSibling.FirstChild.NextSibling
	if comment.Type != html.ElementNode || comment.Data != "abbr" {
		panic("not an aside abbr")
	}
	fixture.Aside = comment.FirstChild.Data

	home := panicClass(fixtureNode.FirstChild, "sp-c-fixture__team--home")
	away := panicClass(fixtureNode.LastChild, "sp-c-fixture__team--away")
	fixture.HomeTeam, fixture.HomeTeamScore = parseTeamScore(home)
	fixture.AwayTeam, fixture.AwayTeamScore = parseTeamScore(away)

	fmt.Println(fixture)
	return fixture
}

func parseGame(competition string, n *html.Node) Fixture {
	if n.Type != html.ElementNode || n.Data != "li" {
		panic("not an list of matches")
	}

	article := n.FirstChild.FirstChild
	if article.Type != html.ElementNode || article.Data != "article" {
		panic("not an article")
	}

	fixture := article.FirstChild
	aside := fixture.NextSibling

	return parseFixture(competition, fixture, aside)
}

func parseComp(compName string, n *html.Node) []Fixture {
	var games []Fixture
	if n.Type != html.ElementNode || n.Data != "ul" {
		panic("not an list of matches")
	}

	for item := n.FirstChild; item != nil; item = item.NextSibling {
		games = append(games, parseGame(compName, item))
	}

	return games
}

func parseResults(n *html.Node, fixtureMap map[string][]Fixture) {
	nextNode := n.FirstChild
	foundLeague := false

	if n.Type == html.ElementNode && n.Data == "h3" {
		for _, a := range n.Attr {
			if a.Key == "class" && strings.Contains(a.Val, "sp-c-match-list-heading") {
				competition := n.FirstChild.Data
				_, present := CompetitionNames[competition]
				if present {
					fixtures := parseComp(competition, n.NextSibling)
					fixtureMap[competition] = fixtures
				}
				foundLeague = true
			}
		}
	}

	if foundLeague {
		return
	}
	for c := nextNode; c != nil; c = c.NextSibling {
		parseResults(c, fixtureMap)
	}
}

// ParseBBCScores parse matches from the BBC sports website
func ParseBBCScores(htmlPage io.Reader) map[string][]Fixture {

	node, err := html.Parse(htmlPage)
	if err != nil {
		log.Fatal(err)
	}

	compFixtures := make(map[string][]Fixture)
	parseResults(node, compFixtures)
	return compFixtures
}
