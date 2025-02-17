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

func (controller *GameManagerController) SwitchActiveGame(index int) {
	controller.KillActiveGame()
	controller.SwitchActiveObjects(index)
	controller.RunActiveGame()
}

// We are going to update the games all the time with the freshest data
func (controller *GameManagerController) RunActiveGame() {
	//Here we want to ProduceGameData to be consumed by our UI.
	ctx, cancel := context.WithCancel(controller.ctx)
	defer cancel()
	controller.goroutineMap.Store(ctx, cancel)
	for {
		select {
		case <-ctx.Done():
			controller.goroutineMap.Delete(ctx)
			//Consume All game Data at end of program
			controller.ConsumeAllGameData()
			return
		default:
			//Update Active game first in go call if UI has consumed all data
			//Update rest of games in other call
			if controller.dataConsumed {
				controller.UpdateActiveDataObjects()
				controller.ProduceActiveGameData()
			}
		}
	}
}

func (controller *GameManagerController) KillActiveGame() {
	controller.goroutineMap.Range(func(key, value interface{}) bool {
		callback, _ := value.(context.CancelFunc)
		callback()
		return true
	})
}

func (controller *GameManagerController) UpdateGameControllers() {
	gameDataObjects = quickio.GoGetGameDataObjectsFromLandingLinks(controller.Landinglinks)
	gameVersesObjects = quickio.GoGetGameVersesDataFromLandingLinks(controller.Landinglinks)
	//Update here
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
