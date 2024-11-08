package internals

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"quickRadio/models"
	"regexp"
	"runtime"
	"testing"
)

func GetGameDataObject() models.GameDataStruct {
	var gameData = &models.GameDataStruct{}
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
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

func TestGetGetRadioLink(t *testing.T) {
	testAbbrev := "LAK"
	want := regexp.MustCompile("Some Regex here for the radio link")
	gdo := GetGameDataObject()
	radioLink, err := GetRadioLink(gdo, testAbbrev)
	if !want.MatchString(radioLink) || err != nil {
		t.Fatalf(`GetRadioLink(GameDataObject, testAbbrev) = %q, %v, want match for %#q, nil`, radioLink, err, want)
	}
}
