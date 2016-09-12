package connector

import (
	"errors"
	"fmt"
	"io"
	"net"
)

const (
	MaxPayloadLen int = 1<<24 - 1
)

var (
	ErrBadConn = errors.New("connection was bad")
)

type Conn struct {
	net.Conn
}

func NewConn(conn net.Conn) *Conn {
	return &Conn{
		Conn: conn,
	}
}

func (c *Conn) ReadPacket() ([]byte, error) {
	header := []byte{0, 0, 0}

	if _, err := io.ReadFull(c.Conn, header); err != nil {
		return nil, ErrBadConn
	}

	length := int(uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16)
	if length < 1 {
		return nil, fmt.Errorf("invalid payload length %d", length)
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(c.Conn, data); err != nil {
		return nil, ErrBadConn
	} else {
		if length < MaxPayloadLen {
			return data, nil
		}

		var buf []byte
		buf, err = c.ReadPacket()
		if err != nil {
			return nil, ErrBadConn
		} else {
			return append(data, buf...), nil
		}
	}
}

func (c *Conn) WritePacket(data []byte) error {
	length := len(data) - 3

	for length >= MaxPayloadLen {
		data[0] = 0xff
		data[1] = 0xff
		data[2] = 0xff

		if n, err := c.Write(data[:3+MaxPayloadLen]); err != nil {
			return ErrBadConn
		} else if n != (3 + MaxPayloadLen) {
			return ErrBadConn
		} else {
			length -= MaxPayloadLen
			data = data[MaxPayloadLen:]
		}
	}

	data[0] = byte(length)
	data[1] = byte(length >> 8)
	data[2] = byte(length >> 16)

	if n, err := c.Write(data); err != nil {
		return ErrBadConn
	} else if n != len(data) {
		return ErrBadConn
	} else {
		return nil
	}
}
