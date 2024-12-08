package main

import (
	"fmt"
	"log/slog"
)

func init() {
	CommandHandlers["NICK"] = handleClientNick
}

func handleClientNick(msg RMsg, client *Client) error {
	if len(msg.Params) < 1 {
		return client.Send(MakeMsg(self, ERR_NEEDMOREPARAMS, "NICK", "Not enough parameters"))
	}
	already, exists := nickToClient.LoadOrStore(msg.Params[0], client)
	if exists {
		if already != client {
			err := client.Send(MakeMsg(self, ERR_NICKNAMEINUSE, client.Nick, msg.Params[0], "Nickname is already in use"))
			if err != nil {
				return err
			}
		}
	} else {
		if (client.State >= ClientStateRegistered || client.Nick != "*") && !nickToClient.CompareAndDelete(client.Nick, client) {
			slog.Error("nick inconsistent", "nick", client.Nick, "client", client)
			return fmt.Errorf("%w: %v", ErrInconsistentClient, client)
		}
		if client.State == ClientStateRegistered {
			err := client.Send(MakeMsg(client, "NICK", msg.Params[0]))
			if err != nil {
				return err
			}
		}
		client.Nick = msg.Params[0]
	}
	if client.State < ClientStateRegistered {
		err := client.checkRegistration()
		if err != nil {
			return err
		}
	}
	return nil
}
