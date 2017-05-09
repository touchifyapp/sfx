// +build verbose

package main

import (
	"log"
)

func verbosef(format string, stuff ...interface{}) {
	log.Printf(format, stuff...)
}

func verboseFatal(err interface{}) {
	log.Fatal(err)
}
