package connector

import (
	"fmt"
	"os"
)

const (
	COM_QUIT = iota
	COM_REGISTER
)

func (c *Conn) WriteQuit() error {
	data := make([]byte, 3, 4)

	data = append(data, COM_QUIT)

	return c.WritePacket(data)
}

func (c *Conn) WriteRegister() error {
	host, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("get hostname error %s", err.Error())
	}

	data := make([]byte, 3, 4+len(host))

	data = append(data, COM_REGISTER)
	data = append(data, []byte(host)...)

	return c.WritePacket(data)
}
