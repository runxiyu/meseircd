package main

import (
	"bufio"
	"log"
	"log/slog"
	"net"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	self = Server{
		conn: nil,
		SID:  "001",
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
					slog.Error("connection routine panicked", "raised", raised)
				}
			}()
			defer conn.Close()
			client, err := NewLocalClient(&conn)
			if err != nil {
				slog.Error("cannot make new local client", "error", err)
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
			slog.Error("error while reading from connection", "error", err)
			(*client.conn).Close()
			return
		}
		slog.Debug("recv", "line", line, "conn", client.conn)
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
