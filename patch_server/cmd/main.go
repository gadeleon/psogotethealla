package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	core "github.com/gadeleon/psogotethealla/client"
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
	maxupbytes := uint32(patchserver.KBToBytes(m))
	log.Printf("Max Upload Speed: %sKB/s (%dBytes/s)", cnf.Config.Section("patch_server").Key("maxup").String(), maxupbytes)

	// TODO: Setup Patch Data Folder
	// TODO: Have this folder be a full path in INI?
	//var ch, ch2 uint32

	log.Print("Setting up patch data")
	patches := patchserver.NewPatchData(cnf.Config.Section("patch_server").Key("directory").Value())
	//patches := []patchserver.PatchData{}
	log.Printf("Created patch struct: %v", patches)

	if len(patches) < 1 {
		log.Fatal("No patches found. At least one patch to check or send is required.")
	}

	log.Print("Loading welcome message...")
	welcome := cnf.Config.Section("login_server").Key("welcome").Value()
	if welcome == "" {
		log.Fatal("Welcome message is empty. Yes, this is required...")
	}
	log.Print(welcome)

	// Start Patch Server
	patchAddr := fmt.Sprintf("%s:%d", cnf.Config.Section("login_server").Key("server").Value(), patchPort)
	l, err := net.Listen("tcp", patchAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	log.Printf("Listening on %s", patchAddr)
	for {
		// Wait for connection
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// Handle Connection, Print Contents
		go func(c net.Conn) {
			defer c.Close()
			raddr, _, _ := net.SplitHostPort(c.RemoteAddr().String())
			log.Printf("%v", c)
			workingConnection := &patchserver.PatchClient{
				Core:      &core.CoreClient{},
				IpAddress: net.ParseIP(raddr),
				Patch:     0,
			}
			log.Printf("Incoming connection from: %s", c.RemoteAddr().String())
			patchserver.StartEncryption(workingConnection)
			log.Print(workingConnection.IpAddress, workingConnection.Patch)

			for {
				buff := make([]byte, 4096)
				c.Read(buff)
				log.Print(buff)
				log.Print("Doing stufF!")
				log.Print("Heeeeey")
				//log.Print(b)
				log.Printf("%v", c)
				break
				// if b[0] == 97 {
				// 	log.Print("Closing")
				// 	c.Close()
				// 	break
				// }
				// if b == "bye" || b == "bye\n" {
				// 	fmt.Println("Okay, bye")
				// 	log.Print("Closed connection.")
				// 	c.Close()
				// 	break
				// }
			}
			// p := make([]byte, 100)
			//io.Copy(c, c)
			// c.Read(p)
			// log.Print(p)
			//handle(c)

		}(conn)
	}

}

func handle(c net.Conn) {
	//defer c.Close()

	// data, _ := bufio.NewReader(c).ReadBytes('\n')
	// fmt.Println(data)
	for {
		b, _ := bufio.NewReader(c).ReadString('\n')
		log.Print(b)
		fmt.Print(b)
		if b == "bye" {
			fmt.Println("Okay, bye")
			log.Print("Closed connection.")
		}

	}

}
