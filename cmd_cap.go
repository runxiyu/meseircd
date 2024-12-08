package main

import (
	"strings"
)

func init() {
	commandHandlers["CAP"] = handleClientCap
}

func handleClientCap(msg RMsg, client *Client) error {
	if len(msg.Params) < 1 {
		err := client.Send(MakeMsg(self, ERR_NEEDMOREPARAMS, "CAP", "Not enough parameters"))
		if err != nil {
			return err
		}
		return nil
	}
	if client.State == ClientStateRemote {
		return ErrRemoteClient
	}
	switch strings.ToUpper(msg.Params[0]) {
	case "LS":
		if client.State == ClientStatePreRegistration {
			client.State = ClientStateCapabilities
		}
		err := client.Send(MakeMsg(self, "CAP", client.Nick, "LS", capls))
		// TODO: Split when too long
		if err != nil {
			return err
		}
	case "REQ":
		if client.State == ClientStatePreRegistration {
			client.State = ClientStateCapabilities
		}
		caps := strings.Split(msg.Params[1], " ")
		for _, c := range caps {
			if c[0] == '-' {
				// TODO: Remove capability
				delete(client.Caps, c)
				continue
			}
			_, ok := Caps[c]
			if ok {
				client.Send(MakeMsg(self, "CAP", client.Nick, "ACK", c))
				client.Caps[c] = struct{}{}
				// TODO: This is terrible
			} else {
				client.Send(MakeMsg(self, "CAP", client.Nick, "NAK", c))
			}
		}
	case "END":
		if client.State != ClientStateCapabilities {
			// Just ignore it
			return nil
		}
		client.State = ClientStateCapabilitiesFinished
		client.checkRegistration()
	}
	return nil
}
