package audio

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"runtime"
	"testing"
	"time"

	"github.com/gopxl/beep/speaker"
)

func TestStream(t *testing.T) {
	var audioDataQueue models.AudioStreamQueue
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	testInput := filepath.Join(dir, "assets", "tests", "game_192k_00001.wav")
	if _, err := os.Stat(testInput); errors.Is(err, os.ErrNotExist) {
		testAAC := filepath.Join(dir, "assets", "tests", "game_192k_00001.aac")
		testInput = quickio.TranscodeToWave(testAAC)
	}
	streamer := InitializeRadio(testInput)
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

func TestTranscodeToWave(t *testing.T) {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	testInput := filepath.Join(dir, "assets", "tests", "game_192k_00001.aac")
	wanted := filepath.Join(dir, "assets", "tests", "game_192k_00001.wav")
	result := quickio.TranscodeToWave(testInput)
	if _, err := os.Stat(wanted); errors.Is(err, os.ErrNotExist) || result != wanted {
		t.Fatalf(`TranscodeToWave(assets/tests/game_192k_00001.aac) = %q, %v, want match for %#q, nil`, result, err, wanted)
	}
	os.Remove(result)
}

func TestExtractQualityStreamSlug(t *testing.T) {
	want := "lak-radio_192K.m3u8"
	wantedAudioQuality := "192K"
	testFile := quickio.GetTestFileObject("lak-radio.m3u8")
	byteValue, _ := io.ReadAll(testFile)
	contents := string(byteValue)
	qualityString, err := quickio.ExtractQualityStreamSlug(contents, wantedAudioQuality)
	if want != qualityString {
		t.Fatalf(`GetQualityStreamSlug(contents, wantedAudioQuality) = %q, %v, want match for %#q, nil`, qualityString, err, want)
	}
}

func TestGetAudioFiles(t *testing.T) {
	want := []string{"lak-radio_192K/00021/lak-radio_192K_00105.aac",
		"lak-radio_192K/00021/lak-radio_192K_00106.aac",
		"lak-radio_192K/00021/lak-radio_192K_00107.aac",
		"lak-radio_192K/00021/lak-radio_192K_00108.aac",
		"lak-radio_192K/00021/lak-radio_192K_00109.aac",
		"lak-radio_192K/00021/lak-radio_192K_00110.aac",
		"lak-radio_192K/00021/lak-radio_192K_00111.aac",
		"lak-radio_192K/00021/lak-radio_192K_00112.aac",
		"lak-radio_192K/00021/lak-radio_192K_00113.aac",
		"lak-radio_192K/00021/lak-radio_192K_00114.aac"}
	testFile := quickio.GetTestFileObject("lak-radio_192k.m3u8")
	byteValue, _ := io.ReadAll(testFile)
	contents := string(byteValue)
	audioFiles, err := quickio.GetAACSlugsFromQualityFile(contents)
	if sliceCompare(audioFiles, want) == false {
		t.Fatalf(`GetQualityStreamSlug(contents, wantedAudioQuality) = %q, %v, want match for %#q, nil`, audioFiles, err, want)
	}
}
