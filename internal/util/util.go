package util

import "log"

func CheckError(e error) {
	if e != nil {
		log.Fatal(e)
	}
}