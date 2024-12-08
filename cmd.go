package main

var commandHandlers = map[string](func(RMsg, *Client) (error)){}

/* Maybe we should make command handlers return their values for easier labelled-reply? */
