package main

var commandHandlers = map[string](func(RMsg, *Client) bool){}

/* Maybe we should make command handlers return their values for easier labelled-reply? */
