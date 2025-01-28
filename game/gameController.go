//go:build ignore
// +build ignore

package game

import (
	"errors"
	"os"
	"quickRadio/models"
	"quickRadio/quickio"
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

//ToDo:Write files stuff and gorountines for handling the files, can just dump data into files for some and touch some other filse.
//Alternatively we can just keep it all in memory but then would require more ui logic.

func (controller *GameController) GetUIDataFromFilename(teamAbbrev string, dataLabel string, defaultReturnValue string) string {
	for f, _ := range os.ReadDir(controller.activeGameDirectory) {
		if strings.Contains(f.Info().Name(), teamAbbrev) && strings.Contains(f.Info().Name(), dataLabel) {
			return Strings.Split(f.Info().Name(), ".")[1]
		}
	}
	return defaultReturnValue
}

func (controller *GameController) DumpGameData() {
	homeScore := filepath.Join(controller.activeGameDirectory, activeGameDataObject.HomeTeam.TeamAbbrev + "_SCORE"+"."activeGameDataObject.HomeTeam.Score)
	awayScore := filepath.Join(controller.activeGameDirectory, activeGameDataObject.AwayTeam.TeamAbbrev + "_SCORE"+"."activeGameDataObject.AwayTeam.Score)
	//Funcify
	f, _ := os.Create(homeScore)
	f.Close()
	f, _ := os.Create(awayScore)
	//Ice Ticker (Game State - Period - TIme Left)
	//Players On ice and in Penalty Box (On home/away side)
	//Game Stats including sog, hits, faceoffs, etc (Center)
}

func (controller *GameController) UpdateDataObjects() {
	controller.gameDataObjects = quickio.GoGetGameDataObjectFromLandingLinks(landingLinks)
	controller.gameVersesObjects = quickio.GoGetGameVersesDataFromLandingLinks(landingLinks)
	controller.UpdateActiveObjects(controller.ActiveGameIndex)
}

func (controller *GameController) UpdateActiveDataObjects(){
	controller.ActiveGameDataObject = quickio.GetGameDataObject(controller.ActiveLandingLink)
	controller.ActiveGameVersesDataObject = quickio.GoGetGameVersesDataFromLandingLinks({controller.ActiveLandingLink})
}

func (controller *GameController) SwitchActiveObjects(gameIndex int) {
	controller.ActiveGameDataObject = controller.gameDataObjects[gameIndex]
	controller.ActiveGameVersesDataObject = controller.gameDataObjects[gameIndex]
	controller.ActiveGameIndex = gameIndex
}

func (controller *GameController) GetActiveRadioLink(teamAbbrev string) {
	if controller.ActiveGameDataObject.AwayTeam.Abbrev == teamAbbrev {
		return controller.ActiveGameDataObject.AwayTeam.RadioLink, nil
	} else if controller.ActiveGameDataObject.HomeTeam.Abbrev == teamAbbrev {
		return controller.ActiveGameDataObject.HomeTeam.RadioLink, nil
	} else {
		return "", errors.New("Couldnt Find a RadioLink from ActiveGameDataObject.")
	}
}

func NewGameController() GameController {
	var controller GameController
	controller.Landinglinks = quickio.GetGameLandingLinks()
	controller.sweaters = quickio.GetSweaters()
	controller.gameDataObjects = quickio.GoGetGameDataObjectFromLandingLinks(controller.Landinglinks)
	controller.gameVersesObjects = quickio.GoGetGameVersesDataFromLandingLinks(controller.landinglinks)
	controller.activeGameDirectory = quickio.GetActiveGameDirectory()
	controller.ActiveGameIndex = 0
	return controller
}
