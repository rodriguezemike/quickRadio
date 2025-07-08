//go:build ignore
// +build ignore

package controllers

import (
	"log"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"reflect"
	"strings"
	"testing"
)

func TestActiveDirectoryValue(t *testing.T) {
	controller := NewGameController()
	want := filepath.Join(os.TempDir(), "QuickRadio", "ActiveGame")
	if want != controller.activeGameDirectory {
		t.Fatalf(`controller.activeGameDirectory != %s || controller.activeGameDirectory = %s`, want, controller.activeGameDirectory)
	}
	quickio.EmptyTmpFolder()
}

func TestEmptyActiveGameDirectory(t *testing.T) {
	controller := NewGameController()
	controller.ProduceActiveGameData()
	controller.EmptyActiveGameDirectory()
	files, _ := os.ReadDir(controller.activeGameDirectory)
	for _, f := range files {
		info, _ := f.Info()
		t.Fatalf(`Found a File. There should no files in the active directory. %s`, info.Name())
	}
}

func TestSwitchActiveDataObjects(t *testing.T) {
	controller := NewGameController()
	gamesLen := len(controller.Landinglinks)
	if gamesLen > 0 {
		wantedActiveGameDataObject := &controller.gameDataObjects[gamesLen-1]
		wantedActiveVersesGameDataObject := &controller.gameVersesObjects[gamesLen-1]
		controller.SwitchActiveObjects(gamesLen - 1)
		if controller.ActiveGameDataObject != wantedActiveGameDataObject {
			t.Fatalf(`&controller.ActiveGameDataObject != &wantedActiveGameDataObject | Address wanted %v | Address got %v`, &wantedActiveGameDataObject, &controller.ActiveGameDataObject)
		}
		if controller.ActiveGameVersesDataObject != wantedActiveVersesGameDataObject {
			t.Fatalf(`&controller.ActiveGameVersesDataObject != wantedActiveVersesGameDataObject | Address wanted %v | Address got %v`, &wantedActiveVersesGameDataObject, &controller.ActiveGameVersesDataObject)
		}
		quickio.EmptyTmpFolder()
	} else {
		log.Println("gameContoller_test::TestSwitchActiveDataObjects", "SSSKKKKKKKKKIPPPP - No active games today :C if only hockey season was 365 days a year.")
	}

}

func TestDumpGameData(t *testing.T) {
	controller := NewGameController()
	gamesLen := len(controller.Landinglinks)
	if gamesLen > 0 {
		controller.ProduceActiveGameData()
		files, _ := os.ReadDir(controller.activeGameDirectory)
		filesFound := 0
		scoreFilesFound := 0
		gameStatesFilesFound := 0
		playersOnIceFilesFound := 0
		teamGameStatsFilesFound := 0

		wantedNumberOfFiles := 6
		wantedScoreFiles := 2
		wantedGameStateFiles := 1
		wantedPlayersOnIceFiles := 2
		wantedTeamGameStatsFiles := 1
		for _, f := range files {
			filesFound += 1
			info, _ := f.Info()
			if strings.Contains(info.Name(), "SCORE") {
				scoreFilesFound += 1
			}
			if strings.Contains(info.Name(), "ACTIVEGAMESTATE") {
				gameStatesFilesFound += 1
			}
			if strings.Contains(info.Name(), "PLAYERSONICE") {
				playersOnIceFilesFound += 1
			}
			if strings.Contains(info.Name(), "TEAMGAMESTATS") {
				teamGameStatsFilesFound += 1
			}
		}
		if filesFound != wantedNumberOfFiles {
			t.Fatalf(`filesFound != wantedNumberOfFiles || Found %d | Wanted %d`, filesFound, wantedNumberOfFiles)
		}
		if scoreFilesFound != wantedScoreFiles {
			t.Fatalf(`scoreFilesFound != wantedScoreFiles || Found %d | Wanted %d`, scoreFilesFound, wantedScoreFiles)
		}
		if gameStatesFilesFound != wantedGameStateFiles {
			t.Fatalf(`gameStatesFilesFound != wantedGameStateFiles || Found %d | Wanted %d`, gameStatesFilesFound, wantedGameStateFiles)
		}
		if playersOnIceFilesFound != wantedPlayersOnIceFiles {
			t.Fatalf(`playersOnIceFilesFound != wantedPlayersOnIceFiles || Found %d | Wanted %d`, playersOnIceFilesFound, wantedPlayersOnIceFiles)
		}
		if teamGameStatsFilesFound != wantedTeamGameStatsFiles {
			t.Fatalf(`teamGameStatsFilesFound != wantedTeamGameStatsFiles || Found %d | Wanted %d`, teamGameStatsFilesFound, wantedTeamGameStatsFiles)
		}
		quickio.EmptyTmpFolder()
	} else {
		log.Println("gameContoller_test::TestDumpGameData", "SSSKKKKKKKKKIPPPP - No active games today :C if only hockey season was 365 days a year.")
	}

}

