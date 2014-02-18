package bitverse

import (
	"log"
)

var debugFlag bool = false

func debug(str string) {
	if debugFlag {
		log.Println(str)
	}
}

func info(str string) {
	log.Println(str)
}

func fatal(str string) {
	log.Println(str)
}
