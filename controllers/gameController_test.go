package controllers

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"reflect"
	"testing"
)

func GetGameDataObjectForTest() models.GameData {
	var gameData = &models.GameData{}
	fileObject := quickio.GetTestFileObject("quicklanding.json")
	defer fileObject.Close()
	byteValue, _ := io.ReadAll(fileObject)
	_ = json.Unmarshal(byteValue, gameData)
	return *gameData
}

func TestActiveDirectoryValue(t *testing.T) {
	controller := NewGameController()
	want := filepath.Join(os.TempDir(), "QuickRadio", "ActiveGame")
	if want != controller.activeGameDirectory {
		t.Fatalf(`controller.activeGameDirectory != %s || controller.activeGameDirectory = %s`, want, controller.activeGameDirectory)
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
