package game

import "quickRadio/models"

type GameController struct {
	Landinglinks         []string
	sweaters             map[string]models.Sweater
	gameDataObjects      []models.GameData
	gameVersesObjects    []models.GameVersesData
	activeLandingLink    string
	activeGameDataObject models.GameData
	activeGameVersesData models.GameVersesData
	//Here we want what we will need to update the ui, filepaths, syncMaps, etc.
}
