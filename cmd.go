package main

var CommandHandlers = map[string](func(RMsg, *Client) error){}

/* Maybe we should make command handlers return their values for easier labelled-reply? */
