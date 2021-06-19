package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gadeleon/psogotethealla/config"
	patchserver "github.com/gadeleon/psogotethealla/patch_server"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Load config
	// TODO: Option Argument instead of argv

	log.Printf("Starting GoTethealla v%s \n", patchserver.SERVER_VERSION)
	cnf, err := config.New(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	if ext := cnf.Config.Section("login_server").Key("external_ip").String(); ext != "" {
		log.Printf("External ip set, binding to %s", ext)

	}
	log.Printf("Ship IP Address: %s", cnf.Config.Section("login_server").Key("server").String())
	// Setup Patch & Data Port
	loginPort := cnf.Config.Section("login_server").Key("port").String()
	p, err := strconv.Atoi(loginPort)
	if err != nil {
		log.Fatal(err)
	}
	patchPort := uint16(p - 1000)
	dataPort := uint16(p - 999)

	log.Printf("Patch Port: %d", patchPort)
	log.Printf("Data Port: %d", dataPort)
	log.Printf("Max Connections: %s", cnf.Config.Section("login_server").Key("maxconn").String())
	// Setup Max Upload Speed in bytes
	mup := cnf.Config.Section("patch_server").Key("maxup").String()
	m, err := strconv.Atoi(mup)
	if err != nil {
		log.Fatal(err)
	}
	maxupbytes := uint32(patchserver.BytesToKB(m))
	log.Printf("Max Upload Speed: %sKB/s (%dBytes/s)", cnf.Config.Section("patch_server").Key("maxup").String(), maxupbytes)

	// TODO: Setup Patch Data Folder
	// TODO: Have this folder be a full path in INI?
	//var ch, ch2 uint32

}
