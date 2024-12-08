package main

import (
	"strconv"
	"strings"
)

func init() {
	CommandHandlers["CAP"] = handleClientCap
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
		var err error
		if len(msg.Params) >= 2 {
			capVersion, err := strconv.ParseUint(msg.Params[1], 10, 64)
			if err == nil && capVersion >= 302 {
				err = client.Send(MakeMsg(self, "CAP", client.Nick, "LS", capls302))
			} else {
				err = client.Send(MakeMsg(self, "CAP", client.Nick, "LS", capls))
			}
		} else {
			err = client.Send(MakeMsg(self, "CAP", client.Nick, "LS", capls))
		}
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
				err := client.Send(MakeMsg(self, "CAP", client.Nick, "ACK", c))
				if err != nil {
					return err
				}
				client.Caps[c] = struct{}{}
				// TODO: This is terrible
			} else {
				err := client.Send(MakeMsg(self, "CAP", client.Nick, "NAK", c))
				if err != nil {
					return err
				}
			}
		}
	case "END":
		if client.State != ClientStateCapabilities {
			// Just ignore it
			return nil
		}
		client.State = ClientStateCapabilitiesFinished
		err := client.checkRegistration()
		if err != nil {
			return err
		}
	}
	return nil
}
