package controllers

import (
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
	Landinglink          string
	GameIndex            int
	dataConsumed         bool
	GameDirectory        string
	HomeTeamController   *TeamController
	AwayTeamController   *TeamController
	gameDataObject       *models.GameData
	gameVersesDataObject *models.GameVersesData
}

func (controller *GameController) EmptyDirectory() {
	quickio.EmptyDirectory(controller.GameDirectory)
}

func (controller *GameController) GetGamestateString() string {
	if controller.gameDataObject.GameState == "LIVE" || controller.gameDataObject.GameState == "CRIT" {
		if controller.gameDataObject.Clock.InIntermission {
			return controller.gameDataObject.GameState + " - " +
				controller.gameDataObject.Venue.Default + ", " + controller.gameDataObject.VenueLocation.Default +
				" - " + "P" + strconv.Itoa(controller.gameDataObject.PeriodDescriptor.Number) + " " + controller.gameDataObject.Clock.TimeRemaining
		} else {
			return controller.gameDataObject.GameState + " - " +
				controller.gameDataObject.Venue.Default + ", " + controller.gameDataObject.VenueLocation.Default +
				" INT " + strconv.Itoa(controller.gameDataObject.PeriodDescriptor.Number) + " " + controller.gameDataObject.Clock.TimeRemaining
		}
	} else {
		if controller.gameDataObject.GameState == "FUT" {
			return controller.gameDataObject.GameState + " - " +
				controller.gameDataObject.GameDate + " - " + controller.gameDataObject.StartTimeUTC +
				controller.gameDataObject.Venue.Default + ", " + controller.gameDataObject.VenueLocation.Default
		} else {
			return controller.gameDataObject.GameState
		}
	}
}

func (controller *GameController) getTeamGameStats() []byte {
	tameGameStats, _ := json.MarshalIndent(controller.gameVersesDataObject.GameInfo.TeamGameStats, "", " ")
	return tameGameStats
}

func (controller *GameController) GetGamestatePath() string {
	return filepath.Join(controller.GameDirectory, "GAMESTATE.label")
}

func (controller *GameController) GetActiveGamestateFromFile() string {
	gameStatePath := controller.GetGamestatePath()
	data, fileObject := quickio.GetDataFromFile(gameStatePath)
	fileObject.Close()
	return string(data)
}

func (controller *GameController) ProduceGameData() {
	var workGroup sync.WaitGroup
	controllers := []TeamController{*controller.HomeTeamController, *controller.AwayTeamController}
	for i := 0; i < 2; i++ {
		workGroup.Add(i)
		go func(teamController TeamController) {
			defer workGroup.Done()
			quickio.TouchFile(teamController.GetScorePath())
			quickio.WriteFile(teamController.GetTeamOnIcePath(), string(teamController.getTeamOnIceJson()))
		}(controllers[i])
	}
	workGroup.Wait()
	gameStatePath := controller.GetGamestatePath()
	//teamGameStatsPath := filepath.Join(controller.GameDirectory, "TEAMGAMESTATS.json")
	quickio.WriteFile(gameStatePath, controller.GetGamestateString())
	//quickio.WriteFile(tameGameStatsPath, string(controller.getTeamGameStats()))
	controller.dataConsumed = false
}

func (controller *GameController) ConsumeGameData() {
	var workGroup sync.WaitGroup
	controllers := []TeamController{*controller.HomeTeamController, *controller.AwayTeamController}
	for i := 0; i < 2; i++ {
		workGroup.Add(i)
		go func(teamController TeamController) {
			defer workGroup.Done()
			teamController.EmptyDirectory()
		}(controllers[i])
	}
	workGroup.Wait()
	controller.dataConsumed = true
}

func CreateNewDefaultGameController() *GameController {
	var controller GameController
	controller.gameDataObject = models.CreateDefaultGameData()
	controller.gameVersesDataObject = models.CreateDefaultVersesData()
	controller.Landinglink = ""
	controller.GameIndex = 0
	controller.GameDirectory = filepath.Join(quickio.GetQuickTmpFolder(), "GAME")
	controller.HomeTeamController = CreateNewDefaultTeamController()
	controller.AwayTeamController = CreateNewDefaultTeamController()
	return &controller
}

func CreateNewGameController(landingLink string, gameIndex int) *GameController {
	var controller GameController
	gdo := quickio.GetGameDataObject(landingLink)
	gvd := quickio.GetGameVersesData(landingLink)
	sweaters := quickio.GetSweaters()
	controller.Landinglink = landingLink
	controller.GameIndex = gameIndex
	controller.GameDirectory = filepath.Join(quickio.GetQuickTmpFolder(), strconv.Itoa(gdo.Id))
	controller.HomeTeamController = CreateNewTeamController(sweaters, landingLink, &gdo, &gvd, true, controller.GameDirectory)
	controller.AwayTeamController = CreateNewTeamController(sweaters, landingLink, &gdo, &gvd, true, controller.GameDirectory)
	controller.gameDataObject = &gdo
	controller.gameVersesDataObject = &gvd
	controller.dataConsumed = false
	log.Println("gameController::NewGameController::controller ", controller)
	return &controller
}
