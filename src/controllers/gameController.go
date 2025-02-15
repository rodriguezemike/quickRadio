package controllers

import (
	"context"
	"encoding/json"
	"log"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"strconv"
	"sync"
)

//Should only have 1 game in here, the Manager Controller can handle the list.

type GameController struct {
	Landinglink        string
	Radiolink          string
	GameIndex          int
	dataConsumed       bool
	GameDirectory      string
	HomeTeamController *TeamController
	AwayTeamController *TeamController
	gameDataObject     *models.GameData
	gameVersesData     *models.GameVersesData
}

func (controller *GameController) EmptyDirectory() {
	quickio.EmptyDirectory(controller.GameDirectory)
}

func (controller *GameController) getActiveGamestateString() string {
	return controller.GetGamestateString(controller.gameDataObject)
}

func (controller *GameController) GetGamestateString(gameDataObject *models.GameData) string {
	if gameDataObject.GameState == "LIVE" || gameDataObject.GameState == "CRIT" {
		if gameDataObject.Clock.InIntermission {
			return gameDataObject.GameState + " - " +
				gameDataObject.Venue.Default + ", " + gameDataObject.VenueLocation.Default +
				" - " + "P" + strconv.Itoa(gameDataObject.PeriodDescriptor.Number) + " " + controller.gameDataObject.Clock.TimeRemaining
		} else {
			return gameDataObject.GameState + " - " +
				gameDataObject.Venue.Default + ", " + gameDataObject.VenueLocation.Default +
				" INT " + strconv.Itoa(gameDataObject.PeriodDescriptor.Number) + " " + controller.gameDataObject.Clock.TimeRemaining
		}
	} else {
		if gameDataObject.GameState == "FUT" {
			return gameDataObject.GameState + " - " +
				gameDataObject.GameDate + " - " + gameDataObject.StartTimeUTC +
				gameDataObject.Venue.Default + ", " + gameDataObject.VenueLocation.Default
		} else {
			return gameDataObject.GameState
		}
	}
}

func (controller *GameController) getTeamOnIceJson(team models.TeamOnIce) []byte {
	onIceJson, _ := json.MarshalIndent(team, "", " ")
	return onIceJson
}

func (controller *GameController) getTeamGameStats() []byte {
	tameGameStats, _ := json.MarshalIndent(controller.gameVersesData.GameInfo.TeamGameStats, "", " ")
	return tameGameStats
}

func (controller *GameController) GetActiveGamestateFromFile() string {
	gameStatePath := filepath.Join(controller.GameDirectory, "ACTIVEGAMESTATE.label")
	data, fileObject := quickio.GetDataFromFile(gameStatePath)
	fileObject.Close()
	return string(data)
}

func (controller *GameController) ProduceGameData() {
	homeScorePath := filepath.Join(controller.GameDirectory, controller.gameDataObject.HomeTeam.Abbrev+"_SCORE."+strconv.Itoa(controller.gameDataObject.HomeTeam.Score))
	awayScorePath := filepath.Join(controller.GameDirectory, controller.gameDataObject.AwayTeam.Abbrev+"_SCORE."+strconv.Itoa(controller.gameDataObject.AwayTeam.Score))
	gameStatePath := filepath.Join(controller.GameDirectory, "ACTIVEGAMESTATE.label")
	homePlayersOnIcePath := filepath.Join(controller.GameDirectory, controller.gameDataObject.HomeTeam.Abbrev+"_PLAYERSONICE.json")
	awayPlayersOnIcePath := filepath.Join(controller.GameDirectory, controller.gameDataObject.AwayTeam.Abbrev+"_PLAYERSONICE.json")
	tameGameStatsPath := filepath.Join(controller.GameDirectory, "TEAMGAMESTATS.json")
	quickio.TouchFile(homeScorePath)
	quickio.TouchFile(awayScorePath)
	quickio.WriteFile(gameStatePath, controller.getActiveGamestateString())
	quickio.WriteFile(homePlayersOnIcePath, string(controller.getTeamOnIceJson(controller.gameDataObject.Summary.IceSurface.HomeTeam)))
	quickio.WriteFile(awayPlayersOnIcePath, string(controller.getTeamOnIceJson(controller.gameDataObject.Summary.IceSurface.AwayTeam)))
	quickio.WriteFile(tameGameStatsPath, string(controller.getTeamGameStats()))
	controller.dataConsumed = false

}
func (controller *GameController) ConsumeGameData() {
	controller.EmptyDirectory(controller.GameDirectory)
	controller.dataConsumed = true
}

func NewGameController() *GameController {
	var controller GameController
	controller.Landinglinks = quickio.GetGameLandingLinks()
	controller.Sweaters = quickio.GetSweaters()
	controller.gameDataObjects = quickio.GoGetGameDataObjectsFromLandingLinks(controller.Landinglinks)
	controller.gameVersesObjects = quickio.GoGetGameVersesDataFromLandingLinks(controller.Landinglinks)
	controller.activeGameDirectory = quickio.GetActiveGameDirectory()
	controller.ActiveGameIndex = 0
	if len(controller.gameDataObjects) > 0 {
		controller.ActiveGameDataObject = &controller.gameDataObjects[controller.ActiveGameIndex]
		controller.ActiveGameVersesDataObject = &controller.gameVersesObjects[controller.ActiveGameIndex]
	} else {
		controller.ActiveGameDataObject = &models.GameData{}
		controller.ActiveGameVersesDataObject = &models.GameVersesData{}
	}

	controller.ctx = context.Background()
	controller.dataConsumed = false
	controller.goroutineMap = &sync.Map{}
	log.Println("gameController::NewGameController::controller ", controller)
	return &controller
}
