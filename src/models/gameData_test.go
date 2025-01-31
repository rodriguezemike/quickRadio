package models_test

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"testing"
)

func TestGameData(t *testing.T) {
	var gameData = &models.GameData{}
	dir := quickio.GetProjectDir()
	testFilePath := filepath.Join(dir, "assets", "tests", "gamelanding.json")
	jsonFileObject, err := os.Open(testFilePath)
	if err != nil {
		t.Fatalf(`Error opening test file %s Error -> %s`, testFilePath, err)
	}
	defer jsonFileObject.Close()
	byteValue, err := io.ReadAll(jsonFileObject)
	if err != nil {
		t.Fatalf(`Error Reading Json file %s Error -> %s`, testFilePath, err)
	}
	err = json.Unmarshal(byteValue, gameData)
	if err != nil {
		t.Fatalf(`Could not Unmarshal Test Json File. %s Error -> %s`, testFilePath, err)
	}
}
