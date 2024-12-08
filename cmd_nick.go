package main

func init() {
	commandHandlers["NICK"] = handleClientNick
}

func handleClientNick(msg RMsg, client *Client) bool {
	if len(msg.Params) < 1 {
		client.Send(MakeMsg(self, ERR_NEEDMOREPARAMS, "NICK", "Not enough parameters"))
		return true
	}
	_, exists := nickToClient.LoadOrStore(msg.Params[0], client)
	if exists {
		client.Send(MakeMsg(self, ERR_NICKNAMEINUSE, client.Nick, msg.Params[0], "Nickname is already in use"))
	} else {
		client.Nick = msg.Params[0]
		if client.State == ClientStateRegistered {
			client.Send(MakeMsg(client, "NICK", msg.Params[0]))
		}
	}
	return true
}
