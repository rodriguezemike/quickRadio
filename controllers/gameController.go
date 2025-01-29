package controllers

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"strconv"
	"strings"
)

type GameController struct {
	Landinglinks               []string
	ActiveGameIndex            int
	ActiveLandingLink          string
	ActiveGameDataObject       models.GameData
	ActiveGameVersesDataObject models.GameVersesData
	sweaters                   map[string]models.Sweater
	gameDataObjects            []models.GameData
	gameVersesObjects          []models.GameVersesData
	//Here we want what we will need to update the ui, filepaths, syncMaps, etc.
	activeGameDirectory string
}

func (controller *GameController) emptyActiveGameDirectory() {
	go quickio.EmptyActiveGameDirectory(controller.activeGameDirectory)
}

func (controller *GameController) GetUIDataFromFilename(teamAbbrev string, dataLabel string, defaultReturnValue string) string {
	files, _ := os.ReadDir(controller.activeGameDirectory)
	for _, f := range files {
		info, _ := f.Info()
		if strings.Contains(info.Name(), teamAbbrev) && strings.Contains(info.Name(), dataLabel) {
			return strings.Split(info.Name(), ".")[1]
		}
	}
	return defaultReturnValue
}

func (controller *GameController) getGamestateString() string {
	if controller.ActiveGameDataObject.GameState == "LIVE" || controller.ActiveGameDataObject.GameState == "CRIT" {
		if !controller.ActiveGameDataObject.Clock.InIntermission {
			return controller.ActiveGameDataObject.GameState + " - " +
				controller.ActiveGameDataObject.Venue.Default + ", " + controller.ActiveGameDataObject.VenueLocation.Default +
				" - " + "P" + strconv.Itoa(controller.ActiveGameDataObject.PeriodDescriptor.Number) + " " + controller.ActiveGameDataObject.Clock.TimeRemaining
		} else {
			return controller.ActiveGameDataObject.GameState + " - " +
				controller.ActiveGameDataObject.Venue.Default + ", " + controller.ActiveGameDataObject.VenueLocation.Default +
				" INT " + strconv.Itoa(controller.ActiveGameDataObject.PeriodDescriptor.Number) + " " + controller.ActiveGameDataObject.Clock.TimeRemaining
		}
	} else {
		if controller.ActiveGameDataObject.GameState == "FUT" {
			return controller.ActiveGameDataObject.GameState + " - " +
				controller.ActiveGameDataObject.GameDate + " - " + controller.ActiveGameDataObject.StartTimeUTC +
				controller.ActiveGameDataObject.Venue.Default + ", " + controller.ActiveGameDataObject.VenueLocation.Default
		} else {
			return controller.ActiveGameDataObject.GameState
		}
	}
}

func (controller *GameController) getTeamOnIceJson(team models.TeamOnIce) []byte {
	onIceJson, _ := json.MarshalIndent(team, "", " ")
	return onIceJson
}

func (controller *GameController) getTeamGameStats() []byte {
	tameGameStats, _ := json.MarshalIndent(controller.ActiveGameVersesDataObject.GameInfo.TeamGameStats, "", " ")
	return tameGameStats
}

func (controller *GameController) DumpGameData() {
	homeScorePath := filepath.Join(controller.activeGameDirectory, controller.ActiveGameDataObject.HomeTeam.Abbrev+"_SCORE."+strconv.Itoa(controller.ActiveGameDataObject.HomeTeam.Score))
	awayScorePath := filepath.Join(controller.activeGameDirectory, controller.ActiveGameDataObject.AwayTeam.Abbrev+"_SCORE."+strconv.Itoa(controller.ActiveGameDataObject.AwayTeam.Score))
	gameStatePath := filepath.Join(controller.activeGameDirectory, "ActiveGameState.label")
	homePlayersOnIcePath := filepath.Join(controller.activeGameDirectory, controller.ActiveGameDataObject.HomeTeam.Abbrev+"_PLAYERSONICE.json")
	awayPlayersOnIcePath := filepath.Join(controller.activeGameDirectory, controller.ActiveGameDataObject.AwayTeam.Abbrev+"_PLAYERSONICE.json")
	tameGameStatsPath := filepath.Join(controller.activeGameDirectory, "TEAMGAMESTATS.json")
	go quickio.TouchFile(homeScorePath)
	go quickio.TouchFile(awayScorePath)
	go quickio.WriteFile(gameStatePath, controller.getGamestateString())
	go quickio.WriteFile(homePlayersOnIcePath, string(controller.getTeamOnIceJson(controller.ActiveGameDataObject.Summary.IceSurface.HomeTeam)))
	go quickio.WriteFile(awayPlayersOnIcePath, string(controller.getTeamOnIceJson(controller.ActiveGameDataObject.Summary.IceSurface.AwayTeam)))
	go quickio.WriteFile(tameGameStatsPath, string(controller.getTeamGameStats()))
}

func (controller *GameController) UpdateDataObjects() {
	controller.gameDataObjects = quickio.GoGetGameDataObjectsFromLandingLinks(controller.Landinglinks)
	controller.gameVersesObjects = quickio.GoGetGameVersesDataFromLandingLinks(controller.Landinglinks)
	controller.SwitchActiveObjects(controller.ActiveGameIndex)
}

func (controller *GameController) UpdateActiveDataObjects() {
	singletonActiveLandingLink := []string{controller.ActiveLandingLink}
	controller.ActiveGameDataObject = quickio.GetGameDataObject(controller.ActiveLandingLink)
	controller.ActiveGameVersesDataObject = quickio.GoGetGameVersesDataFromLandingLinks(singletonActiveLandingLink)[0]
}

func (controller *GameController) SwitchActiveObjects(gameIndex int) {
	controller.ActiveGameDataObject = controller.gameDataObjects[gameIndex]
	controller.ActiveGameVersesDataObject = controller.gameVersesObjects[gameIndex]
	controller.ActiveGameIndex = gameIndex
}

func (controller *GameController) GetActiveRadioLink(teamAbbrev string) (string, error) {
	if controller.ActiveGameDataObject.AwayTeam.Abbrev == teamAbbrev {
		return controller.ActiveGameDataObject.AwayTeam.RadioLink, nil
	} else if controller.ActiveGameDataObject.HomeTeam.Abbrev == teamAbbrev {
		return controller.ActiveGameDataObject.HomeTeam.RadioLink, nil
	} else {
		return "", errors.New("couldnt find a radiolink from activegamedataobject")
	}
}

func NewGameController() *GameController {
	var controller GameController
	controller.Landinglinks = quickio.GetGameLandingLinks()
	controller.sweaters = quickio.GetSweaters()
	controller.gameDataObjects = quickio.GoGetGameDataObjectsFromLandingLinks(controller.Landinglinks)
	controller.gameVersesObjects = quickio.GoGetGameVersesDataFromLandingLinks(controller.Landinglinks)
	controller.activeGameDirectory = quickio.GetActiveGameDirectory()
	controller.ActiveGameIndex = 0
	return &controller
}
