package main

import (
	"flag"
	"log"

	"github.com/jfyne/enpass2gopass"
)

func main() {
	location := flag.String("file", "enpass_output.json", "Your enpass JSON export file")
	flag.Parse()

	export, err := enpass2gopass.NewExport(*location)
	if err != nil {
		log.Fatal("Cannot read export file", err)
	}

	export.Transfer()
}
