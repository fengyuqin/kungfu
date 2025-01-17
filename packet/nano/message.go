/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package nano

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/fengyuqin/kungfu/v2/logger"
	"strings"
)

// MsgType represents the type of message, which could be Request/Notify/Response/Push
type MsgType byte

// Message types
const (
	Request  MsgType = 0x00
	Notify           = 0x01
	Response         = 0x02
	Push             = 0x03
)

const (
	msgRouteCompressMask = 0x01
	msgTypeMask          = 0x07
	msgRouteLengthMask   = 0xFF
	msgHeadLength        = 0x02
)

var types = map[MsgType]string{
	Request:  "Request",
	Notify:   "Notify",
	Response: "Response",
	Push:     "Push",
}

var (
	routes = make(map[string]uint16) // route map to code
	codes  = make(map[uint16]string) // code map to route
)

// Errors that could be occurred in message codec
var (
	ErrWrongMessageType  = errors.New("wrong message type")
	ErrInvalidMessage    = errors.New("invalid message")
	ErrRouteInfoNotFound = errors.New("route info not found in dictionary")
)

// Message represents a unmarshaled message or a message which to be marshaled
type Message struct {
	Type       MsgType // message type
	ID         uint    // unique id, zero while notify mode
	Route      string  // route for locating service
	Data       []byte  // payload
	compressed bool    // is message compressed
}

// NewMessage returns a new message instance
func NewMessage() *Message {
	return &Message{}
}

// String, implementation of fmt.Stringer interface
func (m *Message) String() string {
	return fmt.Sprintf("MsgType: %s, ID: %d, Route: %s, Compressed: %t, BodyLength: %d",
		types[m.Type],
		m.ID,
		m.Route,
		m.compressed,
		len(m.Data))
}

// Encode marshals message to binary format.
func (m *Message) Encode() ([]byte, error) {
	return MsgEncode(m)
}

func routable(t MsgType) bool {
	return t == Request || t == Notify || t == Push
}

func invalidType(t MsgType) bool {
	return t < Request || t > Push

}

// MsgEncode marshals message to binary format. Different message types is corresponding to
// different message header, message types is identified by 2-4 bit of flag field. The
// relationship between message types and message header is presented as follows:
// ------------------------------------------
// |   type   |  flag  |       other        |
// |----------|--------|--------------------|
// | request  |----000-|<message id>|<route>|
// | notify   |----001-|<route>             |
// | response |----010-|<message id>        |
// | push     |----011-|<route>             |
// ------------------------------------------
// The figure above indicates that the bit does not affect the type of message.
// See ref: https://github.com/lonnng/nano/blob/master/docs/communication_protocol.md
func MsgEncode(m *Message) ([]byte, error) {
	if invalidType(m.Type) {
		return nil, ErrWrongMessageType
	}

	buf := make([]byte, 0)
	flag := byte(m.Type) << 1

	code, compressed := routes[m.Route]
	if compressed {
		flag |= msgRouteCompressMask
	}
	buf = append(buf, flag)

	if m.Type == Request || m.Type == Response {
		n := m.ID
		// variant length encode
		for {
			b := byte(n % 128)
			n >>= 7
			if n != 0 {
				buf = append(buf, b+128)
			} else {
				buf = append(buf, b)
				break
			}
		}
	}

	if routable(m.Type) {
		if compressed {
			buf = append(buf, byte((code>>8)&0xFF))
			buf = append(buf, byte(code&0xFF))
		} else {
			buf = append(buf, byte(len(m.Route)))
			buf = append(buf, []byte(m.Route)...)
		}
	}

	buf = append(buf, m.Data...)
	return buf, nil
}

// MsgDecode unmarshal the bytes slice to a message
// See ref: https://github.com/lonnng/nano/blob/master/docs/communication_protocol.md
func MsgDecode(data []byte) (*Message, error) {
	if len(data) < msgHeadLength {
		return nil, ErrInvalidMessage
	}
	m := NewMessage()
	flag := data[0]
	offset := 1
	m.Type = MsgType((flag >> 1) & msgTypeMask)

	if invalidType(m.Type) {
		return nil, ErrWrongMessageType
	}

	if m.Type == Request || m.Type == Response {
		id := uint(0)
		// little end byte order
		// WARN: must can be stored in 64 bits integer
		// variant length encode
		for i := offset; i < len(data); i++ {
			b := data[i]
			id += uint(b&0x7F) << uint(7*(i-offset))
			if b < 128 {
				offset = i + 1
				break
			}
		}
		m.ID = id
	}

	if routable(m.Type) {
		if flag&msgRouteCompressMask == 1 {
			m.compressed = true
			code := binary.BigEndian.Uint16(data[offset:(offset + 2)])
			route, ok := codes[code]
			if !ok {
				return nil, ErrRouteInfoNotFound
			}
			m.Route = route
			offset += 2
		} else {
			m.compressed = false
			rl := data[offset]
			offset++
			m.Route = string(data[offset:(offset + int(rl))])
			offset += int(rl)
		}
	}

	m.Data = data[offset:]
	return m, nil
}

// SetDictionary set routes map which be used to compress route.
// TODO(warning): set dictionary in runtime would be a dangerous operation!!!!!!
func SetDictionary(dict map[string]uint16) {
	for route, code := range dict {
		r := strings.TrimSpace(route)

		// duplication check
		if _, ok := routes[r]; ok {
			logger.Infof("duplicated route(route: %s, code: %d)\n", r, code)
		}

		if _, ok := codes[code]; ok {
			logger.Infof("duplicated route(route: %s, code: %d)\n", r, code)
		}

		// update map, using last value when key duplicated
		routes[r] = code
		codes[code] = r
	}
}
