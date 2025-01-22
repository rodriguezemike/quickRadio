package game

import (
	"encoding/json"
	"errors"
	"fmt"
	"quickRadio/models"
	"quickRadio/quickio"
	"quickRadio/radioErrors"
)

func GetSweaterColors() map[string][]string {
	//In the Future turn this from a map[string][]string to a container of sweater objects with colors attached and team names.
	return quickio.GetSweaterColors()
}

func GetGameLandingLinks() []string {
	linksMap := quickio.GetLinksJson()
	html := quickio.GetGameHtml(linksMap)
	gamecenterBase := fmt.Sprintf("%v", linksMap["gamecenter_api_base"])
	gamecenterLanding := fmt.Sprintf("%v", linksMap["gamecenter_api_slug"])
	gameRegex := fmt.Sprintf("%v", linksMap["game_regex"])
	landingLinks, err := quickio.GetGameLandingLinks(html, gamecenterBase, gamecenterLanding, gameRegex)
	radioErrors.ErrorCheck(err)
	return landingLinks
}

func UIGetGameDataObjects() []models.GameData {
	var gameDataObjects []models.GameData
	landingLinks := GetGameLandingLinks()
	for _, landingLink := range landingLinks {
		gameDataObjects = append(gameDataObjects, GetGameDataObject(landingLink))
	}
	return gameDataObjects
}

func UIGetGameDataObjectsAndGameLandingLinks() ([]models.GameData, []string) {
	landingLinks := GetGameLandingLinks()
	gameDataObjects := GetGameDataObjectFromLandingLinks(landingLinks)
	return gameDataObjects, landingLinks
}

func UIGetGameDataObjectMap() map[string]models.GameData {
	var gameDataMap = make(map[string]models.GameData)
	landingLinks := GetGameLandingLinks()
	for _, landingLink := range landingLinks {
		gameDataMap[landingLink] = GetGameDataObject(landingLink)
	}
	return gameDataMap
}

func GetGameDataObjectFromLandingLinks(landingLinks []string) []models.GameData {
	var gameDataObjects []models.GameData
	for _, landingLink := range landingLinks {
		gameDataObjects = append(gameDataObjects, GetGameDataObject(landingLink))
	}
	return gameDataObjects
}

func GetGameDataObject(gameLandingLink string) models.GameData {
	var gameData = &models.GameData{}
	byteValue := quickio.GetDataFromResponse(gameLandingLink)
	err := json.Unmarshal(byteValue, gameData)
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
