package s2s

import (
	"errors"
	"net"
	"strconv"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("s2s")

var bufferSize = 256 * 1024

const (
	IDLE = iota
)

type Client struct {
	address      string
	connection   *net.TCPConn
	peerState    int
	resourceName string
}

func (c *Client) Establish() error {
	if c.peerState != IDLE {
		return errors.New("Peer state must be idle to establish connection")
	}
	conn, err := dialTCP(c.address)
	if err != nil {
		return err
	}
	err = conn.SetNoDelay(true)
	if err != nil {
		return err
	}
	err = conn.SetReadBuffer(bufferSize)
	if err != nil {
		log.Warning("conn.SetReadBuffer(" + strconv.Itoa(bufferSize) + ") failed (" + err.Error() + ")")
	}
	err = conn.SetWriteBuffer(bufferSize)
	if err != nil {
		log.Warning("conn.SetWriteBuffer(" + strconv.Itoa(bufferSize) + ") failed (" + err.Error() + ")")
	}
	_, err = conn.Write([]byte("NiFi"))
	if err != nil {
		return err
	}
	c.connection = conn
	return nil
}

func (c *Client) initiateResourceNegotiation() error {
	_, err := c.connection.Write([]byte(c.resourceName))
	if err != nil {
		return err
	}
	//_, err = c.connection.Write([]byte(5))
	if err != nil {
		return err
	}
	return nil
}

func dialTCP(address string) (*net.TCPConn, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return conn.(*net.TCPConn), err
}
