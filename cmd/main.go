package main

import (
	"flag"
	"os"
)

// @title ewallet!
// @version 1.0
// @host localhost:8888
// schemes https
func main() {
	var methodFlag = flag.String("method", "server", "select run method for flag method")

	flag.Parse()
	if !flag.Parsed() {
		flag.PrintDefaults()
		os.Exit(2)
	}

	switch *methodFlag {
	case "server":
		RunServer()
	case "dummy_user":
		DummyUser()
	default:
		RunServer()
	}
}
