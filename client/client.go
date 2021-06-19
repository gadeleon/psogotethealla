package client

// Core client contains all the variables that
// the patch, ship, and login server extend
type CoreClient struct {
	plySockfd int32
	// rcvbuf uint8 // probably don't need in golang, should be def
	rcvread             uint16
	expect              uint16
	snddata, sndwritten int32
	cryptOn             int32
	connected           uint32
	todc                int32
	connectionIndex     uint32
}
