package internals

import (
	"fmt"
	"strings"
	"time"
)

// We should have a shared array here maybe? Depends on goRoutines

func playback(landingLink string, radioLink string) {
	//This can be changed to gameState = 'Completed'
	//cloudflareBase - https://d2igy0yla8zi0u.cloudfront.net/lak/20242025/
	//Radio quality
	//Example to Extract : lak-radio_192K.m3u8
	//Then we extract that https://d2igy0yla8zi0u.cloudfront.net/lak/20242025/lak-radio_192K.m3u8
	//Then we get our aac files
	//lak-radio_192K/00021/lak-radio_192K_00118.aac
	cloudflareBase := strings.Join(strings.Split(radioLink, "/")[:len(strings.Split(radioLink, "/"))-1], "/")
	radioQualityLink := GetQualityStreamSlug(radioLink)
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
