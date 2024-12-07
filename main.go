package main

import (
	"bufio"
	"log"
	"log/slog"
	"net"
)

type Client struct {
	conn net.Conn
}

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
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			slog.Error("error while reading from connection", "error", err)
			client.conn.Close()
			return
		}
		_, _ = parseIRCMsg(line)
	}
}
