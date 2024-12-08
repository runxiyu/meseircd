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
	switch strings.ToUpper(msg.Params[0]) {
	case "LS":
		err := client.Send(MakeMsg(self, "CAP", client.Nick, "LS", "sasl=PLAIN,EXTERNAL"))
		if err != nil {
			return err
		}
	case "REQ":
	}
	return nil
}
