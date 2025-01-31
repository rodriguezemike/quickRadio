package main

import (
	"log"
	"quickRadio/views"
)

func main() {
	ui := views.NewQuickRadioView()
	log.Println("In Main line 9")
	ui.CreateAndRunApp()
}
