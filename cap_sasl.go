package main

import (
	"bytes"
	"encoding/base64"
)

type ExtraSasl struct {
	AuthMethod string
}

const (
	RPL_LOGGEDIN    = "900"
	RPL_LOGGEDOUT   = "901"
	ERR_NICKLOCKED  = "902"
	RPL_SASLSUCCESS = "903"
	ERR_SASLFAIL    = "904"
	ERR_SASLTOOLONG = "905"
	ERR_SASLABORTED = "906"
	ERR_SASLALREADY = "907"
	RPL_SASLMECHS   = "908"
)

const (
	panicSaslMethod = "stored illegal SASL method"
)

func init() {
	Caps["sasl"] = "PLAIN,EXTERNAL"
	CommandHandlers["AUTHENTICATE"] = handleClientAuthenticate
}

func handleClientAuthenticate(msg RMsg, client *Client) error {
	_, ok := client.Caps["sasl"]
	if !ok {
		return client.Send(MakeMsg(self, ERR_SASLFAIL, client.Nick, "SASL authentication failed (capability not requested)"))
	}

	if len(msg.Params) < 1 {
		return client.Send(MakeMsg(self, ERR_NEEDMOREPARAMS, "AUTHENTICATE", "Not enough parameters"))
	}

	extraSasl_, ok := client.Extra["sasl"]
	if !ok {
		client.Extra["sasl"] = &ExtraSasl{}
		extraSasl_ = client.Extra["sasl"]
	}
	extraSasl, ok := extraSasl_.(*ExtraSasl)
	if !ok {
		panic(panicType)
	}

	switch extraSasl.AuthMethod {
	case "":
		if msg.Params[0] != "PLAIN" && msg.Params[0] != "EXTERNAL" {
			return client.Send(MakeMsg(self, ERR_SASLFAIL, client.Nick, "SASL authentication failed (invalid method)"))
		}
		extraSasl.AuthMethod = msg.Params[0]
		return client.Send(MakeMsg(self, "AUTHENTICATE", "+"))
	case "*": // Abort
		extraSasl.AuthMethod = ""
		return client.Send(MakeMsg(self, ERR_SASLFAIL, client.Nick, "SASL authentication failed (aborted)"))
	case "EXTERNAL":
		extraSasl.AuthMethod = ""
		return client.Send(MakeMsg(self, ERR_SASLFAIL, client.Nick, "SASL authentication failed"))
	case "PLAIN":
		extraSasl.AuthMethod = ""
		saslPlainData, err := base64.StdEncoding.DecodeString(msg.Params[0])
		if err != nil {
			return client.Send(MakeMsg(self, ERR_SASLFAIL, client.Nick, "SASL authentication failed (base64 decoding error)"))
		}
		saslPlainSegments := bytes.Split(saslPlainData, []byte{0})
		if len(saslPlainSegments) != 3 {
			return client.Send(MakeMsg(self, ERR_SASLFAIL, client.Nick, "SASL authentication failed (not three segments)"))
		}
		_ = string(saslPlainSegments[0]) // authzid unused
		authcid := string(saslPlainSegments[1])
		passwd := string(saslPlainSegments[2])
		if authcid == "runxiyu" && passwd == "hunter2" {
			return client.Send(MakeMsg(self, RPL_SASLSUCCESS, client.Nick, "SASL authentication successful"))
		}
		return client.Send(MakeMsg(self, ERR_SASLFAIL, client.Nick, "SASL authentication failed"))
	default:
		panic(panicSaslMethod)
	}
}
