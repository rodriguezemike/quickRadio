package audio

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"quickRadio/models"
	"quickRadio/quickio"
	"strings"
	"sync"
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
	goroutineMap           sync.Map
}

func (controller *RadioController) initalizeRadioSpeaker(format beep.Format) {
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	controller.speakerInitialized = true
}

func (controller *RadioController) PollRadioFormatLink() {
	log.Println("readioController::PollRadioFormatLink")
	log.Println("PollRadioFormatLink::controller.RadioFormatLink", controller.RadioFormatLink)
	ctx, cancel := context.WithDeadline(controller.ctx, time.Now().Add(10*time.Second))
	defer cancel()
	go quickio.UpdateRadioWavsWithContext(ctx, controller.RadioFormatLink)
	controller.goroutineMap.Store(ctx, cancel)
	log.Println("PollRadioFormatLink::Sleeping for 10 seconds to let it do its work")
	time.Sleep(10 * time.Second)
	controller.goroutineMap.Delete(ctx)
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
	log.Println("RadioController::updateSharedData")
	log.Println("RadioController::updateSharedData::Controller Wavpaths ->", controller.WavPaths)
	log.Println("RadioController::updateSharedData::Controller Streamers ->", controller.streamers)
	latestFileInfo, _ := os.Stat(controller.WavPaths[len(controller.WavPaths)-1])
	files, _ := os.ReadDir(controller.RadioDirectory)
	for _, f := range files {
		info, _ := f.Info()
		if info.ModTime().After(latestFileInfo.ModTime()) && strings.HasSuffix(f.Name(), ".wav") {
			streamer, _ := controller.decodeWaveFile(filepath.Join(controller.RadioDirectory, f.Name()))
			streamers = append(streamers, streamer)
			wavPaths = append(wavPaths, filepath.Join(controller.RadioDirectory, f.Name()))
		}
	}
	controller.streamers = streamers
	if len(wavPaths) != 0 {
		controller.WavPaths = wavPaths
	}
}

func (controller *RadioController) initalPlayback() {
	streamer := controller.initializeRadio(controller.WavPaths[len(controller.WavPaths)-1])
	speaker.Play(&controller.streamQueue)
	controller.addStreamerToPlaybackQueue(streamer)
	go controller.PollRadioFormatLink()
}

func (controller *RadioController) updatePlayback() {
	log.Println("RadioController::updatePlayback")
	for _, streamer := range controller.streamers {
		log.Println("RadioController::updatePlayback::Playing Streamer ", streamer)
		controller.addStreamerToPlaybackQueue(streamer)
	}
}

func (controller *RadioController) stopRadio() {
	speaker.Clear()
	quickio.DeleteRadioLock()
	controller.goroutineMap.Range(func(key, value interface{}) bool {
		callback, _ := value.(context.CancelFunc)
		callback()
		return true
	})
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
				log.Println("PlayRadio::POLLING RADIO FORMAT LINK")
				go controller.PollRadioFormatLink()
				time.Sleep(time.Duration(controller.EmergencySleepInterval) * time.Second)
				continue
			}
			controller.updatePlayback()
			go controller.PollRadioFormatLink()
			time.Sleep(time.Duration(controller.NormalSleepInterval) * time.Second)
		}
	}
}

func StartRadioFun(radioLink string) {
	controller := NewRadioController(radioLink)
	controller.PlayRadio()
}

func StopRadioFun() {
	speaker.Clear()
	quickio.DeleteRadioLock()
}

func NewRadioController(radioLink string) RadioController {
	log.Println("Func -> NewRadioController")
	var controller RadioController
	controller.RadioFormatLink, controller.RadioDirectory, controller.WavPaths = quickio.GetRadioFormatLinkAndDirectory(radioLink)
	controller.NormalSleepInterval = 3
	controller.EmergencySleepInterval = 1
	controller.ctx = context.Background()
	controller.speakerInitialized = false
	log.Println("NewRadioController::controller.WavPaths ", controller.WavPaths)
	log.Println("NewRadioController::controller.RadioDirectory ", controller.RadioDirectory)
	log.Println("NewRadioController::controller.RadioFormatLink ", controller.RadioFormatLink)
	return controller
}

func RadioKillFun() {
	quickio.DeleteRadioLock()
	time.Sleep(3 * time.Second)
	quickio.EmptyTmpFolder()
}
