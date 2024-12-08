package main

import (
// "log/slog"
)

func init() {
	commandHandlers["USER"] = handleClientUser
}

func handleClientUser(msg RMsg, client *Client) bool {
	if len(msg.Params) < 4 {
		client.Send(MakeMsg(self, ERR_NEEDMOREPARAMS, "USER", "Not enough parameters"))
		return true
	}
	client.Ident = "~" + msg.Params[0]
	client.Gecos = msg.Params[3]
	if client.State == ClientStatePreRegistration {
		client.checkRegistration()
	}
	return true
}
