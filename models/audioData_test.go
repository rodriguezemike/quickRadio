package models

import "testing"

func TestAudioData(t *testing.T) {
	var audioQueue = &AudioStreamQueue{}
	if audioQueue == nil {
		t.Fatalf(`Error creating struct AudioStreamQueue.`)
	}
}
