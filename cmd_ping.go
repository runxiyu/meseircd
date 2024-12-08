package main

func init() {
	commandHandlers["PING"] = handleClientPing
}

func handleClientPing(msg RMsg, client *Client) bool {
	if len(msg.Params) < 1 {
		client.Send(MakeMsg(self, ERR_NEEDMOREPARAMS, "PING", "Not enough parameters"))
		return true
	}
	client.Send(MakeMsg(self, "PONG", self.Name, msg.Params[0]))
	return true
}
