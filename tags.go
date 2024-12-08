// Almost everything in this file is adapted from Ergo IRCd
// This is probably considered a derived work for copyright purposes
//
// SPDX-License-Identifier: MIT

package main

import (
	"strings"
	"unicode/utf8"
)

func parseTags(tagsString string) (tags map[string]string, err error) {
	tags = make(map[string]string)
	for 0 < len(tagsString) {
		tagEnd := strings.IndexByte(tagsString, ';')
		endPos := tagEnd
		nextPos := tagEnd + 1
		if tagEnd == -1 {
			endPos = len(tagsString)
			nextPos = len(tagsString)
		}
		tagPair := tagsString[:endPos]
		equalsIndex := strings.IndexByte(tagPair, '=')
		var tagName, tagValue string
		if equalsIndex == -1 {
			// Tag with no value
			tagName = tagPair
		} else {
			tagName, tagValue = tagPair[:equalsIndex], tagPair[equalsIndex+1:]
		}
		// "Implementations [...] MUST NOT perform any validation that would
		//  reject the message if an invalid tag key name is used."
		if validateTagName(tagName) {
			// "Tag values MUST be encoded as UTF8."
			if !utf8.ValidString(tagValue) {
				err = ErrInvalidTagContent
				return
			}
			tags[tagName] = UnescapeTagValue(tagValue)
		}
		tagsString = tagsString[nextPos:]
	}
	return
}

func UnescapeTagValue(inString string) string {
	// buf.Len() == 0 is the fastpath where we have not needed to unescape any chars
	var buf strings.Builder
	remainder := inString
	for {
		backslashPos := strings.IndexByte(remainder, '\\')

		if backslashPos == -1 {
			if buf.Len() == 0 {
				return inString
			} else {
				buf.WriteString(remainder)
				break
			}
		} else if backslashPos == len(remainder)-1 {
			// trailing backslash, which we strip
			if buf.Len() == 0 {
				return inString[:len(inString)-1]
			} else {
				buf.WriteString(remainder[:len(remainder)-1])
				break
			}
		}

		// Non-trailing backslash detected; we're now on the slowpath
		// where we modify the string
		if buf.Len() < len(inString) {
			buf.Grow(len(inString))
		}
		buf.WriteString(remainder[:backslashPos])
		buf.WriteByte(escapedCharLookupTable[remainder[backslashPos+1]])
		remainder = remainder[backslashPos+2:]
	}

	return buf.String()
}

var escapedCharLookupTable [256]byte

func init() {
	for i := 0; i < 256; i += 1 {
		escapedCharLookupTable[i] = byte(i)
	}
	escapedCharLookupTable[':'] = ';'
	escapedCharLookupTable['s'] = ' '
	escapedCharLookupTable['r'] = '\r'
	escapedCharLookupTable['n'] = '\n'
}

// https://ircv3.net/specs/extensions/message-tags.html#rules-for-naming-message-tags
func validateTagName(name string) bool {
	if len(name) == 0 {
		return false
	}
	if name[0] == '+' {
		name = name[1:]
	}
	if len(name) == 0 {
		return false
	}
	// Let's err on the side of leniency here; allow -./ (45-47) in any position
	for i := 0; i < len(name); i++ {
		c := name[i]
		if !(('-' <= c && c <= '/') || ('0' <= c && c <= '9') || ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z')) {
			return false
		}
	}
	return true
}
