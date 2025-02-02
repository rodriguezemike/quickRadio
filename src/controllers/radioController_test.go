package controllers

import (
	"path/filepath"
	"quickRadio/quickio"
	"reflect"
	"testing"
)

func GetRadioController() *RadioController {
	radioLink := "https://d2igy0yla8zi0u.cloudfront.net/lak/20242025/lak-radio.m3u8"
	teamAbbrev := "LAK"
	sampleRate := "192K"
	controller := NewRadioController(radioLink, teamAbbrev, sampleRate)
	return controller
}

func TestRadioDirectory(t *testing.T) {
	controller := GetRadioController()
	wantedRadioDirectory := filepath.Join(quickio.GetQuickTmpFolder(), "20242025", "lak-radio_192K")
	if controller.RadioDirectory != wantedRadioDirectory {
		t.Fatalf("controller.RadioDirectory != wantedRadioDirectory || controller.RadioDirectory = %s wantedRadioDirectory = %s", controller.RadioDirectory, wantedRadioDirectory)
	}
	quickio.EmptyTmpFolder()
}

func TestNewRadioControllerTypes(t *testing.T) {
	controller := GetRadioController()
	if reflect.ValueOf(controller.TeamAbbrev).Kind().String() != "string" {
		t.Fatalf(`reflect.ValueOf(controller.TeamAbbrev).Kind().String() != "string" | reflect.ValueOf(controller.TeamAbbrev).Kind().String()  = %s`, reflect.ValueOf(controller.TeamAbbrev).Kind())
	}
	if reflect.ValueOf(controller.WavPaths).Kind().String() != "slice" {
		t.Fatalf(`reflect.ValueOf(controller.WavPaths).Kind().String() != "slice" | reflect.ValueOf(controller.WavPaths).Kind().String() = %s`, reflect.ValueOf(controller.WavPaths).Kind())
	}
	if reflect.ValueOf(controller.RadioDirectory).Kind().String() != "string" {
		t.Fatalf(`reflect.ValueOf(controller.RadioDirectory).Kind().String() != "string" | reflect.ValueOf(controller.RadioDirectory).Kind()) = %s`, reflect.ValueOf(controller.RadioDirectory).Kind())
	}
	if reflect.ValueOf(controller.RadioFormatLink).Kind().String() != "string" {
		t.Fatalf(`reflect.ValueOf(controller.RadioFormatLink).Kind().String() != "string" | reflect.ValueOf(controller.RadioFormatLink).Kind()) = %s`, reflect.ValueOf(controller.RadioFormatLink).Kind())
	}
	if reflect.ValueOf(controller.EmergencySleepInterval).Kind().String() != "int" {
		t.Fatalf(`reflect.ValueOf(controller.EmergencySleepInterval).Kind().String() != "int" | reflect.ValueOf(controller.EmergencySleepInterval).Kind()) = %s`, reflect.ValueOf(controller.EmergencySleepInterval).Kind())
	}
	if reflect.ValueOf(controller.NormalSleepInterval).Kind().String() != "int" {
		t.Fatalf(`reflect.ValueOf(controller.NormalSleepInterval).Kind().String() != "int" | reflect.ValueOf(controller.NormalSleepInterval).Kind()) = %s`, reflect.ValueOf(controller.NormalSleepInterval).Kind())
	}
	if reflect.ValueOf(controller.streamQueue).Kind().String() != "ptr" {
		t.Fatalf(`reflect.ValueOf(controller.streamQueue).Kind().String() != "ptr" | reflect.ValueOf(controller.streamQueue).Kind()) = %s`, reflect.ValueOf(controller.streamQueue).Kind())
	}
	if reflect.ValueOf(controller.streamers).Kind().String() != "slice" {
		t.Fatalf(`reflect.ValueOf(controller.streamers).Kind().String() != "slice" | reflect.ValueOf(controller.streamers).Kind()) = %s`, reflect.ValueOf(controller.streamers).Kind())
	}
	if reflect.ValueOf(controller.speakerInitialized).Kind().String() != "bool" {
		t.Fatalf(`reflect.ValueOf(controller.speakerInitialized).Kind().String() != "bool" | reflect.ValueOf(controller.speakerInitialized).Kind()) = %s`, reflect.ValueOf(controller.speakerInitialized).Kind())
	}
	if reflect.ValueOf(controller.ctx).Kind().String() != "struct" {
		t.Fatalf(`reflect.ValueOf(controller.ctx).Kind().String() != "struct" | reflect.ValueOf(controller.ctx).Kind()) = %s`, reflect.ValueOf(controller.ctx).Kind())
	}
	if reflect.ValueOf(controller.goroutineMap).Kind().String() != "ptr" {
		t.Fatalf(`reflect.ValueOf(controller.goroutineMap).Kind().String() != "ptr" | reflect.ValueOf(controller.goroutineMap).Kind()) = %s`, reflect.ValueOf(controller.goroutineMap).Kind())
	}
	quickio.EmptyTmpFolder()

}
