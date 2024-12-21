package audio

import (
	"errors"
	"os"
	"os/exec"
	"quickRadio/radioErrors"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gopxl/beep"
	"github.com/gopxl/beep/speaker"
	"github.com/gopxl/beep/wav"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

//For now, work with the sync group, this may not work
//As they may not be a delay in all of this and so we may get breaks
//Unless we want to add in a buffer where we download first and then
//wait and then always be 1 whole file behind in playback.
//This does happen in mobile vs desktop apps.
//When we return this, we may want to wait until the whole manifest is done
//Looking at the polling pattern it looks like it cracks open the file every time.
//We can sort this out this in testing and retieration.

func DownloadAndTranscodeAACs(paths []string) []string {
	wavpaths := make([]string, len(paths))
	var workGroup sync.WaitGroup
	for i := 0; i < len(paths); i++ {
		workGroup.Add(1)
		//We want to return an error maybe if something goes wrong?
		go func(path string) {
			defer workGroup.Done()
			localpath, err := DownloadAAC(path)
			radioErrors.ErrorCheck(err)
			wavpaths = append(wavpaths, TranscodeToWave(localpath))
		}(paths[i])
		workGroup.Wait()
	}
	return wavpaths
}

// Only call this once when setting up the radio
func CreateTmpDirectory(tmpDirectory string) {
	if _, err := os.Stat(tmpDirectory); err == nil {
		return
	} else if os.IsNotExist(err) {
		err := os.MkdirAll(tmpDirectory, 0777)
		radioErrors.ErrorCheck(err)
		return
	} else {
		radioErrors.ErrorCheck(err)
		return
	}
}

func DoesFileExist(filepath string) error {
	if _, err := os.Stat(filepath); err == nil {
		return nil
	} else if os.IsNotExist(err) {
		return errors.New("file does not exist")
	} else {
		radioErrors.ErrorCheck(err)
		return err
	}
}

func DownloadAAC(aacRequestPath string) (string, error) {
	//We want to Move the base os tmp directories into the links.json (Program config.)
	filename := strings.Split(aacRequestPath, "/")[len(strings.Split(aacRequestPath, "/"))-1]
	if runtime.GOOS == "windows" {
		directory := strings.Join(strings.Split(aacRequestPath, "/")[4:], "\\")
		filepath := "C:\\Users\\AppData\\Local\\Temp" + "\\" + "QuickRadio" + "\\" + directory + "\\" + filename
		cmd := exec.Command("wget", aacRequestPath, "-OutFile", filepath)
		_, err := cmd.CombinedOutput()
		radioErrors.ErrorCheck(err)
		err = DoesFileExist(filepath)
		radioErrors.ErrorCheck(err)
		return filepath, nil
	} else if runtime.GOOS == "linux" {
		directory := strings.Join(strings.Split(aacRequestPath, "/")[4:], "/")
		filepath := "/tmp/" + "QuickRadio" + "/" + directory + "/" + filename
		cmd := exec.Command("wget", aacRequestPath, filepath)
		_, err := cmd.CombinedOutput()
		radioErrors.ErrorCheck(err)
		err = DoesFileExist(filepath)
		radioErrors.ErrorCheck(err)
		return filepath, nil
	} else { // Assume Unix default
		directory := strings.Join(strings.Split(aacRequestPath, "/")[4:], "/")
		filepath := "/tmp/" + "QuickRadio" + "/" + directory + "/" + filename
		cmd := exec.Command("wget", aacRequestPath, filepath)
		_, err := cmd.CombinedOutput()
		radioErrors.ErrorCheck(err)
		err = DoesFileExist(filepath)
		radioErrors.ErrorCheck(err)
		return filepath, nil
	}
}

// This will become one of the go rountines
func TranscodeToWave(aacFilepath string) string {
	wavFilepath := strings.Replace(aacFilepath, ".aac", ".wav", 1)
	err := ffmpeg_go.Input(aacFilepath).Output(wavFilepath).OverWriteOutput().ErrorToStdOut().Run()
	radioErrors.ErrorCheck(err)
	return wavFilepath
}

// This is a seq func, that will be broken up in its parallel go routines
// I think were going to want to keep a timer to transode as many as possible
// That is if len of our beep sequence is 5 then we have 50 seconds
// At the end of that 50 we need to play a new sequence, with a done.
// We can have multiple speakers, but I'd prefer to have multiple sequences or have a callback
// that pulls in the next sequence of streamers/callback

//We may lose speaker in some sort of scope issue. Keep an eye for out that./

// For Testing
func playWaveFile(wavFilePath string) {
	streamer := InitializeRadio(wavFilePath)
	PlayRadio(streamer)
}

func DecodeWaveFile(wavFilePath string) (beep.StreamSeekCloser, beep.Format) {
	f, err := os.Open(wavFilePath)
	radioErrors.ErrorCheck(err)
	streamer, format, err := wav.Decode(f)
	radioErrors.ErrorCheck(err)
	return streamer, format
}

func InitalizeRadioSpeaker(format beep.Format) {
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}

func InitializeRadio(wavFilePath string) beep.StreamSeekCloser {
	//This func should hold a bunch of other options if desired.
	streamer, format := DecodeWaveFile(wavFilePath)
	InitalizeRadioSpeaker(format)
	return streamer
}

// For this we are going to follow Queue Example On, adding streamers to the list in order
// This should work. Adding files as long as we have em or until the user presses a button.
// Reason for this is the Streamers need to be dynamic.
// https://github.com/gopxl/beep/wiki/Making-own-streamers
func PlayRadio(streamer beep.StreamSeekCloser) {
	done := make(chan bool)
	defer streamer.Close()
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		//In this callback, we want to either add to the current beep sequence by
		//Checking if we have more files to decode and decoding them adding the sequence
		//To the play. We want to decode in parallel and add in serial. So this just
		//wants pull the next streamer and play
		//If we hit a 'done' state from the game landing then we are done
		//And we can set our done bool to True.
		done <- true
	})))
	<-done
}
