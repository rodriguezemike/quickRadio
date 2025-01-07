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

// For Testing
func playWaveFile(wavFilePath string) {
	streamer := InitializeRadio(wavFilePath)
	PlayWave(streamer)
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
	streamer, format := DecodeWaveFile(wavFilePath)
	InitalizeRadioSpeaker(format)
	return streamer
}

func PlayWave(streamer beep.StreamSeekCloser) {
	done := make(chan bool)
	defer streamer.Close()
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done
}

func PlayRadio() {
	//Built from the TestFunc TestStream in audio data queue
	//We're Initalizing it
	//Then checking to see if we have files to play in the tmp dir
	//Decoding them
	//Adding them to the queue
	//And playing them for the total time we have in a sleep for this go routine
	//Or half the sleep then decode and wait in a infinite loop checking to see
	//The position of the streamer, when iti s done, lock the speaker and add the streamers
	//Then do it again.
	//Then adding more streamers
}
