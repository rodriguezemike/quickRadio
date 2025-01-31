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
