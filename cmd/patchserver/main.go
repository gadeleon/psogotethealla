package main

import (
	"log"

	patchserver "github.com/gadeleon/psogotethealla/patch_server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Load config
	// TODO: Option Argument instead of argv

	log.Printf("Starting GoTethealla v%s \n", patchserver.SERVER_VERSION)
	patchserver.StartServer()
}
