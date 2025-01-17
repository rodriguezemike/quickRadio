package game

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"quickRadio/models"
	"runtime"
	"testing"
)

func sliceCompare(sliceA []string, sliceB []string) bool {
	if len(sliceA) == 0 || len(sliceB) == 0 {
		return false
	}
	for i := 0; i < len(sliceA); i++ {
		if sliceA[i] != sliceB[i] {
			return false
		}
	}
	return true
}

func GetTestFileObject(desiredFilename string) *os.File {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	testFilePath := filepath.Join(dir, "assets", "tests", desiredFilename)
	fileObject, _ := os.Open(testFilePath)
	return fileObject
}

func GetGameDataObject() models.GameData {
	var gameData = &models.GameData{}
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	testFilePath := filepath.Join(dir, "assets", "tests", "gamelanding.json")
	jsonFileObject, _ := os.Open(testFilePath)
	byteValue, _ := io.ReadAll(jsonFileObject)
	_ = json.Unmarshal(byteValue, gameData)
	return *gameData
}

func TestGetLinksJson(t *testing.T) {
	want := "https://nhl.com/"
	linksMap := GetLinksJson()
	baseUrl := fmt.Sprintf("%v", linksMap["base"])
	if want != baseUrl {
		t.Fatalf(`Loading Links Json does not return wanted string. Check links.json in assets and check wanted string Want -> %s | Got -> %s`, want, baseUrl)
	}
}

func TestGetRadioLink(t *testing.T) {
	testAbbrev := "LAK"
	want := "https://d2igy0yla8zi0u.cloudfront.net/lak/20242025/lak-radio.m3u8"
	gdo := GetGameDataObject()
	radioLink, err := GetRadioLink(gdo, testAbbrev)
	if want != radioLink || err != nil {
		t.Fatalf(`GetRadioLink(GameDataObject, testAbbrev) = %q, %v, want match for %#q, nil`, radioLink, err, want)
	}
}
