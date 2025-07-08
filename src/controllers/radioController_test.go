package controllers

import (
	"path/filepath"
	"quickRadio/quickio"
	"reflect"
	"testing"
)

func GetTestRadioController() *RadioController {
	radioLink := "https://d2igy0yla8zi0u.cloudfront.net/lak/20242025/lak-radio.m3u8"
	teamAbbrev := "LAK"
	controller := NewRadioController(radioLink, teamAbbrev)
	sampleRate := "192K"
	controller.RadioFormatLink, controller.RadioDirectory, controller.WavPaths = quickio.GetRadioFormatLinkAndDirectory(controller.RadioLink, sampleRate)
	return controller
}

func TestRadioDirectory(t *testing.T) {
	controller := GetTestRadioController()
	wantedRadioDirectory := filepath.Join(quickio.GetQuickTmpFolder(), "20242025", "lak-radio_192K")
	if controller.RadioDirectory != wantedRadioDirectory {
		t.Fatalf("controller.RadioDirectory != wantedRadioDirectory || controller.RadioDirectory = %s wantedRadioDirectory = %s", controller.RadioDirectory, wantedRadioDirectory)
	}
	quickio.EmptyTmpFolder()
}

func TestNewRadioControllerTypes(t *testing.T) {
	controller := GetTestRadioController()
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

func TestRadioControllerVolumeControl(t *testing.T) {
	controller := GetTestRadioController()

	// Test default volume
	defaultVolume := controller.GetVolume()
	if defaultVolume != 1.0 {
		t.Fatalf("Expected default volume 1.0, got %f", defaultVolume)
	}

	// Test setting volume
	controller.SetVolume(0.7)
	volume := controller.GetVolume()
	if volume != 0.7 {
		t.Fatalf("Expected volume 0.7, got %f", volume)
	}

	// Test volume bounds
	controller.SetVolume(-0.5)
	volume = controller.GetVolume()
	if volume != 0.0 {
		t.Fatalf("Expected volume clamped to 0.0, got %f", volume)
	}

	controller.SetVolume(3.0)
	volume = controller.GetVolume()
	if volume != 2.0 {
		t.Fatalf("Expected volume clamped to 2.0, got %f", volume)
	}

	quickio.EmptyTmpFolder()
}

func TestRadioControllerVolumeWithNilStreamQueue(t *testing.T) {
	controller := GetTestRadioController()
	controller.streamQueue = nil

	// Should handle nil gracefully
	volume := controller.GetVolume()
	if volume != 1.0 {
		t.Fatalf("Expected default volume 1.0 when streamQueue is nil, got %f", volume)
	}

	// Should not panic when setting volume with nil streamQueue
	controller.SetVolume(0.5)
	// If we get here without panicking, the test passes

	quickio.EmptyTmpFolder()
}
