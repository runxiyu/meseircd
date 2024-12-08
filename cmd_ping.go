package main

func init() {
	commandHandlers["PING"] = handleClientPing
}

func handleClientPing(msg RMsg, client *Client) (error) {
	if len(msg.Params) < 1 {
		client.Send(SMsg{Command: ERR_NEEDMOREPARAMS, Params: []string{"PING", "Not enough parameters"}})
	}
	client.Send(SMsg{Command: "PONG", Params: []string{msg.Params[0]}})
	return nil
}
