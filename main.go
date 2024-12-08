package main

import (
	"bufio"
	"log"
	"log/slog"
	"net"
)

func main() {
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
			conn: conn,
		}
		go func() {
			defer func() {
				client.conn.Close()
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
	reader := bufio.NewReader(client.conn)
	messageLoop: for {
		line, err := reader.ReadString('\n')
		if err != nil {
			slog.Error("error while reading from connection", "error", err)
			client.conn.Close()
			return
		}
		msg, err := parseIRCMsg(line)
		if err != nil {
			switch (err) {
			case ErrEmptyMessage:
				continue messageLoop
			case ErrIllegalByte:
				client.Send(SMsg{Command: "ERROR", Params: []string{err.Error()}})
				break messageLoop
			case ErrTagsTooLong:
				fallthrough
			case ErrBodyTooLong:
				client.Send(SMsg{Command: ERR_INPUTTOOLONG, Params: []string{err.Error()}})
				continue messageLoop
			default:
				client.Send(SMsg{Command: "ERROR", Params: []string{err.Error()}})
				break messageLoop
			}
		}

		handler, ok := commandHandlers[msg.Command]
		if !ok {
			client.Send(SMsg{Command: ERR_UNKNOWNCOMMAND, Params: []string{msg.Command, "Unknown command"}})
			continue
		}

		err = handler(msg, client)
		if err != nil {
			client.Send(SMsg{Command: "ERROR", Params: []string{err.Error()}})
			break
		}
	}
}
