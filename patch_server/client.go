package patchserver

import (
	"net"
)

type Client struct {
	Socket       net.Conn
	Data         chan []byte
	Patched      bool
	Encrypted    bool
	ServerCipher CryptSetup
	ClientCipher CryptSetup
}
