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
	switch {
	case client.State < ClientStateRegistered:
		client.Ident = "~" + msg.Params[0]
		client.Gecos = msg.Params[3]
		err := client.checkRegistration()
		if err != nil {
			return err
		}
	case client.State == ClientStateRegistered:
		err := client.Send(MakeMsg(self, ERR_ALREADYREGISTERED, client.Nick, "You may not reregister"))
		if err != nil {
			return err
		}
	case client.State == ClientStateRemote:
	}
	return nil
}
