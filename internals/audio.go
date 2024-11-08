package internals

import (
	"os"
	"strings"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

/*
Download aac files and transcode aac to wav using ffmpeg-go - can tweak wav files
We want the file names to be the same
All this needs to be done in parallel per m3u8 file. Whatever its called.
Once we have the wav files were going to natural sort them
Once naturally sorted
We extend our list of files to play
We playback using beep. - Can tweak playback options

Dev Day 1
Focus on aac to wav transcode and audio playback

*/

// We should have a shared array here maybe? Depends on goRoutines

// This will become one of the go rountines
func TranscodeToWave(aacFilepath string) string {
	wavFilepath := strings.Replace(aacFilepath, ".aac", ".wav", 1)
	err := ffmpeg_go.Input(aacFilepath).Output(wavFilepath).OverWriteOutput().ErrorToStdOut().Run()
	ErrorCheck(err)
	return wavFilepath
}

// This is a seq func, that will be broken up in its parallel go routines
// I think were going to want to keep a timer to transode as many as possible
// That is if len of our beep sequence is 5 then we have 50 seconds
// At the end of that 50 we need to play a new sequence, with a done.
// We can have multiple speakers, but I'd prefer to have multiple sequences or have a callback
// that pulls in the next sequence of streamers/callback
func playWaveFile(wavFilePath string) {
	done := make(chan bool)
	f, err := os.Open(wavFilePath)
	ErrorCheck(err)
	streamer, format, err := wav.Decode(f)
	ErrorCheck(err)
	defer streamer.Close()
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}
