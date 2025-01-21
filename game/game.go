package game

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/radioErrors"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
)

func GetGameHtml(linksMap map[string]interface{}) string {
	var html string
	sleepTimer := linksMap["load_sleep_timer"].(float64)
	baseUrl := fmt.Sprintf("%v", linksMap["base"])

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
	radioErrors.ErrorCheck(err)
	return html
}

func GetSweaterColors() map[string][]string {
	sweaterColors := make(map[string][]string)
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	path := filepath.Join(dir, "assets", "teams", "sweater_colors.txt")

	fileObject, err := os.Open(path)
	radioErrors.ErrorCheck(err)
	defer fileObject.Close()
	scanner := bufio.NewScanner(fileObject)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lineColors := strings.Split(scanner.Text(), ";")
		sweaterColors[lineColors[0]] = lineColors[1:]
	}

	return sweaterColors
}

func GetLinksJson() map[string]interface{} {
	var linksMap map[string]interface{}
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	jsonPath := filepath.Join(dir, "assets", "links", "links.json")

	jsonFileObject, err := os.Open(jsonPath)
	radioErrors.ErrorCheck(err)
	defer jsonFileObject.Close()
	byteValue, _ := io.ReadAll(jsonFileObject)

	json.Unmarshal([]byte(byteValue), &linksMap)
	return linksMap
}

func UIGetGameLandingLinks() []string {
	linksMap := GetLinksJson()
	html := GetGameHtml(linksMap)
	gamecenterBase := fmt.Sprintf("%v", linksMap["gamecenter_api_base"])
	gamecenterLanding := fmt.Sprintf("%v", linksMap["gamecenter_api_slug"])
	gameRegex := fmt.Sprintf("%v", linksMap["game_regex"])
	landingLinks, err := GetGameLandingLinks(html, gamecenterBase, gamecenterLanding, gameRegex)
	radioErrors.ErrorCheck(err)
	return landingLinks
}

func UIGetGameDataObjects() []models.GameData {
	var gameDataObjects []models.GameData
	landingLinks := UIGetGameLandingLinks()
	for _, landingLink := range landingLinks {
		gameDataObjects = append(gameDataObjects, GetGameDataObjectFromResponse(landingLink))
	}
	return gameDataObjects
}

func UIGetGameDataObjectMap() map[string]models.GameData {
	var gameDataMap = make(map[string]models.GameData)
	landingLinks := UIGetGameLandingLinks()
	for _, landingLink := range landingLinks {
		gameDataMap[landingLink] = GetGameDataObjectFromResponse(landingLink)
	}
	return gameDataMap
}

func GetGameLandingLinks(html string, gamecenterBase string, gamecenterLanding string, gameRegex string) ([]string, error) {
	var gameLandingLinks []string
	gameRegexObject, _ := regexp.Compile(gameRegex)
	allGames := gameRegexObject.FindAllString(html, -1)
	currentDate := strings.ReplaceAll(time.Now().Format(time.DateOnly), "-", "/")
	for _, possibleGame := range allGames {
		if strings.Contains(possibleGame, currentDate) {
			gamecenterLink := strings.Trim(possibleGame, "\"")
			landingLink := gamecenterBase + strings.Split(gamecenterLink, "/")[len(strings.Split(gamecenterLink, "/"))-1] + gamecenterLanding
			gameLandingLinks = append(gameLandingLinks, landingLink)
		}
	}
	if len(gameLandingLinks) == 0 {
		return nil, errors.New("couldnt find any game for today")
	}
	return gameLandingLinks, nil
}

func GetGameLandingLink(html string, gamecenterBase string, gamecenterLanding string, gameRegexs []string, teamAbbrev string) (string, error) {
	log.Println()
	log.Println("func getGameLandingLink START")
	var gamecenterLink string
	for _, game := range gameRegexs {
		newGame := strings.Replace(game, "TEAMABBREV", strings.ToLower(teamAbbrev), -1)
		log.Println(teamAbbrev)
		gameRegex, _ := regexp.Compile(newGame)
		allGames := gameRegex.FindAllString(html, -1)
		currentDate := strings.ReplaceAll(time.Now().Format(time.DateOnly), "-", "/")
		for _, possibleGame := range allGames {
			if strings.Contains(possibleGame, currentDate) {
				gamecenterLink = strings.Trim(possibleGame, "\"")
				log.Println("gamecenterLink : ", gamecenterLink)
				break
			}
		}
	}
	if gamecenterLink == "" {
		return "", errors.New("couldnt find game center link")
	}
	gameLandingLink := gamecenterBase + strings.Split(gamecenterLink, "/")[len(strings.Split(gamecenterLink, "/"))-1] + gamecenterLanding
	return gameLandingLink, nil
}

func GetGameDataObjectFromResponse(gameLandingLink string) models.GameData {
	var gameData = &models.GameData{}
	resp, err := http.Get(gameLandingLink)
	log.Println(gameLandingLink)
	radioErrors.ErrorCheck(err)
	defer resp.Body.Close()
	byteValue, err := io.ReadAll(resp.Body)
	radioErrors.ErrorCheck(err)
	err = json.Unmarshal(byteValue, gameData)
	radioErrors.ErrorCheck(err)
	return *gameData
}

func GetRadioLink(gameData models.GameData, teamAbbrev string) (string, error) {
	if gameData.AwayTeam.Abbrev == teamAbbrev {
		return gameData.AwayTeam.RadioLink, nil
	} else if gameData.HomeTeam.Abbrev == teamAbbrev {
		return gameData.HomeTeam.RadioLink, nil
	} else {
		return "", errors.New("couldnt find a radio link in the landing json")
	}
}
