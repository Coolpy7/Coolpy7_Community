package packet

import "fmt"

type ConnackCode uint8

const (
	ConnectionAccepted ConnackCode = iota
	InvalidProtocolVersion
	IdentifierRejected
	ServerUnavailable
	BadUsernameOrPassword
	NotAuthorized
)

func (cc ConnackCode) Valid() bool {
	return cc <= 5
}

func (cc ConnackCode) String() string {
	switch cc {
	case ConnectionAccepted:
		return "connection accepted"
	case InvalidProtocolVersion:
		return "connection refused: unacceptable protocol version"
	case IdentifierRejected:
		return "connection refused: identifier rejected"
	case ServerUnavailable:
		return "connection refused: server unavailable"
	case BadUsernameOrPassword:
		return "connection refused: bad user name or password"
	case NotAuthorized:
		return "connection refused: not authorized"
	}

	return "invalid connack code"
}

type Connack struct {
	SessionPresent bool
	ReturnCode     ConnackCode
}

func NewConnack() *Connack {
	return &Connack{}
}

func (c *Connack) Type() Type {
	return CONNACK
}

func (c *Connack) String() string {
	return fmt.Sprintf("<Connack SessionPresent=%t ReturnCode=%d>",
		c.SessionPresent, c.ReturnCode)
}

func (c *Connack) Len() int {
	return headerLen(2) + 2
}

func (c *Connack) Decode(src []byte) (int, error) {
	// decode header
	total, _, _, err := decodeHeader(src, CONNACK)
	if err != nil {
		return total, err
	}

	// read connack flags
	connackFlags, n, err := readUint8(src[total:], CONNACK)
	total += n
	if err != nil {
		return total, err
	}

	// check flags
	if connackFlags&254 != 0 {
		return total, makeError(CONNACK, "bits 7-1 in acknowledge flags are not 0")
	}

	// set session present
	c.SessionPresent = connackFlags&0x1 == 1

	// read return code
	rc, n, err := readUint8(src[total:], CONNACK)
	total += n
	if err != nil {
		return total, err
	}

	// get return code
	returnCode := ConnackCode(rc)
	if !returnCode.Valid() {
		return total, makeError(CONNACK, "invalid return code (%d)", c.ReturnCode)
	}

	// set return code
	c.ReturnCode = returnCode

	return total, nil
}

func (c *Connack) Encode(dst []byte) (int, error) {
	// encode header
	total, err := encodeHeader(dst, 0, 2, c.Len(), CONNACK)
	if err != nil {
		return total, err
	}

	// get connack flags
	var flags uint8
	if c.SessionPresent {
		flags = 0x1 // 00000001
	} else {
		flags = 0x0 // 00000000
	}

	// write flags
	n, err := writeUint8(dst[total:], flags, CONNACK)
	total += n
	if err != nil {
		return total, err
	}

	// check return code
	if !c.ReturnCode.Valid() {
		return total, makeError(CONNACK, "invalid return code (%d)", c.ReturnCode)
	}

	// write return code
	n, err = writeUint8(dst[total:], uint8(c.ReturnCode), CONNACK)
	total += n
	if err != nil {
		return total, err
	}

	return total, nil
}
