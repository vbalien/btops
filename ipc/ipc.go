package ipc

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net"
)

const (
	defaultBspwmSocket = "/tmp/bspwmunix_0_0-socket"
)

type Subscriber struct {
	Scanner *bufio.Scanner
	conn    *net.UnixConn
}

func (s *Subscriber) Close() error {
	return s.conn.Close()
}

func NewSubscriber() (*Subscriber, error) {
	conn, err := newConn()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(conn)
	_, err = conn.Write(buildPayload("subscribe", "report"))
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		Scanner: scanner,
		conn:    conn,
	}, nil
}

func newConn() (*net.UnixConn, error) {
	raddr, err := net.ResolveUnixAddr("unix", defaultBspwmSocket)
	if err != nil {
		return nil, err
	}
	return net.DialUnix(raddr.Network(), nil, raddr)
}

func buildPayload(cmd ...string) []byte {
	var buffer bytes.Buffer

	for i := range cmd {
		buffer.WriteString(cmd[i])
		buffer.WriteByte('\x00')
	}

	return buffer.Bytes()
}

func Send(cmd ...string) (response []byte, err error) {
	conn, err := newConn()
	if err != nil {
		return response, err
	}
	defer conn.Close()

	_, err = conn.Write(buildPayload(cmd...))
	if err != nil {
		return response, err
	}

	return ioutil.ReadAll(conn)
}
