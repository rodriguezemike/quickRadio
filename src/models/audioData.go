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
	volume    float64
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

func (streamQueue *AudioStreamQueue) SetVolume(volume float64) {
	if volume < 0.0 {
		volume = 0.0
	}
	if volume > 2.0 { // Allow amplification up to 200%
		volume = 2.0
	}
	streamQueue.volume = volume
}

func (streamQueue *AudioStreamQueue) GetVolume() float64 {
	return streamQueue.volume
}

func (streamQueue *AudioStreamQueue) Stream(samples [][2]float64) (n int, ok bool) {
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

		// Apply volume adjustment to the streamed samples
		for i := streamed; i < streamed+n; i++ {
			samples[i][0] *= streamQueue.volume
			samples[i][1] *= streamQueue.volume
		}

		streamed += n
	}
	return len(samples), true
}

func NewAudioStreamQueue() *AudioStreamQueue {
	return &AudioStreamQueue{
		streamers: make([]beep.Streamer, 0),
		volume:    1.0, // Default volume at 100%
	}
}
