package main

import (
	"net"

	"git.sr.ht/~runxiyu/meseircd/meselog"
)

type Server struct {
	conn *net.Conn
	SID  uint32
	Name string
}

func (server *Server) Send(msg SMsg) error {
	return server.SendRaw(msg.ServerSerialize())
}

func (server *Server) SendRaw(s string) error {
	if server == self {
		return ErrSendToSelf
	}
	if server.conn == nil {
		// TODO: Propagate across mesh
		return ErrNotConnectedServer
	}
	meselog.Debug("send", "line", s, "conn", server.conn)
	_, err := (*server.conn).Write([]byte(s))
	if err != nil {
		// TODO: Should shut down the netFd instead but the stdlib
		// doesn't expose a way to do this.
		(*server.conn).Close()
		return err
	}
	return nil
}

func (server *Server) ClientSource() string {
	return server.Name
}

func (server *Server) ServerSource() uint64 {
	return uint64(server.SID) << 32
}

var self *Server
