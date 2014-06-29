package main

import (
	"encoding/json"
	"os"

	"github.com/concourse/tracker-resource/out"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <sources directory>")
		os.Exit(1)
	}

	// sources := os.Args[1]

	var request out.OutRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	json.NewEncoder(os.Stdout).Encode(out.OutResponse{})
}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}
