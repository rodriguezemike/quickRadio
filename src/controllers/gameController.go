package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"strconv"
	"strings"
	"sync"
)

type GameController struct {
	Landinglinks               []string
	ActiveGameIndex            int
	ActiveLandingLink          string
	dataConsumed               bool
	ActiveGameDataObject       *models.GameData
	ActiveGameVersesDataObject *models.GameVersesData
	Sweaters                   map[string]models.Sweater
	gameDataObjects            []models.GameData
	gameVersesObjects          []models.GameVersesData
	ctx                        context.Context
	goroutineMap               *sync.Map

	//Here we want what we will need to update the ui, filepaths, syncMaps, etc.
	activeGameDirectory string
}

func (controller *GameController) GetGameDataObjects() []models.GameData {
	return controller.gameDataObjects
}
func (controller *GameController) GetGameVersesObjects() []models.GameVersesData {
	return controller.gameVersesObjects
}

func (controller *GameController) EmptyActiveGameDirectory() {
	quickio.EmptyActiveGameDirectory(controller.activeGameDirectory)
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

func (controller *GameController) getActiveGamestateString() string {
	return controller.GetGamestateString(controller.ActiveGameDataObject)
}

func (controller *GameController) GetGamestateString(gameDataObject *models.GameData) string {
	if gameDataObject.GameState == "LIVE" || gameDataObject.GameState == "CRIT" {
		if gameDataObject.Clock.InIntermission {
			return gameDataObject.GameState + " - " +
				gameDataObject.Venue.Default + ", " + gameDataObject.VenueLocation.Default +
				" - " + "P" + strconv.Itoa(gameDataObject.PeriodDescriptor.Number) + " " + controller.ActiveGameDataObject.Clock.TimeRemaining
		} else {
			return gameDataObject.GameState + " - " +
				gameDataObject.Venue.Default + ", " + gameDataObject.VenueLocation.Default +
				" INT " + strconv.Itoa(gameDataObject.PeriodDescriptor.Number) + " " + controller.ActiveGameDataObject.Clock.TimeRemaining
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

func (controller *GameController) SwitchActiveGame(index int) {
	controller.KillActiveGame()
	controller.SwitchActiveObjects(index)
	controller.RunActiveGame()
}

func (controller *GameController) getTeamOnIceJson(team models.TeamOnIce) []byte {
	onIceJson, _ := json.MarshalIndent(team, "", " ")
	return onIceJson
}

func (controller *GameController) getTeamGameStats() []byte {
	tameGameStats, _ := json.MarshalIndent(controller.ActiveGameVersesDataObject.GameInfo.TeamGameStats, "", " ")
	return tameGameStats
}

func (controller *GameController) GetActiveGamestateFromFile() string {
	gameStatePath := filepath.Join(controller.activeGameDirectory, "ACTIVEGAMESTATE.label")
	data, fileObject := quickio.GetDataFromFile(gameStatePath)
	fileObject.Close()
	return string(data)
}

func (controller *GameController) ProduceActiveGameData() {
	homeScorePath := filepath.Join(controller.activeGameDirectory, controller.ActiveGameDataObject.HomeTeam.Abbrev+"_SCORE."+strconv.Itoa(controller.ActiveGameDataObject.HomeTeam.Score))
	awayScorePath := filepath.Join(controller.activeGameDirectory, controller.ActiveGameDataObject.AwayTeam.Abbrev+"_SCORE."+strconv.Itoa(controller.ActiveGameDataObject.AwayTeam.Score))
	gameStatePath := filepath.Join(controller.activeGameDirectory, "ACTIVEGAMESTATE.label")
	homePlayersOnIcePath := filepath.Join(controller.activeGameDirectory, controller.ActiveGameDataObject.HomeTeam.Abbrev+"_PLAYERSONICE.json")
	awayPlayersOnIcePath := filepath.Join(controller.activeGameDirectory, controller.ActiveGameDataObject.AwayTeam.Abbrev+"_PLAYERSONICE.json")
	tameGameStatsPath := filepath.Join(controller.activeGameDirectory, "TEAMGAMESTATS.json")
	quickio.TouchFile(homeScorePath)
	quickio.TouchFile(awayScorePath)
	quickio.WriteFile(gameStatePath, controller.getActiveGamestateString())
	quickio.WriteFile(homePlayersOnIcePath, string(controller.getTeamOnIceJson(controller.ActiveGameDataObject.Summary.IceSurface.HomeTeam)))
	quickio.WriteFile(awayPlayersOnIcePath, string(controller.getTeamOnIceJson(controller.ActiveGameDataObject.Summary.IceSurface.AwayTeam)))
	quickio.WriteFile(tameGameStatsPath, string(controller.getTeamGameStats()))
	controller.dataConsumed = false

}
func (controller *GameController) ConsumeActiveGameData() {
	controller.EmptyActiveGameDirectory()
	controller.dataConsumed = true
}

func (controller *GameController) RunActiveGame() {
	//Here we want to ProduceGameData to be consumed by our UI.
	ctx, cancel := context.WithCancel(controller.ctx)
	defer cancel()
	controller.goroutineMap.Store(ctx, cancel)
	for {
		select {
		case <-ctx.Done():
			controller.goroutineMap.Delete(ctx)
			controller.ConsumeActiveGameData()
			return
		default:
			if controller.dataConsumed {
				controller.UpdateActiveDataObjects()
				controller.ProduceActiveGameData()
			}
		}
	}
}

func (controller *GameController) KillActiveGame() {
	controller.goroutineMap.Range(func(key, value interface{}) bool {
		callback, _ := value.(context.CancelFunc)
		callback()
		return true
	})
}

func (controller *GameController) UpdateDataObjects() {
	controller.gameDataObjects = quickio.GoGetGameDataObjectsFromLandingLinks(controller.Landinglinks)
	controller.gameVersesObjects = quickio.GoGetGameVersesDataFromLandingLinks(controller.Landinglinks)
}

func (controller *GameController) UpdateActiveDataObjects() {
	controller.UpdateDataObjects()
	controller.ActiveGameDataObject = &controller.gameDataObjects[controller.ActiveGameIndex]
	controller.ActiveGameVersesDataObject = &controller.gameVersesObjects[controller.ActiveGameIndex]
}

func (controller *GameController) SwitchActiveObjects(gameIndex int) {
	controller.ActiveGameDataObject = &controller.gameDataObjects[gameIndex]
	controller.ActiveGameVersesDataObject = &controller.gameVersesObjects[gameIndex]
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

func (controller *GameController) KillFun() {
	controller.KillActiveGame()
	controller.ConsumeActiveGameData()
	controller.dataConsumed = false
	controller.EmptyActiveGameDirectory()
}

func NewGameController() *GameController {
	var controller GameController
	controller.Landinglinks = quickio.GetGameLandingLinks()
	controller.Sweaters = quickio.GetSweaters()
	controller.gameDataObjects = quickio.GoGetGameDataObjectsFromLandingLinks(controller.Landinglinks)
	controller.gameVersesObjects = quickio.GoGetGameVersesDataFromLandingLinks(controller.Landinglinks)
	controller.activeGameDirectory = quickio.GetActiveGameDirectory()
	controller.ActiveGameIndex = 0
	controller.ActiveGameDataObject = &controller.gameDataObjects[controller.ActiveGameIndex]
	controller.ActiveGameVersesDataObject = &controller.gameVersesObjects[controller.ActiveGameIndex]
	controller.ctx = context.Background()
	controller.dataConsumed = false
	controller.goroutineMap = &sync.Map{}
	return &controller
}
