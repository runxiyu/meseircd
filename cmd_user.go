package main

import (
// "log/slog"
)

func init() {
	commandHandlers["USER"] = handleClientUser
}

func handleClientUser(msg RMsg, client *Client) error {
	if len(msg.Params) < 4 {
		return client.Send(MakeMsg(self, ERR_NEEDMOREPARAMS, "USER", "Not enough parameters"))
	}
	switch client.State {
	case ClientStatePreRegistration:
		client.Ident = "~" + msg.Params[0]
		client.Gecos = msg.Params[3]
		err := client.checkRegistration()
		if err != nil {
			return err
		}
	case ClientStateRegistered:
		err := client.Send(MakeMsg(self, ERR_ALREADYREGISTERED, client.Nick, "You may not reregister"))
		if err != nil {
			return err
		}
	case ClientStateRemote:
	}
	return nil
}
