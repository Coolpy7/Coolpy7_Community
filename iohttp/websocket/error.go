package websocket

import (
	"errors"
	"fmt"
)

var (
	// ErrUpgradeTokenNotFound .
	ErrUpgradeTokenNotFound = errors.New("websocket: the client is not using the websocket protocol: 'upgrade' token not found in 'Connection' header")

	// ErrUpgradeMethodIsGet .
	ErrUpgradeMethodIsGet = errors.New("websocket: the client is not using the websocket protocol: request method is not GET")

	// ErrUpgradeInvalidWebsocketVersion .
	ErrUpgradeInvalidWebsocketVersion = errors.New("websocket: unsupported version: 13 not found in 'Sec-Websocket-Version' header")

	// ErrUpgradeUnsupportedExtensions .
	ErrUpgradeUnsupportedExtensions = errors.New("websocket: application specific 'Sec-WebSocket-Extensions' headers are unsupported")

	// ErrUpgradeOriginNotAllowed .
	ErrUpgradeOriginNotAllowed = errors.New("websocket: request origin not allowed by Upgrader.CheckOrigin")

	// ErrUpgradeMissingWebsocketKey .
	ErrUpgradeMissingWebsocketKey = errors.New("websocket: not a websocket handshake: 'Sec-WebSocket-Key' header is missing or blank")

	// ErrUpgradeNotHijacker .
	ErrUpgradeNotHijacker = errors.New("websocket: response does not implement http.Hijacker")

	// ErrInvalidControlFrame .
	ErrInvalidControlFrame = errors.New("websocket: invalid control frame")

	// ErrInvalidWriteCalling .
	ErrInvalidWriteCalling = errors.New("websocket: invalid write calling, should call WriteMessage instead")

	// ErrReserveBitSet .
	ErrReserveBitSet = errors.New("websocket: reserved bit set it frame")

	// ErrReservedOpcodeSet .
	ErrReservedOpcodeSet = errors.New("websocket: reserved opcode received")

	// ErrControlMessageFragmented .
	ErrControlMessageFragmented = errors.New("websocket: control messages must not be fragmented")

	// ErrFragmentsShouldNotHaveBinaryOrTextOpcode .
	ErrFragmentsShouldNotHaveBinaryOrTextOpcode = errors.New("websocket: fragments should not have opcode of text or binary")

	// ErrInvalidCloseCode .
	ErrInvalidCloseCode = errors.New("websocket: invalid close code")

	// ErrBadHandshake .
	ErrBadHandshake = errors.New("websocket: bad handshake")

	// ErrInvalidCompression .
	ErrInvalidCompression = errors.New("websocket: invalid compression negotiation")

	// ErrMalformedURL .
	ErrMalformedURL = errors.New("malformed ws or wss URL")

	// ErrMessageTooLarge.
	ErrMessageTooLarge = errors.New("message exceeds the configured limit")

	// ErrMessageSendQuqueIsFull .
	ErrMessageSendQuqueIsFull = errors.New("message send queue is full")
)

// CloseError .
type CloseError struct {
	Code   int
	Reason string
}

// Error .
func (ce CloseError) Error() string {
	return fmt.Sprintf("websocket: close code=%d and reason=%q", ce.Code, ce.Reason)
}

// CloseCode .
func CloseCode(err error) int {
	var ce CloseError
	if errors.As(err, &ce) {
		return ce.Code
	}
	return -1
}

// CloseReason .
func CloseReason(err error) string {
	var ce CloseError
	if errors.As(err, &ce) {
		return ce.Reason
	}
	return ""
}
