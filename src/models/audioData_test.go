package models

import (
	"testing"
	"time"

	"github.com/gopxl/beep/speaker"
)

func TestAudioDataInit(t *testing.T) {
	var testStream = &AudioStreamQueue{}
	speaker.Play(testStream)
	time.Sleep(5 * time.Second)
}
