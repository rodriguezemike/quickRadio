package audio

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestTranscodeToWave(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	testInput := filepath.Join(dir, "assets", "tests", "game_192k_00001.aac")
	wanted := filepath.Join(dir, "assets", "tests", "game_192k_00001.wav")
	result := TranscodeToWave(testInput)
	if _, err := os.Stat(wanted); errors.Is(err, os.ErrNotExist) || result != wanted {
		t.Fatalf(`TranscodeToWave(assets/tests/game_192k_00001.aac) = %q, %v, want match for %#q, nil`, result, err, wanted)
	}
	os.Remove(result)
}

func TestPlayWaveFile(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	testInput := filepath.Join(dir, "assets", "tests", "game_192k_00001.wav")
	if _, err := os.Stat(testInput); errors.Is(err, os.ErrNotExist) {
		testAAC := filepath.Join(dir, "assets", "tests", "game_192k_00001.aac")
		testInput = TranscodeToWave(testAAC)
	}
	if _, err := os.Stat(testInput); errors.Is(err, os.ErrNotExist) {
		t.Fatalf(`Our Test Input Wav file was not generated. Check err %s`, err)
	}
	playWaveFile(testInput)
	os.Remove(testInput)
}
