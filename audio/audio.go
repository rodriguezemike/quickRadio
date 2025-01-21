package audio

import (
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"quickRadio/radioErrors"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

func DecodeWaveFile(wavFilepath string) (beep.StreamSeekCloser, beep.Format) {
	wavFilepath = strings.TrimSpace(wavFilepath)
	f, err := os.Open(wavFilepath)
	radioErrors.ErrorCheck(err)
	streamer, format, err := wav.Decode(f)
	radioErrors.ErrorCheck(err)
	return streamer, format
}

func InitalizeRadioSpeaker(format beep.Format) {
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}

func InitializeRadio(wavFilePath string) beep.StreamSeekCloser {
	streamer, format := DecodeWaveFile(wavFilePath)
	InitalizeRadioSpeaker(format)
	return streamer
}

func AddStreamerToPlaybackQueue(audioPlaybackQueue models.AudioStreamQueue, streamer beep.StreamSeekCloser, qualityLink string) models.AudioStreamQueue {
	speaker.Lock()
	audioPlaybackQueue.Add(streamer)
	speaker.Unlock()
	go quickio.UpdateRadioWavs(qualityLink)
	return audioPlaybackQueue
}

func InitalPlayback(wavPaths []string, audioPlaybackQueue models.AudioStreamQueue, qualityLink string) (models.AudioStreamQueue, string) {
	var radioDirectory string
	speakerInitalized := false
	for _, wavPath := range wavPaths {
		if strings.HasSuffix(wavPath, ".wav") {
			streamer, format := DecodeWaveFile(wavPath)
			if !speakerInitalized {
				InitalizeRadioSpeaker(format)
				speaker.Play(&audioPlaybackQueue)
				speakerInitalized = true
				radioDirectory = filepath.Dir(wavPath)
			}
			audioPlaybackQueue = AddStreamerToPlaybackQueue(audioPlaybackQueue, streamer, qualityLink)
			time.Sleep(3 * time.Second)
		}
	}
	return audioPlaybackQueue, radioDirectory

}

func UpdatePlayback(audioPlaybackQueue models.AudioStreamQueue, wavPaths []string, radioDirectory string, qualityLink string) (models.AudioStreamQueue, []string, string) {
	var audioStreamers []beep.StreamSeekCloser
	audioStreamers, wavPaths = UpdateSharedData(radioDirectory, wavPaths)
	if len(audioStreamers) == 0 {
		go quickio.UpdateRadioWavs(qualityLink)
		time.Sleep(10 * time.Second)
	} else {
		for _, streamer := range audioStreamers {
			AddStreamerToPlaybackQueue(audioPlaybackQueue, streamer, qualityLink)
			time.Sleep(3 * time.Second)
		}
	}
	return audioPlaybackQueue, wavPaths, radioDirectory
}

func PlayRadio(wavPaths []string, qualityLink string) {
	var audioPlaybackQueue models.AudioStreamQueue

	audioPlaybackQueue, radioDirectory := InitalPlayback(wavPaths, audioPlaybackQueue, qualityLink)
	for {
		audioPlaybackQueue, wavPaths, radioDirectory = UpdatePlayback(audioPlaybackQueue, wavPaths, radioDirectory, qualityLink)
	}
}

func UpdateSharedData(radioDirectory string, wavPaths []string) ([]beep.StreamSeekCloser, []string) {
	var streamers []beep.StreamSeekCloser
	latestFileInfo, err := os.Stat(wavPaths[len(wavPaths)-1])
	radioErrors.ErrorCheck(err)
	files, err := os.ReadDir(radioDirectory)
	radioErrors.ErrorCheck(err)
	for _, f := range files {
		info, err := f.Info()
		radioErrors.ErrorCheck(err)
		if info.ModTime().After(latestFileInfo.ModTime()) && strings.HasSuffix(f.Name(), ".wav") {
			streamer, _ := DecodeWaveFile(filepath.Join(radioDirectory, f.Name()))
			streamers = append(streamers, streamer)
			wavPaths = append(wavPaths, filepath.Join(radioDirectory, f.Name()))
		}
	}
	return streamers, wavPaths
}

func StopRadio() {
	speaker.Clear()
}

func KillFun() {
	StopRadio()
	speaker.Close()
	quickio.EmptyTmpFolder()
}

func StartFun(radioLink string) {
	if radioLink != "" {
		//Check where we're at as we may want to resume fun instead of starting over
		qualitySlug := quickio.GetQualityStreamSlug(radioLink, "192K")
		//ToDo: Once we have the Slug we want to build the name of the link we need to grap
		qualityRadioPath := quickio.BuildQualityRadioPath(radioLink, qualitySlug)
		aacPaths := quickio.GetAACPaths(qualityRadioPath)
		wavPaths := quickio.DownloadAndTranscodeAACs(aacPaths)
		PlayRadio(wavPaths, qualityRadioPath)
	}
}
