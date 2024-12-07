package main

func trimInitialSpaces(line string) string {
	var i int
	for i = 0; i < len(line) && line[i] == ' '; i++ {
	}
	return line[i:]
}
