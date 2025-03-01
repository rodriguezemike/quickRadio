package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"quickRadio/radioErrors"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Should only have 1 game in here, the Manager Controller can handle the list.
type GameController struct {
	Landinglink          string
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

func (controller *GameController) IsLive() bool {
	return controller.gameDataObject.GameState == "LIVE" || controller.gameDataObject.GameState == "CRIT"
}

func (controller *GameController) IsFuture() bool {
	return controller.gameDataObject.GameState == "FUT"
}

func (controller *GameController) IsPregame() bool {
	return controller.gameDataObject.GameState == "PRE"
}

func (controller *GameController) IsDone() bool {
	return controller.gameDataObject.GameState == "FINAL" || controller.gameDataObject.GameState == "OFF"
}

func (controller *GameController) GetGamestateString() string {
	if controller.gameDataObject.GameState == "LIVE" || controller.gameDataObject.GameState == "CRIT" {
		if !controller.gameDataObject.Clock.InIntermission {
			return controller.gameDataObject.GameState + " - " +
				controller.gameDataObject.Venue.Default + ", " + controller.gameDataObject.VenueLocation.Default +
				"\n" + "P" + strconv.Itoa(controller.gameDataObject.PeriodDescriptor.Number) + " " + controller.gameDataObject.Clock.TimeRemaining
		} else {
			return controller.gameDataObject.GameState + " - " +
				controller.gameDataObject.Venue.Default + ", " + controller.gameDataObject.VenueLocation.Default +
				" INT " + strconv.Itoa(controller.gameDataObject.PeriodDescriptor.Number) + " " + controller.gameDataObject.Clock.TimeRemaining
		}
	} else {
		if controller.IsFuture() || controller.IsPregame() {
			startTime, err := time.Parse(time.RFC3339, controller.gameDataObject.StartTimeUTC)
			radioErrors.ErrorLog(err)
			return "Upcoming" + "\n" + "Game Time: " + controller.gameDataObject.GameDate + " @ " + fmt.Sprint(startTime.UTC().Local().Format("03:04PM")) + "\n" +
				"Live from " + controller.gameDataObject.Venue.Default + ", " + controller.gameDataObject.VenueLocation.Default
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
	teamGameStats, _ := json.MarshalIndent(controller.gameVersesDataObject.TeamGameStats, "", " ")
	return teamGameStats
}

func (controller *GameController) GetTeamSeasonStat() []byte {
	teamSeasonStats, _ := json.MarshalIndent(controller.gameVersesDataObject.TeamSeasonStats, "", " ")
	return teamSeasonStats
}

func (controller *GameController) GetSeriesWinsStats() []byte {
	seriesWins, _ := json.MarshalIndent(controller.gameVersesDataObject.SeasonSeriesWins, "", " ")
	return seriesWins
}

func (controller *GameController) GetTeamGameStatsObjects() []models.TeamGameStat {
	return controller.gameVersesDataObject.TeamGameStats
}

func (controller *GameController) GetTeamSeasonStatObject() models.TeamSeasonStats {
	return controller.gameVersesDataObject.TeamSeasonStats
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

func (controller *GameController) getDefaultValueFromObjectName(objectName string) string {
	if strings.Contains(strings.ToLower(objectName), models.DEFAULT_HOME_PREFIX) {
		if strings.Contains(strings.Split(strings.ToLower(objectName), models.VALUE_DELIMITER)[1], models.DEFAULT_HOME_PREFIX) {
			return models.DEFAULT_WINNING_STAT
		} else {
			return models.DEFAULT_LOSING_STAT
		}
	} else if strings.Contains(strings.ToLower(objectName), models.DEFAULT_TIED_PREFIX) {
		return models.DEFAULT_WINNING_STAT
	} else if strings.Contains(strings.ToLower(objectName), models.DEFAULT_AWAY_PREFIX) {
		if strings.Contains(strings.Split(strings.ToLower(objectName), models.VALUE_DELIMITER)[1], models.DEFAULT_AWAY_PREFIX) {
			return models.DEFAULT_WINNING_STAT
		} else {
			return models.DEFAULT_LOSING_STAT
		}
	} else {
		return models.DEFAULT_LOSING_STAT
	}
}

func (controller *GameController) GetTextFromObjectNameFilepath(objectName string) string {
	log.Println(objectName)
	files, _ := os.ReadDir(controller.GameDirectory)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "."+objectName) {
			return strings.Split(f.Name(), ".")[0]
		}
	}
	return controller.getDefaultValueFromObjectName(objectName)
}
func (controller *GameController) GetGameStatFromFilepath(categoryName string) (int, int, int, bool) {
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
	if strings.Contains(strings.ToLower(categoryName), models.DEFAULT_HOME_PREFIX) {
		return models.DEFAULT_WINNING_STAT_INT, models.DEFAULT_LOSTING_STAT, models.DEFAULT_TOTAL_STAT_INT, true
	} else if strings.Contains(strings.ToLower(categoryName), models.DEFAULT_AWAY_PREFIX) {
		return models.DEFAULT_WINNING_STAT_INT / 2, models.DEFAULT_LOSTING_STAT / 2, models.DEFAULT_TOTAL_STAT_INT, true
	} else {
		return models.DEFAULT_LOSTING_STAT, models.DEFAULT_WINNING_STAT_INT, models.DEFAULT_TOTAL_STAT_INT, false
	}
}

// Next season : Write a Controller abst object that game controller, team controller, gamestat and game stats controller can use
// To produce data from the Game controller which calls all other controllers produce data and then this can done in parallel with a single
// wait group at this level. this will become the standardized model for All other apps using this MVC style arch
// Produces a path to be touched holding all necessary data for the UI to update. Avoids file opening and closing operations.
func (controller *GameController) getGameStatPath(gameStat *models.TeamGameStat) string {
	var homeHandle string
	anyAwayValue, _ := gameStat.AwayValue.(string)
	anyHomeValue, _ := gameStat.HomeValue.(string)
	awayValue, err := strconv.Atoi(anyAwayValue)
	radioErrors.ErrorLog(err)
	homeValue, err := strconv.Atoi(anyHomeValue)
	radioErrors.ErrorLog(err)
	maxValue := strconv.Itoa(awayValue + homeValue)
	if homeValue > awayValue {
		homeHandle = strconv.FormatBool(true)
	} else {
		homeHandle = strconv.FormatBool(false)
	}
	filename := anyHomeValue + models.VALUE_DELIMITER +
		anyAwayValue + models.VALUE_DELIMITER +
		maxValue + models.VALUE_DELIMITER +
		homeHandle + "." +
		gameStat.Category
	path := filepath.Join(controller.GameDirectory, filename)
	return path
}

func (controller *GameController) getHomeStatPath(gameStat *models.TeamGameStat) string {
	return filepath.Join(controller.GameDirectory, models.DEFAULT_HOME_PREFIX+models.VALUE_DELIMITER+gameStat.Category)
}

func (controller *GameController) getAwayStatPath(gameStat *models.TeamGameStat) string {
	return filepath.Join(controller.GameDirectory, models.DEFAULT_AWAY_PREFIX+models.VALUE_DELIMITER+gameStat.Category)
}

func (controller *GameController) updateGameData() {
	gdo := quickio.GetGameDataObject(controller.Landinglink)
	gvd := quickio.GetGameVersesData(controller.Landinglink)
	controller.gameDataObject = nil
	controller.gameVersesDataObject = nil
	controller.gameDataObject = &gdo
	controller.gameVersesDataObject = &gvd
	go controller.AwayTeamController.UpdateTeamController(&gdo, &gvd)
	go controller.HomeTeamController.UpdateTeamController(&gdo, &gvd)
}

func (controller *GameController) ProduceGameData() {
	var workGroup sync.WaitGroup
	var gameStatWorkGroup sync.WaitGroup
	workGroupCounter := 0
	gameStatWorkGroupCounter := 0
	controller.updateGameData()
	controllers := []TeamController{*controller.HomeTeamController, *controller.AwayTeamController}
	for i := range controllers {
		workGroup.Add(1)
		workGroupCounter += 1
		go func(teamController TeamController) {
			defer workGroup.Done()
			quickio.TouchFile(teamController.GetScorePath())
			quickio.TouchFile(teamController.GetSOGPath())
			quickio.WriteFile(teamController.GetTeamOnIcePath(), string(teamController.getTeamOnIceJson()))
		}(controllers[i])
	}
	teamGameStatObjects := controller.GetTeamGameStatsObjects()
	log.Println("TeameVerses Gamedata Object - Team Game Stats Objects", teamGameStatObjects)
	for i := range teamGameStatObjects {
		gameStatWorkGroup.Add(1)
		gameStatWorkGroupCounter += 1
		go func(gameStat models.TeamGameStat) {
			defer gameStatWorkGroup.Done()
			quickio.TouchFile(controller.getGameStatPath(&gameStat))
			quickio.TouchFile(controller.getHomeStatPath(&gameStat))
			quickio.TouchFile(controller.getAwayStatPath(&gameStat))
		}(teamGameStatObjects[i])
	}
	if workGroupCounter > 0 {
		workGroup.Wait()
	}
	if gameStatWorkGroupCounter > 0 {
		gameStatWorkGroup.Wait()
	}
	gameStatePath := controller.GetGamestatePath()
	quickio.WriteFile(gameStatePath, controller.GetGamestateString())
	controller.dataConsumed = false
}

func (controller *GameController) ConsumeGameData() {
	var workGroup sync.WaitGroup
	controllers := []TeamController{*controller.HomeTeamController, *controller.AwayTeamController}
	for i := range controllers {
		workGroup.Add(1)
		go func(teamController TeamController) {
			defer workGroup.Done()
			teamController.EmptyDirectory()
		}(controllers[i])
	}
	workGroup.Wait()
	controller.EmptyDirectory()
	controller.dataConsumed = true
}

func CreateNewDefaultGameController() *GameController {
	var controller GameController
	controller.gameDataObject = models.CreateDefaultGameData()
	controller.gameVersesDataObject = models.CreateDefaultVersesData()
	controller.Landinglink = ""
	controller.GameDirectory = filepath.Join(quickio.GetQuickTmpFolder(), "GAME")
	controller.HomeTeamController = CreateNewDefaultTeamController()
	controller.AwayTeamController = CreateNewDefaultTeamController()
	return &controller
}

func CreateNewGameController(landingLink string) *GameController {
	var controller GameController
	gdo := quickio.GetGameDataObject(landingLink)
	gvd := quickio.GetGameVersesData(landingLink)
	sweaters := quickio.GetSweaters()
	controller.Landinglink = landingLink
	controller.GameDirectory = filepath.Join(quickio.GetQuickTmpFolder(), strconv.Itoa(gdo.Id))
	controller.HomeTeamController = CreateNewTeamController(sweaters, landingLink, &gdo, &gvd, true, controller.GameDirectory)
	controller.AwayTeamController = CreateNewTeamController(sweaters, landingLink, &gdo, &gvd, false, controller.GameDirectory)
	controller.gameDataObject = &gdo
	controller.gameVersesDataObject = &gvd
	controller.dataConsumed = false
	return &controller
}
