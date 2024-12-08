package main

func init() {
	commandHandlers["PING"] = handleClientPing
}

func handleClientPing(msg RMsg, client *Client) error {
	if len(msg.Params) < 1 {
		return client.Send(MakeMsg(self, ERR_NEEDMOREPARAMS, "PING", "Not enough parameters"))
	}
	return client.Send(MakeMsg(self, "PONG", self.Name, msg.Params[0]))
}
