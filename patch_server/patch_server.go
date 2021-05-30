// ************************************************************************
//   Tethealla Patch Server
//   Copyright (C) 2008  Terry Chatman Jr.

//   This program is free software: you can redistribute it and/or modify
//   it under the terms of the GNU General Public License version 3 as
//   published by the Free Software Foundation.

//   This program is distributed in the hope that it will be useful,
//   but WITHOUT ANY WARRANTY; without even the implied warranty of
//   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//   GNU General Public License for more details.

//   You should have received a copy of the GNU General Public License
//   along with this program.  If not, see <http://www.gnu.org/licenses/>.
// ************************************************************************

package patchserver

import (
	"errors"
	"log"
	"net"
)

type PatchServer interface {
	Serve()
}

const (
	MAX_PATCHES                    = 4096
	PATCH_COMPILED_MAX_CONNECTIONS = 300
	SERVER_VERSION                 = "0.0.1"
	SEND_PACKET                    = 0x00
	RECEIVE_PACKET_02              = 0x01
	RECEIVE_PACKET_10              = 0x02
	SEND_PACKET_0B                 = 0x03
	MAX_SENDCHECK                  = 0x04
	TCP_BUFFER_SIZE                = 65530
	SOCKET_ERROR                   = -1
	MAX_SIMULTANEOUS_CONNECTIONS   = 6
)

// "Encyption Data Struct"
// Scaffolding for communication with PSO client...?
type CryptSetup struct {
	keys    [1024]uint32 //encryption stream
	pc_posn uint32       //
}

var Crypt_PC_GetNextKey func(c *CryptSetup) uint32

func send_to_server(sock int, packet []byte) error {
	log.Print("Sending to patch server...")
	//pktlen := len(packet)
	// C code sends the message and compares size of response
	// TODO: Suss out sending messages/socks in GO properly
	err := errors.New("send_to_server(): failure")
	log.Print(err)
	return err
}

// convert string of IP to ...something
// In the C code, serverIP is an array of 4 integers ([4]int in Go)
// I don't know if I need to recreate that, but we'll see.

// Parses IP from config file, if it can't parse
// then it grabs IPv4 from net.LookupIP
func parseIPString(ip string) (net.IP, error) {
	log.Printf("Using '%s' for IP/Host\n", ip)
	if ip == "" {
		return nil, errors.New("IP/Host provided is blank")
	}
	// Parse IP Address
	// If addr is nil, lookup the IPv4
	addr := net.ParseIP(ip)
	if addr == nil {
		log.Println("Provided string is likely a hostname, looking up...")
		lookup, err := net.LookupIP(ip)
		log.Println("IP lookup produced these IPS:", lookup)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for _, ipv4 := range lookup {
			if ipv4.To4() != nil {
				addr = ipv4
			}
		}
	}
	return addr, nil
}

// Simple func to convert int of bytes into KB
func bytesToKB(bytes int) int {
	return bytes * 1024
}
