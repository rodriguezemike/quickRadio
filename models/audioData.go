package models

import "github.com/gopxl/beep"

//This will hold all the streams and make it easier to maniuplate filepaths,
//Streamers, making a Sequence and handle its callback.
//We may want this to hold a buffer instead and move into buffing.
//This file just holds the idea for now.
//But should have something like addToBUffer or addToSequence
//etc.
type AudioData struct {
	streams []beep.StreamSeekCloser
}
