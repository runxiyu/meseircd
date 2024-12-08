package main

import (
	"strings"
)

func init() {
	commandHandlers["CAP"] = handleClientCap
}

func handleClientCap(msg RMsg, client *Client) bool {
	if len(msg.Params) < 1 {
		client.Send(MakeMsg(self, ERR_NEEDMOREPARAMS, "CAP", "Not enough parameters"))
		return true
	}
	switch (strings.ToUpper(msg.Params[0])) {
	case "LS":
		client.Send(MakeMsg(self, "CAP", client.Nick, "LS", "sasl=PLAIN,EXTERNAL"))
	case "REQ":
	}
	return true
}
