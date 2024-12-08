package main

import (
	"strings"
)

type RMsg struct {
	RawSource  string
	Command    string
	Tags       map[string]string
	Params     []string
}

type Sourceable interface {
	ClientSource() string
	ServerSource() string
}

type SMsg struct {
	Source     *Sourceable
	Command    string
	Tags       map[string]string
	Params     []string
}

func (msg *SMsg) ClientSerialize() (final string) {
	if msg.Tags != nil && len(msg.Tags) != 0 {
		final = "@"
		for k, v := range msg.Tags{
			// TODO: Tag values must be escaped
			final += k + "=" + v + ";"
		}
		final += " "
	}
	if msg.Source != nil {
		final += ":" + (*msg.Source).ClientSource() + " "
	}
	final += msg.Command + " "

	if len(msg.Params) > 0 {
		for i := 0; i < len(msg.Params) - 1; i++ {
			final += msg.Params[i] + " "
		}
		final += ":" + msg.Params[len(msg.Params) - 1]
	}
	final += "\n"
	return
}

// Partially adapted from https://github.com/ergochat/irc-go.git
func parseIRCMsg(line string) (msg RMsg, err error) {
	msg = RMsg{}

	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")

	if len(line) == 0 {
		err = ErrEmptyMessage
		return
	}

	for _, v := range line {
		if v == '\x00' || v == '\r' || v == '\n' {
			err = ErrIllegalByte
			return
		}
	}

	// IRCv3 tags
	if line[0] == '@' {
		tagEnd := strings.IndexByte(line, ' ')
		if tagEnd == -1 {
			err = ErrEmptyMessage
			return
		}
		tagsString := line[1:tagEnd]
		if 0 < MaxlenTagData && MaxlenTagData < len(tagsString) {
			err = ErrTagsTooLong
			return
		}
		msg.Tags, err = parseTags(tagsString)
		if err != nil {
			return
		}
		// Skip over the tags and the separating space
		line = line[tagEnd+1:]
	}

	if len(line) > MaxlenBody {
		err = ErrBodyTooLong
		line = line[:MaxlenBody]
	}

	line = trimInitialSpaces(line)

	// Source
	if 0 < len(line) && line[0] == ':' {
		sourceEnd := strings.IndexByte(line, ' ')
		if sourceEnd == -1 {
			err = ErrEmptyMessage
			return
		}
		msg.RawSource = line[1:sourceEnd]
		// Skip over the source and the separating space
		line = line[sourceEnd+1:]
	}

	// Command
	commandEnd := strings.IndexByte(line, ' ')
	paramStart := commandEnd + 1
	if commandEnd == -1 {
		commandEnd = len(line)
		paramStart = len(line)
	}
	baseCommand := line[:commandEnd]
	if len(baseCommand) == 0 {
		err = ErrEmptyMessage
		return
	}
	// TODO: Actually must be either letters or a 3-digit numeric
	if !isASCII(baseCommand) {
		err = ErrIllegalByte
		return
	}
	msg.Command = strings.ToUpper(baseCommand)
	line = line[paramStart:]

	// Other arguments
	for {
		line = trimInitialSpaces(line)
		if len(line) == 0 {
			break
		}
		// Trailing
		if line[0] == ':' {
			msg.Params = append(msg.Params, line[1:])
			break
		}
		paramEnd := strings.IndexByte(line, ' ')
		if paramEnd == -1 {
			msg.Params = append(msg.Params, line)
			break
		}
		msg.Params = append(msg.Params, line[:paramEnd])
		line = line[paramEnd+1:]
	}

	return
}
