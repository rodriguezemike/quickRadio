package controllers

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"strconv"
	"strings"
)

type TeamController struct {
	Landinglink    string
	Sweater        *models.Sweater
	Team           *models.TeamData
	TeamOnIce      *models.TeamOnIce
	gameDataObject *models.GameData
	gameVersesData *models.GameVersesData
	home           bool
	teamDirectory  string
}

func (controller *TeamController) UpdateTeamController(gdo *models.GameData, gvd *models.GameVersesData) {
	controller.gameDataObject = nil
	controller.gameVersesData = nil
	controller.Team = nil
	controller.gameDataObject = gdo
	controller.gameVersesData = gvd
	if controller.home {
		controller.Team = &gdo.HomeTeam
		controller.TeamOnIce = &gdo.Summary.IceSurface.HomeTeam
	} else {
		controller.Team = &gdo.AwayTeam
		controller.TeamOnIce = &gdo.Summary.IceSurface.AwayTeam
	}
}

func (controller *TeamController) EmptyDirectory() {
	quickio.EmptyDirectory(controller.teamDirectory)
}

func (controller *TeamController) GetScorePath() string {
	return filepath.Join(controller.teamDirectory, controller.Team.Abbrev+"_SCORE."+strconv.Itoa(controller.Team.Score))
}

func (controller *TeamController) GetSOGPath() string {
	log.Println("SOG PATH ->", controller.teamDirectory, controller.Team.Abbrev+"_SOG."+strconv.Itoa(controller.Team.Sog))
	return filepath.Join(controller.teamDirectory, controller.Team.Abbrev+"_SOG."+strconv.Itoa(controller.Team.Sog))
}

func (controller *TeamController) GetStatsPath() string {
	return filepath.Join(controller.teamDirectory, "STATS.json")
}

func (controller *TeamController) GetTeamOnIcePath() string {
	return filepath.Join(controller.teamDirectory, controller.Team.Abbrev+"_TEAMONICE.json")
}

func (controller *TeamController) GetUIDataFromFilename(dataLabel string, defaultReturnValue string) string {
	files, _ := os.ReadDir(controller.teamDirectory)
	for _, f := range files {
		info, _ := f.Info()
		if strings.Contains(info.Name(), controller.Team.Abbrev) && strings.Contains(info.Name(), dataLabel) {
			return strings.Split(info.Name(), ".")[1]
		}
	}
	return defaultReturnValue
}

func (controller *TeamController) getTeamOnIceJson() []byte {
	onIceJson, _ := json.MarshalIndent(controller.Team, "", " ")
	return onIceJson
}

func CreateNewDefaultTeamController() *TeamController {
	gameDirectory := filepath.Join(quickio.GetQuickTmpFolder(), "GAME")
	controller := TeamController{}
	gdo := models.CreateDefaultGameData()
	gvd := models.CreateDefaultVersesData()
	nhlf := quickio.GetSweaters()["NHLF"]
	controller.Landinglink = ""
	controller.gameDataObject = gdo
	controller.gameVersesData = gvd
	controller.Sweater = &nhlf
	controller.Team = models.CreateDefaultTeam()
	controller.TeamOnIce = models.CreateDefaultTeamOnIce()
	controller.teamDirectory = filepath.Join(gameDirectory, controller.Team.Abbrev)
	return &controller
}

func CreateNewTeamController(sweaters map[string]models.Sweater, landingLink string, gdo *models.GameData, gvd *models.GameVersesData, home bool, gameDirectory string) *TeamController {
	controller := TeamController{}
	controller.Landinglink = landingLink
	controller.gameDataObject = gdo
	controller.gameVersesData = gvd
	controller.home = home
	if home {
		sweater := sweaters[gdo.HomeTeam.Abbrev]
		controller.Sweater = &sweater
		controller.TeamOnIce = &gdo.Summary.IceSurface.HomeTeam
		controller.teamDirectory = filepath.Join(gameDirectory, gdo.HomeTeam.Abbrev)
		controller.Team = &gdo.HomeTeam
	} else {
		sweater := sweaters[gdo.AwayTeam.Abbrev]
		controller.Sweater = &sweater
		controller.TeamOnIce = &gdo.Summary.IceSurface.AwayTeam
		controller.teamDirectory = filepath.Join(gameDirectory, gdo.AwayTeam.Abbrev)
		controller.Team = &gdo.AwayTeam
	}
	return &controller
}

func CreateNewTeamControllersFromLiveGameData(landingLink string, gameDataObject *models.GameData, gameVersesDataObject models.GameVersesData) (*TeamController, *TeamController) {
	gameDirectory := filepath.Join(quickio.GetQuickTmpFolder(), strconv.Itoa(gameDataObject.Id))
	sweaters := quickio.GetSweaters()
	return CreateNewTeamController(sweaters, landingLink, gameDataObject, &gameVersesDataObject, true, gameDirectory),
		CreateNewTeamController(sweaters, landingLink, gameDataObject, &gameVersesDataObject, false, gameDirectory)
}

func CreateNewTeamControllersFromLandingLink(landingLink string) (*TeamController, *TeamController) {
	gameDataObject := quickio.GetGameDataObject(landingLink)
	gameVersesDataObject := quickio.GetGameVersesData(landingLink)
	gameDirectory := filepath.Join(quickio.GetQuickTmpFolder(), strconv.Itoa(gameDataObject.Id))
	sweaters := quickio.GetSweaters()
	return CreateNewTeamController(sweaters, landingLink, &gameDataObject, &gameVersesDataObject, true, gameDirectory),
		CreateNewTeamController(sweaters, landingLink, &gameDataObject, &gameVersesDataObject, false, gameDirectory)
}
