package main

import (
	"net"
)

type Server struct {
	conn *net.Conn
	SID  [3]byte
	Name string
}

func (server *Server) Send(msg SMsg) error {
	return server.SendRaw(msg.ServerSerialize())
}

func (server *Server) SendRaw(s string) error {
	if server == &self {
		return ErrSendToSelf
	}
	if server.conn == nil {
		// TODO: Propagate across mesh
		return ErrNotConnectedServer
	}
	_, err := (*server.conn).Write([]byte(s))
	if err != nil {
		// TODO: Should shut down the netFd instead but the stdlib
		// doesn't expose a way to do this.
		(*server.conn).Close()
		return err
	}
	return nil
}

func (server Server) ClientSource() string {
	return server.Name
}

func (server Server) ServerSource() string {
	return string(server.SID[:])
}

var self Server