package controllers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"quickRadio/radioErrors"
	"strconv"
	"strings"
	"sync"
)

// Should only have 1 game in here, the Manager Controller can handle the list.
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
			if controller.gameDataObject.GameState != "" {
				return controller.gameDataObject.GameState
			} else {
				return models.DEFAULT_GAMESTATE_STRING
			}
		}
	}
}

func (controller *GameController) GetTeamGameStats() []byte {
	tameGameStats, _ := json.MarshalIndent(controller.gameVersesDataObject.GameInfo.TeamGameStats, "", " ")
	return tameGameStats
}

func (controller *GameController) GetTeamGameStatsObjects() []models.TeamGameStat {
	return controller.gameVersesDataObject.GameInfo.TeamGameStats
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

func (controller *GameController) GetGameStatFromFile(categoryName string) (int, int, int, bool) {
	//Should be a file that exists in game directory that has the infor in the file name ending with gamecategorySlider
	//Abstract this further to save a file I/O operation.
	files, _ := os.ReadDir(controller.GameDirectory)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "."+categoryName) {
			values := strings.Split(strings.Split(f.Name(), ".")[0], models.VALUE_DELIMITER)
			//Shouldnt Crash? Maybe for testing but this hould have some sort default if all the values are not there.
			//If we crash we need to have something to dump all of memories and exit gracefully avoiding mem leaks.
			homeValue, err := strconv.Atoi(values[0])
			radioErrors.ErrorLog(err)
			awayValue, err := strconv.Atoi(values[0])
			radioErrors.ErrorLog(err)
			if homeValue-awayValue == 0 {
				homeValue = homeValue / 2
				awayValue = awayValue / 2
			}
			maxValue, err := strconv.Atoi(values[2])
			radioErrors.ErrorLog(err)
			homeHandle, err := strconv.ParseBool(values[3])
			radioErrors.ErrorLog(err)
			return homeValue, awayValue, maxValue, homeHandle
		}
	}
	if strings.Contains(strings.ToLower(categoryName), "home") {
		return models.DEFAULT_WINNING_STAT_INT, models.DEFAULT_LOSTING_STAT, models.DEFAULT_TOTAL_STAT_INT, true
	} else if strings.Contains(strings.ToLower(categoryName), "tied") {
		return models.DEFAULT_WINNING_STAT_INT / 2, models.DEFAULT_LOSTING_STAT / 2, models.DEFAULT_TOTAL_STAT_INT, true
	} else {
		return models.DEFAULT_LOSTING_STAT, models.DEFAULT_WINNING_STAT_INT, models.DEFAULT_TOTAL_STAT_INT, false

	}
}

// Next season : Write a Controller abst object that game controller, team controller, gamestat and game stats controller can use
// To produce data from the Game controller which calls all other controllers produce data and then this can done in parallel with a single
// wait group at this level. this will become the standardized model for All other apps using this MVC style arch
// Produces a path to be touched holding all necessary data for the UI to update. Avoids file opening and closing operations.
func (controller *GameController) getGameStatPath(gameStat models.TeamGameStat) string {
	var homeHandle string
	awayValue, err := strconv.Atoi(gameStat.AwayValue)
	radioErrors.ErrorLog(err)
	homeValue, err := strconv.Atoi(gameStat.HomeValue)
	radioErrors.ErrorLog(err)
	maxValue := strconv.Itoa(awayValue + homeValue)
	if homeValue > awayValue {
		homeHandle = strconv.FormatBool(true)
	} else {
		homeHandle = strconv.FormatBool(false)
	}
	filename := gameStat.HomeValue + models.VALUE_DELIMITER +
		gameStat.AwayValue + models.VALUE_DELIMITER +
		maxValue + models.VALUE_DELIMITER +
		homeHandle + "." +
		gameStat.Category
	path := filepath.Join(controller.GameDirectory, filename)
	return path
}

func (controller *GameController) ProduceGameData() {
	var workGroup sync.WaitGroup
	var gameStatWorkGroup sync.WaitGroup
	controllers := []TeamController{*controller.HomeTeamController, *controller.AwayTeamController}
	teamGameStatObjects := controller.GetTeamGameStatsObjects()
	for i := 0; i < 2; i++ {
		workGroup.Add(i)
		go func(teamController TeamController) {
			defer workGroup.Done()
			quickio.TouchFile(teamController.GetScorePath())
			quickio.WriteFile(teamController.GetTeamOnIcePath(), string(teamController.getTeamOnIceJson()))
		}(controllers[i])
	}
	for i := 0; i < len(teamGameStatObjects); i++ {
		gameStatWorkGroup.Add(i)
		go func(gameStat models.TeamGameStat) {
			defer gameStatWorkGroup.Done()
			quickio.TouchFile(controller.getGameStatPath(gameStat))
		}(teamGameStatObjects[i])
	}
	workGroup.Wait()
	gameStatWorkGroup.Wait()
	gameStatePath := controller.GetGamestatePath()
	quickio.WriteFile(gameStatePath, controller.GetGamestateString())
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
