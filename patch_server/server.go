package patchserver

import (
	"log"
	"net"
)

type PatchServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func (p *PatchServer) listen() {
	for {
		select {
		case connection := <-p.register:
			p.clients[connection] = true
		case connection := <-p.unregister:
			if _, ok := p.clients[connection]; ok {
				close(connection.Data)
				delete(p.clients, connection)
				log.Printf("Connection from %s terminated", connection.Socket.RemoteAddr().String())
			}
		case message := <-p.broadcast:
			for connection := range p.clients {
				select {
				case connection.Data <- message:
				default:
					close(connection.Data)
					delete(p.clients, connection)
				}
			}
		}
	}
}

func (p *PatchServer) receive(client *Client) {
	for {
		message := make([]byte, 4096)
		length, err := client.Socket.Read(message)
		if err != nil {
			p.unregister <- client
			client.Socket.Close()
			break
		}
		//p.welcome(client)
		if !client.Encrypted {
			log.Print("Starting encryption.")
			SetUpKeys(client)
			log.Printf("Before cipher: %s", string(message))
			Crypt_PC_CryptData(&client.ClientCipher, message)
			//log.Printf("After cipher: %s", string(message))
			log.Print("Did we decrypt it?")

		}

		if length > 0 {
			log.Printf("Received: %s", string(message))
			//log.Print(message)
			//p.broadcast <- message
		}
	}
}

func (p *PatchServer) send(client *Client) {
	defer client.Socket.Close()
	for {
		select {
		case message, ok := <-client.Data:
			if !ok {
				return
			}
			client.Socket.Write(message)

		}
	}
}

func (p *PatchServer) welcome(client *Client) {
	welpacket := GetPacket02()
	// for c := 0; c < 8; c++ {
	// 	welpacket[0x44+c] = uint8(rand.Int() % 255)
	// }
	client.Socket.Write(welpacket)
}

func StartServer() {
	log.Print("Setting up Listener")
	listener, err := net.Listen("tcp", "10.211.55.2:11000")
	if err != nil {
		log.Fatal(err)
	}
	server := PatchServer{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	log.Print("Server made, listening.")
	go server.listen()

	for {
		connection, _ := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		client := &Client{
			Socket: connection,
			Data:   make(chan []byte),
		}
		log.Printf("Accepted connection from: %s", client.Socket.RemoteAddr().String())
		server.register <- client
		log.Print("Clients,", server.clients)
		if !client.Patched {
			server.welcome(client)
		}

		go server.receive(client)
		go server.send(client)
	}
}
