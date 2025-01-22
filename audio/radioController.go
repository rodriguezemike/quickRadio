package audio

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
)

type RadioController struct {
	WavPaths               []string
	RadioDirectory         string
	RadioFormatLink        string
	EmergencySleepInterval int
	NormalSleepInterval    int
	streamQueue            models.AudioStreamQueue
	streamers              []beep.StreamSeekCloser
	speakerInitialized     bool
	ctx                    context.Context
	cancelFunc             context.CancelFunc
}

func (controller *RadioController) initalizeRadioSpeaker(format beep.Format) {
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	controller.speakerInitialized = true
}

func (controller *RadioController) PollRadioFormatLink() {
	ctx, cancel := context.WithDeadline(controller.ctx, time.Now().Add(10*time.Second))
	defer cancel()
	go quickio.UpdateRadioWavsWithContext(ctx, controller.RadioFormatLink)
	controller.cancelFunc = cancel
}

func (controller *RadioController) addStreamerToPlaybackQueue(streamer beep.StreamSeekCloser) {
	speaker.Lock()
	controller.streamQueue.Add(streamer)
	speaker.Unlock()
}
func (controller *RadioController) decodeWaveFile(wavFilepath string) (beep.StreamSeekCloser, beep.Format) {
	wavFilepath = strings.TrimSpace(wavFilepath)
	f, _ := os.Open(wavFilepath)
	streamer, format, _ := wav.Decode(f)
	return streamer, format
}

func (controller *RadioController) initializeRadio(wavFilePath string) beep.StreamSeekCloser {
	streamer, format := controller.decodeWaveFile(wavFilePath)
	controller.initalizeRadioSpeaker(format)
	return streamer
}

func (controller *RadioController) updateSharedData() {
	var streamers []beep.StreamSeekCloser
	var wavPaths []string
	latestFileInfo, _ := os.Stat(controller.WavPaths[len(controller.WavPaths)-1])
	files, _ := os.ReadDir(controller.RadioDirectory)
	for _, f := range files {
		info, _ := f.Info()
		if info.ModTime().After(latestFileInfo.ModTime()) && strings.HasSuffix(f.Name(), ".wav") {
			streamer, _ := DecodeWaveFile(filepath.Join(controller.RadioDirectory, f.Name()))
			streamers = append(streamers, streamer)
			wavPaths = append(wavPaths, filepath.Join(controller.RadioDirectory, f.Name()))
		}
	}
	controller.streamers = streamers
	controller.WavPaths = wavPaths
}

func (controller *RadioController) initalPlayback() {
	streamer := controller.initializeRadio(controller.WavPaths[len(controller.WavPaths)-1])
	speaker.Play(&controller.streamQueue)
	controller.addStreamerToPlaybackQueue(streamer)
	controller.PollRadioFormatLink()
}

func (controller *RadioController) updatePlayback() {
	log.Println("FUNC - UPDATE PLAYBACK")
	for _, streamer := range controller.streamers {
		log.Println("Playing Streamer ", streamer)
		controller.addStreamerToPlaybackQueue(streamer)
	}
}

func (controller *RadioController) stopRadio() {
	speaker.Clear()
	quickio.DeleteRadioLock()
	controller.cancelFunc()
	time.Sleep(1 * time.Second)
	quickio.EmptyRadioDirectory(controller.RadioDirectory)
}

func (controller *RadioController) PlayRadio() {
	if !quickio.IsRadioLocked() {
		quickio.CreateRadioLock()
		controller.initalPlayback()
		time.Sleep(time.Duration(controller.NormalSleepInterval) * time.Second)
		for {
			if !quickio.IsRadioLocked() {
				controller.stopRadio()
				return
			}
			controller.updateSharedData()
			if len(controller.streamers) == 0 {
				controller.PollRadioFormatLink()
				time.Sleep(time.Duration(controller.EmergencySleepInterval) * time.Second)
				continue
			}
			controller.updatePlayback()
		}
	}
}

func StartRadioFun(radioLink string) {
	controller := NewRadioController(radioLink)
	controller.PlayRadio()
}

func StopRadioFun() {
	quickio.DeleteRadioLock()
}

func NewRadioController(radioLink string) RadioController {
	var controller RadioController
	controller.RadioFormatLink, controller.RadioDirectory, controller.WavPaths = quickio.GetRadioFormatLinkAndDirectory(radioLink)
	controller.NormalSleepInterval = 3
	controller.EmergencySleepInterval = 1
	controller.ctx = context.Background()
	controller.speakerInitialized = false
	return controller
}

func RadioKillFun() {
	quickio.DeleteRadioLock()
	time.Sleep(5 * time.Second)
	quickio.EmptyTmpFolder()
}
