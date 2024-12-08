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
	switch client.State {
	case ClientStatePreRegistration:
		client.Ident = "~" + msg.Params[0]
		client.Gecos = msg.Params[3]
		client.checkRegistration()
	case ClientStateRegistered:
		client.Send(MakeMsg(self, ERR_ALREADYREGISTERED, client.Nick, "You may not reregister"))
	case ClientStateRemote:
	}
	return true
}
