package controllers

import (
	"encoding/json"
	"os"
	"quickRadio/models"
	"strings"
)

type TeamController struct {
	Landinglink    string
	Sweater        *models.Sweater
	Team           *models.TeamData
	playersOnIce   []models.PlayerOnIce
	gameDataObject *models.GameData
	gameVersesData *models.GameVersesData
	gameDirectory  string
}

func (controller *TeamController) GetTeam() *models.TeamData {
	return controller.team
}

func (controller *TeamController) GetUIDataFromFilename(dataLabel string, defaultReturnValue string) string {
	files, _ := os.ReadDir(controller.gameDirectory)
	for _, f := range files {
		info, _ := f.Info()
		if strings.Contains(info.Name(), controller.team.Abbrev) && strings.Contains(info.Name(), dataLabel) {
			return strings.Split(info.Name(), ".")[1]
		}
	}
	return defaultReturnValue
}

func (controller *TeamController) getTeamOnIceJson() []byte {
	onIceJson, _ := json.MarshalIndent(controller.team, "", " ")
	return onIceJson
}
