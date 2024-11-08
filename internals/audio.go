package internals

import (
	"fmt"
	"strings"
	"time"

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

//This will become one of the go rountines

func TranscodeToWave(aacFilepath string) string {
	wavFilepath := strings.Replace(aacFilepath, ".aac", ".wav", 1)
	err := ffmpeg_go.Input(aacFilepath).Output(wavFilepath).OverWriteOutput().ErrorToStdOut().Run()
	ErrorCheck(err)
	return wavFilepath
}
func playWaveFile(wavFilePath string) string {
	return ""
}

func playback(landingLink string, radioLink string) {
	//This can be changed to gameState = 'Completed'
	//cloudflareBase - https://d2igy0yla8zi0u.cloudfront.net/lak/20242025/
	//Radio quality
	//Example to Extract : lak-radio_192K.m3u8
	//Then we extract that https://d2igy0yla8zi0u.cloudfront.net/lak/20242025/lak-radio_192K.m3u8
	//Then we get our aac files
	//lak-radio_192K/00021/lak-radio_192K_00118.aac
	cloudflareBase := strings.Join(strings.Split(radioLink, "/")[:len(strings.Split(radioLink, "/"))-1], "/")
	radioQualityLink, err := GetQualityStreamSlug(radioLink)
	ErrorCheck(err)
	streamLink := cloudflareBase + "/" + radioQualityLink
	for true == true {
		//While were running in the game loop
		//We to sleep main thread, really this would update the UI
		//Looking for updates in the game landing for game status, time, etc.
		//check for new undownloaded radio files
		//download those
		//Play Any AAC file were missing.
		//Delete old AAC files
		//Update our data structs.
		time.Sleep(1 * time.Second)
		fmt.Println(streamLink)
	}
}
