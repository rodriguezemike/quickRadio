package controllers

import (
	"os"
	"path/filepath"
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

}

func TestEmptyActiveGameDirectory(t *testing.T) {
	controller := NewGameController()
	controller.DumpGameData()
	controller.EmptyActiveGameDirectory()
	files, _ := os.ReadDir(controller.activeGameDirectory)
	for _, f := range files {
		info, _ := f.Info()
		t.Fatalf(`Found a File. There should no files in the active directory. %s`, info.Name())
	}
}

func TestSwitchActiveDataObjects(t *testing.T) {
	controller := NewGameController()
	wantedActiveGameDataObject := &controller.gameDataObjects[2]
	wantedActiveVersesGameDataObject := &controller.gameVersesObjects[2]
	controller.SwitchActiveObjects(2)
	if controller.ActiveGameDataObject != wantedActiveGameDataObject {
		t.Fatalf(`&controller.ActiveGameDataObject != &wantedActiveGameDataObject | Address wanted %v | Address got %v`, &wantedActiveGameDataObject, &controller.ActiveGameDataObject)
	}
	if controller.ActiveGameVersesDataObject != wantedActiveVersesGameDataObject {
		t.Fatalf(`&controller.ActiveGameVersesDataObject != wantedActiveVersesGameDataObject | Address wanted %v | Address got %v`, &wantedActiveVersesGameDataObject, &controller.ActiveGameVersesDataObject)
	}
}

func TestDumpGameData(t *testing.T) {
	controller := NewGameController()
	controller.DumpGameData()
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
	if reflect.ValueOf(controller.ActiveGameDataObject).Kind().String() != "struct" {
		t.Fatalf(`reflect.ValueOf(controller.ActiveGameDataObject).Kind().String() != "struct" | reflect.ValueOf(controller.ActiveGameDataObject).Kind()) = %s`, reflect.ValueOf(controller.ActiveGameDataObject).Kind())
	}
	if reflect.ValueOf(controller.ActiveGameVersesDataObject).Kind().String() != "struct" {
		t.Fatalf(`reflect.ValueOf(controller.ActiveGameVersesDataObject).Kind().String() != "struct" | reflect.ValueOf(controller.ActiveGameVersesDataObject).Kind()) = %s`, reflect.ValueOf(controller.ActiveGameVersesDataObject).Kind())
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
