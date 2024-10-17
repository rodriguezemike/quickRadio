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
	team := "LAK"
	teamKeys := []string{"awayTeam", "homeTeam"}
	abvrKey := "abbrev"
	var landingMap map[string]interface{}
	var teamMap map[string]interface{}
	defer resp.Body.Close()
	byteValue, err := io.ReadAll(resp.Body)
	errorCheck(err)
	json.Unmarshal([]byte(byteValue), &landingMap)
	for _, key := range teamKeys {
		//unmarshal the inner dictionary
		if landingMap[key][abvrKey] == team {
			return landingMap[key][radioLink].string, nil
		}
	}
	return "", errors.New("Couldnt find a radio link in the landig json.")
}

// We should have a shared array here maybe? Depends on goRoutines
func downloadAudioFiles(radioLink string) {
	//Handler of downloading audio files to a temp file location for playback
	return
}
func downloadAudioFile(audioFile string) bool {
	return false
}

func playback(landingLink string, radioLink string) {
	//This can be changed to gameState = 'Completed'
	//cloudflareBase - https://d2igy0yla8zi0u.cloudfront.net/lak/20242025/
	//Radio quality
	//Example to Extract : lak-radio_192K.m3u8
	//Then we extract that https://d2igy0yla8zi0u.cloudfront.net/lak/20242025/lak-radio_192K.m3u8
	//Then we get our aac files
	//lak-radio_192K/00021/lak-radio_192K_00118.aac
	cloudflareBase := strings.Join(strings.Split(radioLink, "/")[:len(strings.Split(radioLink, "/"))-1], "/")
	radioQualityLink := getQualityStreamSlug(radioLink)
	streamLink := cloudflareBase + "/" + radioQualityLink
	for true == true {
		//While were running in the game loop
		//We to sleep main thread, really this would update the UI
		//Looking for updates in the game landing for game status, time, etc.
		//check for new undownloaded radio files
		//download those
		//Play Any AAC file were missing.
		//Delete old AAC files
		//Update our data structs.
		time.Sleep(1 * time.Second)
		fmt.Println(streamLink)
	}
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
	//This is where we start to play the radio, could get a tad interesting.
	playback(landingLink, radioLink)
}
