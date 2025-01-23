package audio

import (
	"log"
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

func AddStreamerToPlaybackQueue(audioPlaybackQueue *models.AudioStreamQueue, streamer beep.StreamSeekCloser, qualityLink string) {
	speaker.Lock()
	audioPlaybackQueue.Add(streamer)
	speaker.Unlock()
	go quickio.UpdateRadioWavs(qualityLink)
	time.Sleep(3 * time.Second)

}

func InitalPlayback(wavPaths []string, audioPlaybackQueue *models.AudioStreamQueue, qualityLink string) string {
	var radioDirectory string
	speakerInitalized := false
	for _, wavPath := range wavPaths[len(wavPaths)-2:] {
		if !quickio.IsRadioLocked() {
			//Call Funcs to clean up any gorountines
			return ""
		}
		if strings.HasSuffix(wavPath, ".wav") {
			streamer, format := DecodeWaveFile(wavPath)
			if !speakerInitalized {
				InitalizeRadioSpeaker(format)
				speaker.Play(audioPlaybackQueue)
				speakerInitalized = true
				radioDirectory = filepath.Dir(wavPath)
			}
			AddStreamerToPlaybackQueue(audioPlaybackQueue, streamer, qualityLink)
		}
	}
	return radioDirectory

}

func UpdatePlayback(audioPlaybackQueue *models.AudioStreamQueue, audioStreamers []beep.StreamSeekCloser, qualityLink string) {
	log.Println("FUNC - UPDATE PLAYBACK")
	for _, streamer := range audioStreamers {
		log.Println("Playing Streamer ", streamer)
		AddStreamerToPlaybackQueue(audioPlaybackQueue, streamer, qualityLink)
	}
}

func PlayRadio(wavPaths []string, qualityLink string) {
	var audioPlaybackQueue models.AudioStreamQueue
	var audioStreamers []beep.StreamSeekCloser
	if !quickio.IsRadioLocked() {
		quickio.CreateRadioLock()
		radioDirectory := InitalPlayback(wavPaths, &audioPlaybackQueue, qualityLink)
		for {
			if !quickio.IsRadioLocked() {
				//Here we sohould gather up the gorountines and end them gracefully.
				return
			}
			audioStreamers, wavPaths = UpdateSharedData(radioDirectory, wavPaths)
			log.Println("Audio Streamers After Update ", audioStreamers)
			if len(audioStreamers) == 0 {
				go quickio.UpdateRadioWavs(qualityLink)
				time.Sleep(1 * time.Second)
				continue
			}
			UpdatePlayback(&audioPlaybackQueue, audioStreamers, qualityLink)
		}
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
	quickio.DeleteRadioLock()
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
