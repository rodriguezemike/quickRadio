package game

import (
	"encoding/json"
	"errors"
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
)

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
}
func DownloadAudioFile(audioFile string) bool {
	return false
}
