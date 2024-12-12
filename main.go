package main

import (
	"bufio"
	"log"
	"net"

	"git.sr.ht/~runxiyu/meseircd/meselog"
)

const VERSION = "MeseIRCd-0.0.0"

func main() {
	setupCapls()

	self = &Server{
		conn: nil,
		SID:  0,
		Name: "irc.runxiyu.org",
	}

	listener, err := net.Listen("tcp", ":6667")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer func() {
				raised := recover()
				if raised != nil {
					meselog.Error("connection routine panicked", "raised", raised)
				}
			}()
			defer conn.Close()
			client, err := NewLocalClient(&conn)
			if err != nil {
				meselog.Error("cannot make new local client", "error", err)
			}
			defer client.Teardown()
			client.handleConnection()
		}()
	}
}

func (client *Client) handleConnection() {
	reader := bufio.NewReader(*client.conn)
messageLoop:
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			meselog.Error("error while reading from connection", "error", err)
			(*client.conn).Close()
			return
		}
		meselog.Debug("recv", "line", line, "client", client.CID)
		msg, err := parseIRCMsg(line)
		if err != nil {
			switch err {
			case ErrEmptyMessage:
				continue messageLoop
			case ErrIllegalByte:
				err := client.Send(MakeMsg(self, "ERROR", err.Error()))
				if err != nil {
					meselog.Error("error while reporting illegal byte", "error", err, "client", client)
					return
				}
				return
			case ErrTagsTooLong:
				fallthrough
			case ErrBodyTooLong:
				err := client.Send(MakeMsg(self, ERR_INPUTTOOLONG, err.Error()))
				if err != nil {
					meselog.Error("error while reporting body too long", "error", err, "client", client)
					return
				}
				continue messageLoop
			default:
				err := client.Send(MakeMsg(self, "ERROR", err.Error()))
				if err != nil {
					meselog.Error("error while reporting parser error", "error", err, "client", client)
				}
				return
			}
		}

		handler, ok := CommandHandlers[msg.Command]
		if !ok {
			err := client.Send(MakeMsg(self, ERR_UNKNOWNCOMMAND, msg.Command, "Unknown command"))
			if err != nil {
				meselog.Error("error while reporting unknown command", "error", err, "client", client)
				return
			}
			continue
		}

		err = handler(msg, client)
		if err != nil {
			meselog.Error("handler error", "error", err, "client", client)
			return
		}
	}
}
