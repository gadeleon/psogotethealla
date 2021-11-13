package client

// Core client contains all the variables that
// the patch, ship, and login server extend
type CoreClient struct {
	PlySockfd int32
	// rcvbuf uint8 // probably don't need in golang, should be def
	Rcvread             uint16
	Expect              uint16
	Snddata, Sndwritten int32
	CryptOn             int32
	Connected           uint32
	Todc                int32
	ConnectionIndex     uint32
}
