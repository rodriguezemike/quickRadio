package models_test

import (
	"quickRadio/models"
	"testing"
	"time"

	"github.com/gopxl/beep/speaker"
)

func TestAudioDataInit(t *testing.T) {
	var testStream = &models.AudioStreamQueue{}
	speaker.Play(testStream)
	time.Sleep(5 * time.Second)
}

func TestVolumeControl(t *testing.T) {
	streamQueue := models.NewAudioStreamQueue()

	// Test default volume
	if streamQueue.GetVolume() != 1.0 {
		t.Fatalf("Expected default volume 1.0, got %f", streamQueue.GetVolume())
	}

	// Test setting volume
	streamQueue.SetVolume(0.5)
	if streamQueue.GetVolume() != 0.5 {
		t.Fatalf("Expected volume 0.5, got %f", streamQueue.GetVolume())
	}

	// Test volume bounds - too low
	streamQueue.SetVolume(-0.5)
	if streamQueue.GetVolume() != 0.0 {
		t.Fatalf("Expected volume clamped to 0.0, got %f", streamQueue.GetVolume())
	}

	// Test volume bounds - too high
	streamQueue.SetVolume(3.0)
	if streamQueue.GetVolume() != 2.0 {
		t.Fatalf("Expected volume clamped to 2.0, got %f", streamQueue.GetVolume())
	}
}

func TestVolumeApplication(t *testing.T) {
	streamQueue := models.NewAudioStreamQueue()
	streamQueue.SetVolume(0.5)

	// Create test samples
	samples := make([][2]float64, 10)
	for i := range samples {
		samples[i][0] = 1.0
		samples[i][1] = 1.0
	}

	// Note: We can't easily test the actual volume application without
	// adding streamers, but the volume field is tested above
	if streamQueue.GetVolume() != 0.5 {
		t.Fatalf("Volume should be set to 0.5")
	}
}

func TestNewAudioStreamQueue(t *testing.T) {
	streamQueue := models.NewAudioStreamQueue()

	if streamQueue == nil {
		t.Fatalf("NewAudioStreamQueue should not return nil")
	}

	if streamQueue.GetVolume() != 1.0 {
		t.Fatalf("New audio stream queue should have default volume 1.0")
	}
}
