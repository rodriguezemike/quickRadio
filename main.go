package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func errorCheck(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getLinksJson() map[string]interface{} {
	var linksMap map[string]interface{}

	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	jsonPath := filepath.Join(dir, "assets", "links", "links.json")

	jsonFileObject, err := os.Open(jsonPath)
	errorCheck(err)
	defer jsonFileObject.Close()
	byteValue, _ := io.ReadAll(jsonFileObject)

	json.Unmarshal([]byte(byteValue), &linksMap)
	return linksMap
}

func getGameLandingLink(html string, gamecenterBase string, gamecenterLanding string, gameRegexs []string) (string, error) {
	log.Println()
	log.Println("func getGameLandingLink START")
	var gamecenterLink string
	for _, game := range gameRegexs {
		gameRegex, _ := regexp.Compile(game)
		allGames := gameRegex.FindAllString(html, -1)
		log.Println("allGames : ", allGames)
		currentDate := strings.ReplaceAll(time.Now().Format(time.DateOnly), "-", "/")
		log.Println("currentDate : ", currentDate)
		for _, possibleGame := range allGames {
			if strings.Contains(possibleGame, currentDate) {
				gamecenterLink = strings.Trim(possibleGame, "\"")
				log.Println("gamecenterLink : ", gamecenterLink)
				break
			}
		}
		if gamecenterLink != "" {
			break
		}
	}
	if gamecenterLink == "" {
		return "", errors.New("couldnt find game center link")
	}
	gameLandingLink := gamecenterBase + strings.Split(gamecenterLink, "/")[len(strings.Split(gamecenterLink, "/"))-1] + gamecenterLanding
	log.Println("gameLandingLink : ", gameLandingLink)
	log.Println("func getGameLandingLink END")
	log.Println("")
	return gameLandingLink, nil
}
func getRadioLink(resp *http.Response) (string, error) {
	return "", nil
}

func main() {
	var html string
	linksMap := getLinksJson()
	baseUrl := fmt.Sprintf("%v", linksMap["base"])
	gamecenterBase := fmt.Sprintf("%v", linksMap["gamecenter_api_base"])
	gamecenterLanding := fmt.Sprintf("%v", linksMap["gamecenter_api_slug"])
	gameRegexs := []string{fmt.Sprintf("%v", linksMap["home_game_regex"]), fmt.Sprintf("%v", linksMap["away_game_regex"])}
	sleepTimer := linksMap["load_sleep_timer"].(float64)
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate(baseUrl),
		chromedp.Sleep(time.Duration(sleepTimer)*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			rootNode, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			html, err = dom.GetOuterHTML().WithNodeID(rootNode.NodeID).Do(ctx)
			return err
		}),
	)
	errorCheck(err)
	landingLink, err := getGameLandingLink(html, gamecenterBase, gamecenterLanding, gameRegexs)
	errorCheck(err)
	resp, err := http.Get(landingLink)
	errorCheck(err)
	radioLink, err := getRadioLink(resp)
	errorCheck(err)
	log.Println(radioLink)
}
