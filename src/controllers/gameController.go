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

func (controller *GameController) IsInIntermission() bool {
	return controller.gameDataObject.GameState == "LIVE" && controller.gameDataObject.Clock.InIntermission
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
				"\n" + "INT " + strconv.Itoa(controller.gameDataObject.PeriodDescriptor.Number) + " " + controller.gameDataObject.Clock.TimeRemaining
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

func (controller *GameController) GetFloatFraction(fraction string) string {
	numerator, _ := strconv.ParseFloat(strings.Split(fraction, "/")[0], 64)
	denominator, _ := strconv.ParseFloat(strings.Split(fraction, "/")[1], 64)
	floatFraction := strconv.FormatFloat(numerator/denominator, 'f', 2, 64)
	return floatFraction
}

func (controller *GameController) GetTextFromObjectNameFilepath(objectName string, defaultString string) string {
	files, _ := os.ReadDir(controller.GameDirectory)
	for _, f := range files {
		if strings.HasPrefix(f.Name(), objectName) {
			text := strings.ReplaceAll(strings.ReplaceAll(strings.Split(f.Name(), ".")[1], "-", "/"), "#", ".")
			if strings.Contains(text, ".") && !strings.Contains(text, "%") {
				stat, _ := strconv.ParseFloat(text, 64)
				stat = stat * 100.0
				text = strconv.FormatFloat(stat, 'f', 2, 64) + "%"
			}
			return text
		}
	}
	return defaultString
}

func (controller *GameController) GetGameStatFloatsFromFilepath(categoryName string) (float64, float64, float64, bool) {
	//Should be a file that exists in game directory that has the infor in the file name ending with gamecategorySlider
	//Abstract this further to save a file I/O operation.
	files, _ := os.ReadDir(controller.GameDirectory)
	for _, f := range files {
		//Rewrite this.
		if strings.HasPrefix(f.Name(), "."+categoryName) {
			values := strings.Split(strings.Split(f.Name(), ".")[0], models.VALUE_DELIMITER)
			//Shouldnt Crash? Maybe for testing but this hould have some sort default if all the values are not there.
			//If we crash we need to have something to dump all of memories and exit gracefully avoiding mem leaks.
			homeValue, err := strconv.ParseFloat(values[1], 64)
			radioErrors.ErrorLog(err)
			awayValue, err := strconv.ParseFloat(values[1], 64)
			radioErrors.ErrorLog(err)
			if homeValue-awayValue == 0 {
				homeValue = homeValue / 2
				awayValue = awayValue / 2
			}
			if len(values) > 2 {
				maxValue, err := strconv.ParseFloat(values[2], 64)
				radioErrors.ErrorLog(err)
				homeHandle, err := strconv.ParseBool(values[3])
				radioErrors.ErrorLog(err)
				log.Println("GameController::GetGameStatFloatsFromFilepath::homeValue, awayValue, maxValue, homeHandle ", homeValue, awayValue, maxValue, homeHandle)
				return homeValue, awayValue, maxValue, homeHandle
			} else {
				maxValue := 1.0
				homeHandle := true
				log.Println("GameController::GetGameStatFloatsFromFilepath::homeValue, awayValue, maxValue, homeHandle ", homeValue, awayValue, maxValue, homeHandle)
				return homeValue, awayValue, maxValue, homeHandle
			}

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

// We need to clean up all code around sliders and file i/o for them.
func (controller *GameController) GetGameStatFromFilepath(categoryName string) (int, int, int, bool) {
	//Should be a file that exists in game directory that has the infor in the file name ending with gamecategorySlider
	//Abstract this further to save a file I/O operation.
	files, _ := os.ReadDir(controller.GameDirectory)
	for _, f := range files {
		if strings.HasSuffix(f.Name(), "."+categoryName) {
			values := strings.Split(strings.Split(f.Name(), ".")[0], models.VALUE_DELIMITER)
			//Shouldnt Crash? Maybe for testing but this hould have some sort default if all the values are not there.
			//If we crash we need to have something to dump all of memories and exit gracefully avoiding mem leaks.
			log.Println("GameController::GetGameStatFromFilepath::CategoryName And Values, len(values) ", categoryName, values, len(values))
			if len(values) > 3 {
				homeValue, err := strconv.Atoi(values[0])
				radioErrors.ErrorLog(err)
				awayValue, err := strconv.Atoi(values[1])
				radioErrors.ErrorLog(err)
				maxValue, err := strconv.Atoi(values[2])
				radioErrors.ErrorLog(err)
				homeHandle, err := strconv.ParseBool(values[3])
				radioErrors.ErrorLog(err)
				log.Println("GameController::GetGameStatFromFilepath::Returning from len > 3 || categoryName, homeValue, awayValue, maxValue, homeHandle ", categoryName, homeValue, awayValue, maxValue, homeHandle)
				return homeValue, awayValue, maxValue, homeHandle
			} else {
				homeValue, err := strconv.Atoi(values[0])
				radioErrors.ErrorLog(err)
				awayValue, err := strconv.Atoi(values[0])
				radioErrors.ErrorLog(err)
				maxValue := 1
				homeHandle := true
				log.Println("GameController::GetGameStatFromFilepath::IN BUGGED LOGIC BRANCH- From fraction. or naked decimal(Replace with # like for labels)::categoryName, homeValue, awayValue, maxValue, homeHandle ", categoryName, homeValue, awayValue, maxValue, homeHandle)
				return homeValue, awayValue, maxValue, homeHandle
			}

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

func (controller *GameController) ConvertAnyStatToGameStatString(gameStat *models.TeamGameStat, home bool) string {
	if home {
		switch value := gameStat.HomeValue.(type) {
		case string:
			return value
		case int:
			return strconv.Itoa(value)
		case float64:
			return strconv.FormatFloat(value, 'G', 10, 64)
		case bool:
			return strconv.FormatBool(value)
		case nil:
			return "nil"
		default:
			return fmt.Sprintf("%v", value)
		}
	} else {
		switch value := gameStat.AwayValue.(type) {
		case string:
			return value
		case int:
			return strconv.Itoa(value)
		case float64:
			return strconv.FormatFloat(value, 'G', 10, 64)
		case bool:
			return strconv.FormatBool(value)
		case nil:
			return "nil"
		default:
			return fmt.Sprintf("%v", value)
		}
	}
}

func (controller *GameController) convertAnyStatToGameStatString(gameStat *models.TeamGameStat, home bool) string {
	if home {
		switch value := gameStat.HomeValue.(type) {
		case string:
			return strings.ReplaceAll(value, "/", "-")
		case int:
			return strconv.Itoa(value)
		case float64:
			return strings.ReplaceAll(strconv.FormatFloat(value, 'G', 10, 64), ".", "#")
		case bool:
			return strconv.FormatBool(value)
		case nil:
			return "nil"
		default:
			return fmt.Sprintf("%v", value)
		}
	} else {
		switch value := gameStat.AwayValue.(type) {
		case string:
			return strings.ReplaceAll(value, "/", "-")
		case int:
			return strconv.Itoa(value)
		case float64:
			return strings.ReplaceAll(strconv.FormatFloat(value, 'G', 10, 64), ".", "#")
		case bool:
			return strconv.FormatBool(value)
		case nil:
			return "nil"
		default:
			return fmt.Sprintf("%v", value)
		}
	}
}

func (controller *GameController) convertAnyStatToGameStatFilename(gameStat *models.TeamGameStat) string {
	var homeHandle string
	switch homeValue := gameStat.HomeValue.(type) {
	case string:
		//Here were gonna get things like powerplay fractions
		//Were gonna wana parse that into floats and then replace the / with a - and then set it
		//Back on the UI side.
		/*
			filename := anyHomeValue + models.VALUE_DELIMITER +
				anyAwayValue + models.VALUE_DELIMITER +
				maxValue + models.VALUE_DELIMITER +
				homeHandle + "." +
				gameStat.Category
			filename pattern must be held. Maybe even having a type in the name so
			  string_1-3_2-4_7_true.ppFraction
			Save the floats for conversion and slider quantization for later
		*/
		log.Println("gameController::convertAnyStatToGameStatFilename::CASE STRING ::", homeValue)
		return homeValue
	case int:
		awayValue := gameStat.AwayValue.(int)
		if homeValue > awayValue {
			homeHandle = strconv.FormatBool(true)
		} else {
			homeHandle = strconv.FormatBool(false)
		}
		maxValue := strconv.Itoa(homeValue + homeValue)
		filename := strconv.Itoa(homeValue) + models.VALUE_DELIMITER +
			strconv.Itoa(awayValue) + models.VALUE_DELIMITER +
			maxValue + models.VALUE_DELIMITER +
			homeHandle + "." + gameStat.Category
		log.Println("gameController::convertAnyStatToGameStatFilename::CASE INT ::", filename)
		return filename
	case float32, float64:
		homeValueFloat := gameStat.HomeValue.(float64)
		awayValue := gameStat.AwayValue.(float64)
		if homeValueFloat > awayValue {
			homeHandle = strconv.FormatBool(true)
		} else {
			homeHandle = strconv.FormatBool(false)
		}
		homeValueString := strconv.FormatFloat(homeValueFloat, 'G', 10, 64)
		awayValueString := strconv.FormatFloat(awayValue, 'G', 10, 64)
		maxValueString := strconv.FormatFloat(homeValueFloat+awayValue, 'G', 10, 64)
		filename := homeValueString + models.VALUE_DELIMITER +
			awayValueString + models.VALUE_DELIMITER +
			maxValueString + models.VALUE_DELIMITER +
			homeHandle + "." + gameStat.Category
		log.Println("gameController::convertAnyStatToGameStatFilename::CASE FLOAT64 or FLOAT32 ::", filename)
		return filename
	case bool:
		log.Println("gameController::convertAnyStatToGameStatFilename::CASE BOOL ::", gameStat.Category)
		return strconv.FormatBool(homeValue)
	case nil:
		log.Println("gameController::convertAnyStatToGameStatFilename::CASE nil ::", gameStat.Category)
		return "nil"
	default:
		log.Println("gameController::convertAnyStatToGameStatFilename::CASE DEFAULT ::", gameStat.Category)
		return fmt.Sprintf("%v", homeValue)
	}

}

// Next season : Write a Controller abst object that game controller, team controller, gamestat and game stats controller can use
// To produce data from the Game controller which calls all other controllers produce data and then this can done in parallel with a single
// wait group at this level. this will become the standardized model for All other apps using this MVC style arch
// Produces a path to be touched holding all necessary data for the UI to update. Avoids file opening and closing operations.
func (controller *GameController) getGameStatPath(gameStat *models.TeamGameStat) string {
	filename := controller.convertAnyStatToGameStatFilename(gameStat)
	path := filepath.Join(controller.GameDirectory, filename)
	return path
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

func (controller *GameController) getHomeStatPath(category string, value string) string {
	return filepath.Join(controller.GameDirectory, models.DEFAULT_HOME_PREFIX+models.VALUE_DELIMITER+category+"."+value)
}

func (controller *GameController) getAwayStatPath(category string, value string) string {
	return filepath.Join(controller.GameDirectory, models.DEFAULT_AWAY_PREFIX+models.VALUE_DELIMITER+category+"."+value)
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

func (controller *GameController) ProduceTeamGameStats() {
	var gameStatWorkGroup sync.WaitGroup
	teamGameStatObjects := controller.GetTeamGameStatsObjects()
	for i := range teamGameStatObjects {
		gameStatWorkGroup.Add(1)
		go func(gameStat models.TeamGameStat) {
			defer gameStatWorkGroup.Done()
			quickio.TouchFile(controller.getGameStatPath(&gameStat))
			quickio.TouchFile(controller.getHomeStatPath(gameStat.Category, controller.convertAnyStatToGameStatString(&gameStat, true)))
			quickio.TouchFile(controller.getAwayStatPath(gameStat.Category, controller.convertAnyStatToGameStatString(&gameStat, false)))
		}(teamGameStatObjects[i])
	}
	gameStatWorkGroup.Wait()
}

// We want a Go Update Players func that updates players on ice in parallel.
// Each Label per Player also needs to be in parallel
func (controller *GameController) ProduceGameData() {
	var workGroup sync.WaitGroup
	controller.updateGameData()
	controllers := []TeamController{*controller.HomeTeamController, *controller.AwayTeamController}
	for i := range controllers {
		workGroup.Add(1)
		go func(teamController TeamController) {
			defer workGroup.Done()
			quickio.TouchFile(teamController.GetScorePath())
			quickio.TouchFile(teamController.GetSOGPath())
			for index, player := range teamController.GetAllPlayersOnIce() {
				quickio.TouchFile(teamController.GetSweaterNumberPath(index, player))
				quickio.TouchFile(teamController.GetPlayerNamePath(index, player))
				quickio.TouchFile(teamController.GetPositioncodePath(index, player))
			}
			for index, player := range teamController.TeamOnIce.PenaltyBox {
				pentalyBoxIndex := index + 10
				quickio.TouchFile(teamController.GetSweaterNumberPath(pentalyBoxIndex, player))
				quickio.TouchFile(teamController.GetPlayerNamePath(pentalyBoxIndex, player))
				quickio.TouchFile(teamController.GetPositioncodePath(pentalyBoxIndex, player))
			}
		}(controllers[i])
	}
	workGroup.Wait()
	controller.ProduceTeamGameStats()
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
