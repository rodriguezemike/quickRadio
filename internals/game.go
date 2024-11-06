package internals

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"quickRadio/models"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func GetLinksJson() map[string]interface{} {
	var linksMap map[string]interface{}
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	jsonPath := filepath.Join(dir, "assets", "links", "links.json")

	jsonFileObject, err := os.Open(jsonPath)
	ErrorCheck(err)
	defer jsonFileObject.Close()
	byteValue, _ := io.ReadAll(jsonFileObject)

	json.Unmarshal([]byte(byteValue), &linksMap)
	return linksMap
}

func GetGameLandingLink(html string, gamecenterBase string, gamecenterLanding string, gameRegexs []string) (string, error) {
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

func GetGameDataObjectFromResponse(gameLandingLink string) models.GameDataStruct {
	var gameData = &models.GameDataStruct{}
	resp, err := http.Get(gameLandingLink)
	ErrorCheck(err)
	defer resp.Body.Close()
	byteValue, err := io.ReadAll(resp.Body)
	ErrorCheck(err)
	err = json.Unmarshal(byteValue, gameData)
	ErrorCheck(err)
	return *gameData
}

func GetRadioLink(gameData models.GameDataStruct, teamAbbrev string) (string, error) {

	if gameData.AwayTeam.Abbrev == teamAbbrev {
		return gameData.AwayTeam.RadioLink, nil
	} else if gameData.HomeTeam.Abbrev == teamAbbrev {
		return gameData.HomeTeam.RadioLink, nil
	} else {
		return "", errors.New("Couldnt find a radio link in the landig json.")
	}
}

func GetQualityStreamSlug(radioLink string) (string, error) {
	return "", nil
}

func DownloadAudioFiles(radioLink string) {
	//Handler of downloading audio files to a temp file location for playback
	return
}
func DownloadAudioFile(audioFile string) bool {
	return false
}
