package main

func trimInitialSpaces(line string) string {
	var i int
	for i = 0; i < len(line) && line[i] == ' '; i++ {
	}
	return line[i:]
}

func isASCII(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] > 127 {
			return false
		}
	}
	return true
}
