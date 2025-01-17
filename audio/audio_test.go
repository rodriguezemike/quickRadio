package audio

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func GetTestFileObject(desiredFilename string) *os.File {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filepath.Dir(filename))
	testFilePath := filepath.Join(dir, "assets", "tests", desiredFilename)
	fileObject, _ := os.Open(testFilePath)
	return fileObject
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

func TestDownloadAACs(t *testing.T) {
	//This is a bit difficule to test as we'd need some random live game to test downloading and transcoding AACS to wavs
	// Ready for play.
	//This might be the only way to do it and if so we'd need to handle that gracefully, minimizing the amount of errors.
}
func TestGetQualityStreamSlug(t *testing.T) {
	want := "lak-radio_192K.m3u8"
	wantedAudioQuality := "192K"
	testFile := GetTestFileObject("lak-radio.m3u8")
	byteValue, _ := io.ReadAll(testFile)
	contents := string(byteValue)
	qualityString, err := GetQualityStreamSlug(contents, wantedAudioQuality)
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
	testFile := GetTestFileObject("lak-radio_192k.m3u8")
	byteValue, _ := io.ReadAll(testFile)
	contents := string(byteValue)
	audioFiles, err := GetAudioFiles(contents)
	if sliceCompare(audioFiles, want) == false {
		t.Fatalf(`GetQualityStreamSlug(contents, wantedAudioQuality) = %q, %v, want match for %#q, nil`, audioFiles, err, want)
	}
}
