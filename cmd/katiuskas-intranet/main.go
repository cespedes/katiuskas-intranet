package main

import (
	"log"
	"os"

	"github.com/cespedes/katiuskas-intranet/katintranet"
)

func main() {
	if err := katintranet.Run(os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
