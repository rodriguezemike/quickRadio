package models

import "github.com/gopxl/beep"

//This will hold all the streams and make it easier to maniuplate filepaths,
//Streamers, making a Sequence and handle its callback.
//We may want this to hold a buffer instead and move into buffing.
//This file just holds the idea for now.
//But should have something like addToBUffer or addToSequence
//etc.
type AudioStreamQueue struct {
	streamers []beep.Streamer
}

func (streamQueue *AudioStreamQueue) Err() error {
	return nil
}

func (streamQueue *AudioStreamQueue) Add(streamers ...beep.Streamer) {
	streamQueue.streamers = append(streamQueue.streamers, streamers...)
}

func (streamQueue *AudioStreamQueue) PopStream() {
	streamQueue.streamers = streamQueue.streamers[1:]
}

func (streamQueue *AudioStreamQueue) stream(samples [][2]float64) (n int, ok bool) {
	streamed := 0
	for streamed < len(samples) {
		if len(streamQueue.streamers) == 0 {
			for i := range samples[streamed:] {
				samples[i][0] = 0
				samples[i][0] = 0
			}
			break
		}
		n, ok := streamQueue.streamers[0].Stream(samples[streamed:])
		if !ok {
			streamQueue.PopStream()
		}
		streamed += n
		return len(samples), true
	}
}
