package main

import (
	"log/slog"
)

func init() {
	commandHandlers["NICK"] = handleClientNick
}

func handleClientNick(msg RMsg, client *Client) bool {
	if len(msg.Params) < 1 {
		client.Send(MakeMsg(self, ERR_NEEDMOREPARAMS, "NICK", "Not enough parameters"))
		return true
	}
	already, exists := nickToClient.LoadOrStore(msg.Params[0], client)
	if exists {
		if already != client {
			client.Send(MakeMsg(self, ERR_NICKNAMEINUSE, client.Nick, msg.Params[0], "Nickname is already in use"))
		}
	} else {
		if client.State == ClientStateRegistered {
			if !nickToClient.CompareAndDelete(client.Nick, client) {
				slog.Error("nick inconsistent", "nick", client.Nick, "client", client)
				return false
			}
			client.Send(MakeMsg(client, "NICK", msg.Params[0]))
		}
		client.Nick = msg.Params[0]
	}
	return true
}
