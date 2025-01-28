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

func (controller *GameController) getGamestateString() string {
	if controller.ActiveGameDataObject.GameState == "LIVE" || controller.ActiveGameDataObject.GameState == "CRIT"{
		if !controller.ActiveGameDataObject.InIntermission{
			return controller.ActiveGameDataObject.GameState + " - " \
					+ controller.ActiveGameDataObject.Venue.Default + ", " + controller.ActiveGameDataObject.VenueLocation.Default \
					+ " - " + "P" + strconv.Itoa(gdo.PeriodDescriptor.Number) + " " + gdo.Clock.TimeRemaining)
		} else {
			return controller.ActiveGameDataObject.GameState + " - " \
			+ controller.ActiveGameDataObject.Venue.Default + ", " + controller.ActiveGameDataObject.VenueLocation.Default \
			+ " INT " + strconv.Itoa(gdo.PeriodDescriptor.Number) + " " + gdo.Clock.TimeRemaining)
		}
	} else {
		if controller.ActiveGameDataObject.GameState == "FUT" { 
			return controller.ActiveGameDataObject.GameState + " - "\
			+ controller.ActiveGameDataObject.GameDate + " - " + controller.ActiveGameDataObject.StartTimeUTC\
			+ controller.ActiveGameDataObject.Venue.Default + ", " + controller.ActiveGameDataObject.VenueLocation.Default
		} else {
			return controller.ActiveGameDataObject.GameState
		}
	}
}

func (controller *GameController) getTeamOnIceJson(team models.TeamOnIce)[]byte{
	onIceJson, _ := json.MarshalIndent(team, "", " ")
	return onIceJson
}

func (controller *GameController) getTeamGameStats()[]byte{
	tameGameStats, _ := json.MarshalIndent(controller.ActiveGameVersesDataObject.GameInfo.TeamGameStats)
	return teamGameStats
}

func (controller *GameController) DumpGameData() {
	homeScorePath := filepath.Join(controller.activeGameDirectory, activeGameDataObject.HomeTeam.TeamAbbrev + "_SCORE."+ activeGameDataObject.HomeTeam.Score)
	awayScorePath := filepath.Join(controller.activeGameDirectory, activeGameDataObject.AwayTeam.TeamAbbrev + "_SCORE."+ activeGameDataObject.AwayTeam.Score)
	gameStatePath := filepath.Join(controller.activeGameDirectory, "ActiveGameState.label")
	homePlayersOnIcePath := filepath.Join(controller.activeGameDirectory, activeGameDataObject.HomeTeam.TeamAbbrev + "_playersOnIce."+"json")
	awayPlayersOnIcePath := filepath.Join(controller.activeGameDirectory, activeGameDataObject.AwayTeam.TeamAbbrev + "_playersOnIce."+"json")
	go quickio.touchFile(homeScorePath)
	go quickio.touchFile(awayScorePath)
	go quickio.writeFile(gameStatePath, controller.getGamestateString())
	go quickio.writeFile(homePlayersOnIcePath, controller.getTeamOnIceJson(activeGameDataObject.Summary.IceSurface.HomeTeam))
	go quickio.write File(awayPlayersOnIcePath, controller.getTeamOnIceJson(activeGameDataObject.Summary.IceSurface.AwayTeam))

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