func TestNewGameContollerTypes(t *testing.T) {
	controller := NewGameController()
	if reflect.ValueOf(controller.Landinglinks).Kind().String() != "slice" {
		t.Fatalf(`reflect.ValueOf(controller.Landinglinks).Kind().String() != "slice" | reflect.ValueOf(controller.Landinglinks).Kind()) = %s`, reflect.ValueOf(controller.Landinglinks).Kind())
	}
	if reflect.ValueOf(controller.ActiveGameIndex).Kind().String() != "int" {
		t.Fatalf(`reflect.ValueOf(controller.ActiveGameIndex).Kind().String() != "int" | reflect.ValueOf(controller.ActiveGameIndex).Kind()) = %s`, reflect.ValueOf(controller.ActiveGameIndex).Kind())
	}
	if reflect.ValueOf(controller.ActiveLandingLink).Kind().String() != "string" {
		t.Fatalf(`reflect.ValueOf(controller.ActiveLandingLink).Kind().String() != "string" | reflect.ValueOf(controller.ActiveLandingLink).Kind()) = %s`, reflect.ValueOf(controller.ActiveLandingLink).Kind())
	}
	if reflect.ValueOf(controller.ActiveGameDataObject).Kind().String() != "ptr" {
		t.Fatalf(`reflect.ValueOf(controller.ActiveGameDataObject).Kind().String() != "ptr" | reflect.ValueOf(controller.ActiveGameDataObject).Kind()) = %s`, reflect.ValueOf(controller.ActiveGameDataObject).Kind())
	}
	if reflect.ValueOf(controller.ActiveGameVersesDataObject).Kind().String() != "ptr" {
		t.Fatalf(`reflect.ValueOf(controller.ActiveGameVersesDataObject).Kind().String() != "ptr" | reflect.ValueOf(controller.ActiveGameVersesDataObject).Kind()) = %s`, reflect.ValueOf(controller.ActiveGameVersesDataObject).Kind())
	}
	if reflect.ValueOf(controller.Sweaters).Kind().String() != "map" {
		t.Fatalf(`reflect.ValueOf(controller.Sweaters).Kind().String() != "map" | reflect.ValueOf(controller.Sweaters).Kind()) = %s`, reflect.ValueOf(controller.Sweaters).Kind())
	}
	if reflect.ValueOf(controller.GetGameDataObjects()).Kind().String() != "slice" {
		t.Fatalf(`reflect.ValueOf(controller.GetGameDataObjects()).Kind().String() != "slice" | reflect.ValueOf(controller.GetGameDataObjects()).Kind()) = %s`, reflect.ValueOf(controller.GetGameDataObjects()).Kind())
	}
	if reflect.ValueOf(controller.GetGameVersesObjects()).Kind().String() != "slice" {
		t.Fatalf(`reflect.ValueOf(controller.GetGameVersesObjects()).Kind().String() != "slice" | reflect.ValueOf(controller.GetGameVersesObjects()).Kind()) = %s`, reflect.ValueOf(controller.GetGameVersesObjects()).Kind())
	}
	if reflect.ValueOf(controller.GetGameVersesObjects()).Kind().String() != "slice" {
		t.Fatalf(`reflect.ValueOf(controller.GetGameVersesObjects()).Kind().String() != "slice" | reflect.ValueOf(controller.GetGameVersesObjects()).Kind()) = %s`, reflect.ValueOf(controller.GetGameVersesObjects()).Kind())
	}
	if reflect.ValueOf(controller.activeGameDirectory).Kind().String() != "string" {
		t.Fatalf(`reflect.ValueOf(controller.activeGameDirectory).Kind().String() != "string" | reflect.ValueOf(controller.activeGameDirectory).Kind()) = %s`, reflect.ValueOf(controller.activeGameDirectory).Kind())
	}
}

func TestPowerplayDetection(t *testing.T) {
	controller := CreateNewDefaultGameController()

	// Test default game data (no powerplay)
	if controller.IsPowerplayActive() {
		t.Fatalf("Default game data should not have powerplay active")
	}

	// Test with powerplay data
	controller.gameDataObject.Situation.HomeTeam.SituationDescriptions = []string{models.POWERPLAY_INDICATOR}
	controller.gameDataObject.Situation.HomeTeam.Abbrev = "LAK"
	controller.gameDataObject.Situation.AwayTeam.Abbrev = "VGK"

	if !controller.IsPowerplayActive() {
		t.Fatalf("Should detect powerplay when home team has PP")
	}

	powerplayTeam := controller.GetPowerplayTeam()
	if powerplayTeam != "LAK" {
		t.Fatalf("Expected powerplay team 'LAK', got '%s'", powerplayTeam)
	}

	// Test away team powerplay
	controller.gameDataObject.Situation.HomeTeam.SituationDescriptions = []string{}
	controller.gameDataObject.Situation.AwayTeam.SituationDescriptions = []string{models.POWERPLAY_INDICATOR}

	if !controller.IsPowerplayActive() {
		t.Fatalf("Should detect powerplay when away team has PP")
	}

	powerplayTeam = controller.GetPowerplayTeam()
	if powerplayTeam != "VGK" {
		t.Fatalf("Expected powerplay team 'VGK', got '%s'", powerplayTeam)
	}

	// Test no powerplay
	controller.gameDataObject.Situation.AwayTeam.SituationDescriptions = []string{}

	if controller.IsPowerplayActive() {
		t.Fatalf("Should not detect powerplay when no team has PP")
	}

	powerplayTeam = controller.GetPowerplayTeam()
	if powerplayTeam != "" {
		t.Fatalf("Expected empty powerplay team, got '%s'", powerplayTeam)
	}

	quickio.EmptyTmpFolder()
}

func TestPowerplayPathGeneration(t *testing.T) {
	controller := CreateNewDefaultGameController()
	controller.gameDataObject.Situation.HomeTeam.SituationDescriptions = []string{models.POWERPLAY_INDICATOR}
	controller.gameDataObject.Situation.HomeTeam.Abbrev = "LAK"

	powerplayPath := controller.GetPowerplayPath()
	expectedPath := filepath.Join(controller.GameDirectory, models.POWERPLAY_PREFIX+".LAK")

	if powerplayPath != expectedPath {
		t.Fatalf("Expected powerplay path '%s', got '%s'", expectedPath, powerplayPath)
	}

	quickio.EmptyTmpFolder()
}
