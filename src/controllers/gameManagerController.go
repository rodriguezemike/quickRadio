//go:build ignore
// +build ignore

package controllers

import (
	"context"
	"errors"
	"quickRadio/models"
	"quickRadio/quickio"
	"sync"
)

type GameManagerController struct {
	Landinglinks         []string
	ActiveGameIndex      int
	ActiveGameController *GameController
	ActiveRadioLink      string
	Sweaters             map[string]*models.Sweater
	gameControllers      []GameController
	ctx                  context.Context
	goroutineMap         *sync.Map
}

func (controller *GameManagerController) GetGameController(index int) GameController {
	return controller.gameControllers[index]
}
func (controller *GameManagerController) GetGameControllers() []GameController {
	return controller.gameControllers
}

func (controller *GameManagerController) UpdateGameControllers() {
	gameDataObjects = quickio.GoGetGameDataObjectsFromLandingLinks(controller.Landinglinks)
	gameVersesObjects = quickio.GoGetGameVersesDataFromLandingLinks(controller.Landinglinks)
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
