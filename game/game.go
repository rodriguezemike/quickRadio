package game

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

func UIGetGameDataObjects() []models.GameDataStruct {
	var gameDataObjects []models.GameDataStruct
	landingLinks := UIGetGameLandingLinks()
	for _, landingLink := range landingLinks {
		gameDataObjects = append(gameDataObjects, GetGameDataObjectFromResponse(landingLink))
	}
	return gameDataObjects
}

func UIGetGameDataObjectMap() map[string]models.GameDataStruct {
	var gameDataMap = make(map[string]models.GameDataStruct)
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
	log.Println("allGames : ", allGames)
	for _, possibleGame := range allGames {
		if strings.Contains(possibleGame, currentDate) {
			gamecenterLink := strings.Trim(possibleGame, "\"")
			log.Println(strings.Split(gamecenterLink, "/"))
			log.Println(gamecenterLanding)
			landingLink := gamecenterBase + strings.Split(gamecenterLink, "/")[len(strings.Split(gamecenterLink, "/"))-1] + gamecenterLanding
			gameLandingLinks = append(gameLandingLinks, landingLink)
			log.Println(landingLink)
		}
	}
	if len(gameLandingLinks) == 0 {
		return nil, errors.New("couldnt find any game for today.")
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
		log.Println("game: ", game)
		log.Println("New Game :", newGame)
		gameRegex, _ := regexp.Compile(newGame)
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

func GetGameDataObjectFromResponse(gameLandingLink string) models.GameDataStruct {
	var gameData = &models.GameDataStruct{}
	resp, err := http.Get(gameLandingLink)
	radioErrors.ErrorCheck(err)
	defer resp.Body.Close()
	byteValue, err := io.ReadAll(resp.Body)
	radioErrors.ErrorCheck(err)
	err = json.Unmarshal(byteValue, gameData)
	radioErrors.ErrorCheck(err)
	return *gameData
}

func GetRadioLink(gameData models.GameDataStruct, teamAbbrev string) (string, error) {
	if gameData.AwayTeam.Abbrev == teamAbbrev {
		return gameData.AwayTeam.RadioLink, nil
	} else if gameData.HomeTeam.Abbrev == teamAbbrev {
		return gameData.HomeTeam.RadioLink, nil
	} else {
		return "", errors.New("couldnt find a radio link in the landing json")
	}
}

func GetQualityStreamSlugFromResponse(radioLink string, audioQuality string) string {
	resp, err := http.Get(radioLink)
	radioErrors.ErrorCheck(err)
	defer resp.Body.Close()
	//This will just get the file from download so this is a placeholder for now.
	byteValue, err := io.ReadAll(resp.Body)
	radioErrors.ErrorCheck(err)
	audioQualitySlug, err := GetQualityStreamSlug(string(byteValue), audioQuality)
	radioErrors.ErrorCheck(err)
	return audioQualitySlug
}

func GetQualityStreamSlug(m3uContents string, audioQuality string) (string, error) {
	for _, line := range strings.Split(m3uContents, "\n") {
		if strings.Contains(line, audioQuality) && strings.Contains(line, ".m3u8") {
			return line, nil
		}
	}
	return "", errors.New("couldnt find audio quality string")
}

func GetAudioFiles(m3uContents string) ([]string, error) {
	var audioFiles []string
	for _, line := range strings.Split(m3uContents, "\n") {
		if !strings.Contains(line, "#") && strings.Contains(line, ".aac") {
			audioFiles = append(audioFiles, line)
		}
	}
	if len(audioFiles) == 0 {
		return nil, errors.New("couldnt find audio files")
	}
	return audioFiles, nil
}

func DownloadAudioFiles(radioLink string) {
	//Handler of downloading audio files to a temp file location for playback
	//Should cross platform, we want to find the tmp location for each of the following
	//windows, linux, default will be unix
	//https://pkg.go.dev/runtime#GOOS is what will be used
	//C:\Users\AppData\Local\Temp is for windows
	// /tmp/ for linux
	// Assume /tmp/ for everything else.
	//Subdir will be quickRadio
	//We can store em by game and build that string from the gamedata object.
}
