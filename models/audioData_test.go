package models

import (
	"errors"
	"os"
	"path/filepath"
	"quickRadio/audio"
	"runtime"
	"testing"
	"time"

	"github.com/gopxl/beep/speaker"
)

func TestStream(t *testing.T) {
	var audioDataQueue AudioStreamQueue
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	testInput := filepath.Join(dir, "assets", "tests", "game_192k_00001.wav")
	if _, err := os.Stat(testInput); errors.Is(err, os.ErrNotExist) {
		testAAC := filepath.Join(dir, "assets", "tests", "game_192k_00001.aac")
		testInput = audio.TranscodeToWave(testAAC)
	}
	streamer := audio.InitializeRadio(testInput)
	defer streamer.Close()
	speaker.Play(&audioDataQueue)
	for i := 0; i <= 1; i++ {
		speaker.Lock()
		audioDataQueue.Add(streamer)
		speaker.Unlock()
		time.Sleep(10 * time.Second)
		streamer.Seek(0)
	}
	streamer.Close()
	os.Remove(testInput)
}
