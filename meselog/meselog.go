package meselog

import (
	"fmt"
)

func log(str string, keyvals []any) {
	fmt.Print(str + " ")
	for i, j := range keyvals {
		if i&1 == 0 {
			fmt.Printf("%v=", j)
		} else if i == len(keyvals)-1 {
			fmt.Printf("%#v", j)
		} else {
			fmt.Printf("%#v ", j)
		}
	}
	fmt.Print("\n")
}

func Error(str string, keyvals ...any) {
	log("ERROR "+str, keyvals)
}

func Debug(str string, keyvals ...any) {
	log("DEBUG "+str, keyvals)
}
