package patchserver

import (
	"encoding/binary"
	"log"
	"math/rand"
	"time"
)

func SetUpKeys(pc *Client) {
	log.Printf("Setting up keypair for %s", pc.Socket.RemoteAddr().String())
	// Modify last bytes of welcome packet for cipher
	// Not sure if actually a nonce but it's not a crib
	var nonce uint32
	packet02 := GetPacket02()
	for c := 0; c < 8; c++ {
		packet02[0x44+c] = uint8(rand.Int() % 255)
	}
	// Convert data we want into uint32 for crypto methods
	crypto := packet02[0x44 : 0x44+4]
	nonce = binary.BigEndian.Uint32(crypto)
	Create_PC_Keys(&pc.ServerCipher, nonce)

	crypto = packet02[0x44+4 : 0x44+8]
	nonce = binary.BigEndian.Uint32(crypto)
	Create_PC_Keys(&pc.ClientCipher, nonce)
	log.Print("Key pair created...?")

}

func StartEncryption(pc *PatchClient) {
	var um uint32
	log.Printf("Encrypting! %s", pc.IpAddress)
	packet02 := GetPacket02()
	copy(pc.sndbuf[:], packet02)
	// Add some flavor
	for c := 0; c < 8; c++ {
		pc.sndbuf[0x44+c] = uint8(rand.Int() % 255)
	}
	// Add Size Packet02 which is 76
	pc.Core.Snddata += 76
	// Convert data we want into uint32 for crypto
	crypto := []byte(pc.sndbuf[0x44 : 0x44+4])
	um = binary.BigEndian.Uint32(crypto)
	Create_PC_Keys(&pc.serverCipher, um)

	crypto = []byte(pc.sndbuf[0x44+4 : 0x44+8])
	um = binary.BigEndian.Uint32(crypto)
	Create_PC_Keys(&pc.clientCipher, um)
	pc.Core.CryptOn = 1
	pc.sendCheck[SEND_PACKET_02] = 1
	pc.Core.Connected = uint32(time.Now().Unix()) // this may be 64 bit.
	log.Print("Encryption Complete?")
	log.Print(pc.clientCipher)

}

// OG C
//void start_encryption(BANANA* connect)
// memcpy (&connect->sndbuf[0], &Packet02[0], sizeof (Packet02));
// for (c=0;c<8;c++)
//   connect->sndbuf[0x44+c] = (uint8_t) rand() % 255;
// connect->snddata += sizeof (Packet02);

// memcpy (&c, &connect->sndbuf[0x44], 4);
// CRYPT_PC_CreateKeys(&connect->server_cipher,c);
// memcpy (&c, &connect->sndbuf[0x48], 4);
// CRYPT_PC_CreateKeys(&connect->client_cipher,c);
// connect->crypt_on = 1;
// connect->sendCheck[SEND_PACKET_02] = 1;
// connect->connected = (unsigned) servertime;

func Create_PC_Keys(c *CryptSetup, val uint32) {
	// Setup...registers?
	var esi, ebx, edi, eax, edx, var1 uint32
	esi = 1
	ebx = val
	edi = 0x15 // aka 15 in decimal
	c.keys[56] = ebx
	c.keys[55] = ebx
	for edi <= 0x46E { // huh
		eax = edi
		var1 = eax / 55
		edx = eax - (var1 * 55)
		ebx = ebx - esi
		edi = edi + 0x15
		c.keys[edx] = esi
		esi = ebx
		ebx = c.keys[edx]
	}
	Crypt_PC_MixKeys(c)
	Crypt_PC_MixKeys(c)
	Crypt_PC_MixKeys(c)
	Crypt_PC_MixKeys(c)
	c.pcPosn = 56

}

// mix...the keys?
func Crypt_PC_MixKeys(c *CryptSetup) {
	var esi, edi, eax, ebp, edx uint32
	edi = 1
	edx = 0x18
	eax = edi
	for edx > 0 {
		esi = c.keys[eax+0x1F]
		ebp = c.keys[eax]
		ebp = ebp - esi
		c.keys[eax] = ebp
		eax++
		edx--
	}
	edi = 0x19
	edx = 0x1F
	eax = edi
	for edx > 0 {
		esi = c.keys[eax-0x18]
		ebp = c.keys[eax]
		ebp = ebp - esi
		c.keys[eax] = ebp
		eax++
		edx--
	}

}

func Crypt_PC_GetNextKey(c *CryptSetup) uint32 {
	var re uint32
	if c.pcPosn == 56 {
		Crypt_PC_MixKeys(c)
		c.pcPosn = 1
	}
	re = c.keys[c.pcPosn]
	c.pcPosn++
	return re
}

// What's the type in the message???
func Crypt_PC_CryptData(c *CryptSetup, message []byte) {
	// Here is where 64 v. 32 bit comes into play.
	// This will break at some point.

	var x uint32
	// XOR the message
	cnt := 0
	out := make([]byte, len(message))

	for cnt < len(message) {
		um := binary.BigEndian.Uint32(message[cnt : cnt+4])
		data := um + x
		data ^= Crypt_PC_GetNextKey(c)
		x += 4
		cnt += 4
		binary.LittleEndian.PutUint32(out, data)
	}
	log.Printf("After cipher: %s", string(out))

}
