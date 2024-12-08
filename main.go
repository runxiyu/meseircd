package main

import (
	"bufio"
	"log"
	"log/slog"
	"net"
)

func main() {
	self = Server{
		conn: nil,
		SID:  [3]byte{'1', 'H', 'C'},
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

		client := &Client{
			conn:   &conn,
			Server: self,
		}
		go func() {
			defer func() {
				(*client.conn).Close()
				// TODO: Unified client clean-up
			}()
			defer func() {
				raised := recover()
				if raised != nil {
					slog.Error("connection routine panicked", "raised", raised)
				}
			}()
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
			slog.Error("error while reading from connection", "error", err)
			(*client.conn).Close()
			return
		}
		msg, err := parseIRCMsg(line)
		if err != nil {
			switch err {
			case ErrEmptyMessage:
				continue messageLoop
			case ErrIllegalByte:
				client.Send(MakeMsg(self, "ERROR", err.Error()))
				break messageLoop
			case ErrTagsTooLong:
				fallthrough
			case ErrBodyTooLong:
				client.Send(MakeMsg(self, ERR_INPUTTOOLONG, err.Error()))
				continue messageLoop
			default:
				client.Send(MakeMsg(self, "ERROR", err.Error()))
				break messageLoop
			}
		}

		handler, ok := commandHandlers[msg.Command]
		if !ok {
			client.Send(MakeMsg(self, ERR_UNKNOWNCOMMAND, msg.Command, "Unknown command"))
			continue
		}

		cont := handler(msg, client)
		if !cont {
			break
		}
	}
}
